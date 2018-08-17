package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	promApi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
)

type stats struct {
	sr        float64
	rr        uint64
	latencies map[float64]uint64
	mem       uint64
	cpu       float64
}

// protocol => app => stats
type statsReport map[string]map[string]*stats

const (
	successRateQuery = "sum(rate(successes{job=~\"strest-client|slow-cooker\"}[5m])) by (app, job) / sum(rate(requests{job=~\"strest-client|slow-cooker\"}[5m])) by (app, job)"
	requestRateQuery = "sum(rate(requests{job=~\"strest-client|slow-cooker\"}[5m])) by (app, job)"
	latencyQuery     = "histogram_quantile(%f, sum(rate(latency_us_bucket{job=~\"strest-client|slow-cooker\"}[5m])) by (le, app, job))"
	memoryQuery      = "sum(container_memory_working_set_bytes{}) by (container_name, pod_name)"
	cpuQuery         = "sum(rate(container_cpu_usage_seconds_total{}[5m])) by (container_name, pod_name)"
)

var (
	namespace *string

	protocolMap = map[model.LabelValue]string{
		"slow-cooker":   "h1",
		"strest-client": "h2",
	}
	latencies      = []float64{0.5, 0.75, 0.9, 0.95, 0.99, 0.999}
	dropContainers = []string{
		"grafana",
		"helloworld",
		"POD",
		"prometheus",
		"slow-cooker",
		"strest-client",
		"strest-server",
	}
)

func main() {
	namespace = flag.String("namespace", "l5d-perf", "namespace where the performance test is running")
	logLevel := flag.String("log-level", log.InfoLevel.String(), "log level, must be one of: panic, fatal, error, warn, info, debug")
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("Invalid log-level: %s", *logLevel)
	}
	log.SetLevel(level)

	// wait 7 minutes prior to first report (2 minutes for warmup, 5 minutes for test)
	// publish report every minute after that
	log.Infof("Waiting 7 minutes to publish first report for namespace %s...", *namespace)
	<-time.After(7 * time.Minute)

	err = publishReport()
	if err != nil {
		log.Errorf("Failed to publish report: %s", err)
	}

	for range time.Tick(time.Minute) {
		err := publishReport()
		if err != nil {
			log.Errorf("Failed to publish report: %s", err)
		}
	}
}

func publishReport() error {
	report, err := getStats()
	if err != nil {
		log.Errorf("Failed to retreive stats: %s", err)
		return err
	}

	fmt.Print(report)

	return nil
}

