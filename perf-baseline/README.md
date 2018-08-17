# Performance baseline

Demonstrates baseline performance metrics for the Linkerd2 Proxy.

<img width="1228" alt="l5d-perf" src="https://user-images.githubusercontent.com/236915/44286934-ba6a1380-a21f-11e8-9055-3ed294fe637f.png">

## Deploy

```bash
cat perf-baseline.yaml | kubectl apply -f -
```

## Observe

### Reporter

The [`reporter/`](reporter/) process queries Prometheus and prints a report to
stdout:

```bash
kubectl -n l5d-perf logs -f $(kubectl -n l5d-perf get po --selector=app=reporter -o jsonpath='{.items[*].metadata.name}')
```

### Grafana

```bash
kubectl -n l5d-perf port-forward $(kubectl -n l5d-perf get po --selector=app=grafana -o jsonpath='{.items[*].metadata.name}') 3000:3000
open http://localhost:3000
```

### Prometheus

```bash
kubectl -n l5d-perf port-forward $(kubectl -n l5d-perf get po --selector=app=prometheus -o jsonpath='{.items[*].metadata.name}') 9090:9090
open http://localhost:9090
```

## Test setup

- each pod has 3 containers: load generator -> linkerd2 -> backend
- 1000 RPS spread across 10 connections
- HTTP/1.1 load via [slow_cooker](https://github.com/BuoyantIO/slow_cooker)
- HTTP/2 load via [strest-grpc](https://github.com/BuoyantIO/strest-grpc)
- Observability via Prometheus and Grafana
- Baseline (no proxy) config for comparison

## Testing changes to the Linkerd2 proxy

For community members interested in testing their own performance changes in the
[Linkerd2 Proxy repo](https://github.com/linkerd/linkerd2-proxy):
1. Build a Docker image from the
   [Linkerd2 Proxy repo](https://github.com/linkerd/linkerd2-proxy).
2. In [`perf-baseline.yaml`](perf-baseline.yaml), replace references to
   `gcr.io/linkerd-io/proxy:*` with your image.

## Hardware requirements

This test suite boots:
- 2 Linkerd's
- 4 load testers at 1000 RPS
- 4 backend servers
- 1 Prometheus
- 1 Grafana

Reommended hardware:
- 16 cores
- 8GB memory

## Tuning

Achieving consistent performance results in a scheduled environment like
Kubernetes requires some tuning. Several strategies are available to help enable
this.

### Node affinity

Node affinity pins a pod to a specific node. For example, to pin your pod to
a node named `node2`, add this to your PodSpec:

```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - node2
```

### Taints and tolerations

Complementary to node affnity, taints prevent other processes from running on a
given node.

To set a taint on a node named `node2`:

```bash
kubectl taint nodes node2 dedicated=groupName:NoSchedule
```

Then in a PodSpec, specify tolerations for that taint:

```yaml
tolerations:
- key: "dedicated"
  operator: "Equal"
  value: "groupName"
  effect: "NoSchedule"
```

## Related configs

### Reporter

* [`reporter/`](reporter/)

This tool gathers performance metrics from Prometheus and prints a report to
stdout.

### Grafana

* [`grafana/`](grafana/)

Files to build a Grafana image with dashboards specific to Linkerd performance
testing.

### Linkerd1

* [`linkerd1-perf/`](linkerd1-perf/)

Performance tests for Linkerd1.
