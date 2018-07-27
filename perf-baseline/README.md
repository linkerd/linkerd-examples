# Performance baseline

Demonstrates baseline performance metrics for the Linkerd2 Proxy.

## Test setup

- 1 pod with 3 containers: load generator -> linkerd2 -> backend
- 1000qps spread across 10 connections
- HTTP/1.1 load via [slow_cooker](https://github.com/BuoyantIO/slow_cooker)
- HTTP/2 load via [strest-grpc](https://github.com/BuoyantIO/strest-grpc)
- Observability via Prometheus and Grafana
- Baseline (no proxy) config for comparison

## Deploy

```bash
cat perf-baseline.yaml | kubectl apply -f -
```

## Observe

### Prometheus

```bash
kubectl -n perf-baseline port-forward $(kubectl -n perf-baseline get po --selector=app=prometheus -o jsonpath='{.items[*].metadata.name}') 9090:9090
open http://localhost:9090
```

### Grafana

```bash
kubectl -n perf-baseline port-forward $(kubectl -n perf-baseline get po --selector=app=grafana -o jsonpath='{.items[*].metadata.name}') 3000:3000
open http://localhost:3000
```