func getStats() (statsReport, error) {
	promURL := fmt.Sprintf("http://prometheus.%s.svc.cluster.local:9090", *namespace)

	log.Debugf("Connecting to Prometheus at %s", promURL)

	promClient, err := promApi.NewClient(promApi.Config{Address: promURL})
	if err != nil {
		log.Errorf("Failed to create Prometheus client: %s", err)
		return nil, err
	}

	promAPI := promv1.NewAPI(promClient)

	sr, err := queryProm(promAPI, successRateQuery)
	if err != nil {
		log.Errorf("Prometheus query (%s) failed: %s", successRateQuery, err)
		return nil, err
	}

	rr, err := queryProm(promAPI, requestRateQuery)
	if err != nil {
		log.Errorf("Prometheus query (%s) failed: %s", requestRateQuery, err)
		return nil, err
	}

	latencyResults := map[float64]model.Vector{}
	for _, latency := range latencies {
		query := fmt.Sprintf(latencyQuery, latency)
		latencyResults[latency], err = queryProm(promAPI, query)
		if err != nil {
			log.Errorf("Prometheus query (%s) failed: %s", query, err)
			return nil, err
		}
	}

	mem, err := queryProm(promAPI, memoryQuery)
	if err != nil {
		log.Errorf("Prometheus query (%s) failed: %s", memoryQuery, err)
		return nil, err
	}

	cpu, err := queryProm(promAPI, cpuQuery)
	if err != nil {
		log.Errorf("Prometheus query (%s) failed: %s", cpuQuery, err)
		return nil, err
	}

	// protocol => app => stats
	report := statsReport{}
	for _, protocol := range protocolMap {
		report[protocol] = map[string]*stats{}
	}

	for _, sample := range sr {
		app := string(sample.Metric["app"])
		if protocol, ok := protocolMap[sample.Metric["job"]]; ok {
			if report[protocol][app] == nil {
				report[protocol][app] = &stats{}
			}
			report[protocol][app].sr = float64(sample.Value)
		}
	}
	for _, sample := range rr {
		app := string(sample.Metric["app"])
		if protocol, ok := protocolMap[sample.Metric["job"]]; ok {
			if report[protocol][app] == nil {
				report[protocol][app] = &stats{}
			}
			report[protocol][app].rr = uint64(sample.Value)
		}
	}
	for latency, latencyResult := range latencyResults {
		for _, sample := range latencyResult {
			app := string(sample.Metric["app"])
			if protocol, ok := protocolMap[sample.Metric["job"]]; ok {
				if report[protocol][app] == nil {
					report[protocol][app] = &stats{}
				}
				if report[protocol][app].latencies == nil {
					report[protocol][app].latencies = map[float64]uint64{}
				}
				report[protocol][app].latencies[latency] = uint64(sample.Value)
			}
		}
	}
	for _, sample := range mem {
		container := string(sample.Metric["container_name"])
		for _, protocol := range protocolMap {
			if strings.Contains(string(sample.Metric["pod_name"]), protocol) {
				if report[protocol][container] == nil {
					report[protocol][container] = &stats{}
				}
				report[protocol][container].mem = uint64(sample.Value)
			}
		}
	}
	for _, sample := range cpu {
		container := string(sample.Metric["container_name"])
		for _, protocol := range []string{"h1", "h2"} {
			if strings.Contains(string(sample.Metric["pod_name"]), protocol) {
				if report[protocol][container] == nil {
					report[protocol][container] = &stats{}
				}
				report[protocol][container].cpu = float64(sample.Value)
			}
		}
	}

	for _, protocol := range protocolMap {
		for _, drop := range dropContainers {
			delete(report[protocol], drop)
		}
	}

	return report, nil
}

func queryProm(promAPI promv1.API, query string) (model.Vector, error) {
	log.Debugf("Prometheus query: %s", query)

	res, err := promAPI.Query(context.Background(), query, time.Time{})
	if err != nil {
		log.Errorf("Prometheus query(%+v) failed with: %+v", query, err)
		return nil, err
	}
	log.Debugf("Prometheus query response: %+v", res)

	if res.Type() != model.ValVector {
		err = fmt.Errorf("Unexpected query result type (expected Vector): %s", res.Type())
		log.Error(err)
		return nil, err
	}

	return res.(model.Vector), nil
}

func (s statsReport) String() string {
	buffer := bytes.NewBufferString("Stats report:\n")

	for protocol, protocolReport := range s {
		buffer.WriteString(fmt.Sprintf("Protocol: %s:\n", protocol))
		for app, stats := range protocolReport {
			buffer.WriteString(fmt.Sprintf("  %s\n", app))
			buffer.WriteString(fmt.Sprintf("    Success rate:     %3.2f%%\n", stats.sr*100))
			buffer.WriteString(fmt.Sprintf("    Request rate:     %d\n", stats.rr))
			for _, latency := range latencies {
				buffer.WriteString(fmt.Sprintf("    p%d latency (us): %d\n", int(latency*100), stats.latencies[latency]))
			}
			buffer.WriteString(fmt.Sprintf("    Memory (bytes):   %d\n", stats.mem))
			buffer.WriteString(fmt.Sprintf("    CPU (cores):      %.3f\n", stats.cpu))
		}
	}

	return buffer.String()
}
