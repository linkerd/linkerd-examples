package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requests = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "requests",
	Help: "Number of requests",
})

type handler struct {
	successRate float64
	latency     time.Duration
}

func printError(desc string, e error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", desc, e)
}

func (h *handler) HandleRequest(w http.ResponseWriter, req *http.Request) {

	requests.Inc()
	time.Sleep(h.latency)

	// if erroring, just return
	if rand.Float64() > h.successRate {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal service error"))
		return
	}

	// just return success
	w.Write([]byte("pong"))
}

func init() {
	prometheus.MustRegister(requests)
}

func main() {
	addr := flag.String("addr", ":8501", "service port to run on")
	successRate := flag.Float64("success-rate", 1.0, "service success rate")
	latency := flag.Duration("latency", time.Duration(0), "latency to add to each request")
	flag.Parse()

	fmt.Printf("serving on %s, success rate: %2f, latency: %v", *addr, *successRate, *latency)

	httpHandler := handler{successRate: *successRate, latency: *latency}
	http.HandleFunc("/", httpHandler.HandleRequest)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*addr, nil)
}
