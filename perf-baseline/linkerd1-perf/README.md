# Linkerd1 performance testing

Test Linkerd1 performance against a released version. For demonstration
purposes, this example hard-codes an OpenJ9 variant of Linkerd1.

## Test setup

- each pod has 3 containers: load generator -> linkerd1 -> backend
- 1000 RPS spread across 10 connections
- HTTP/1.1 load via [slow_cooker](https://github.com/BuoyantIO/slow_cooker)
- HTTP/2 load via [strest-grpc](https://github.com/BuoyantIO/strest-grpc)
- Observability via Prometheus and Grafana
- Baseline (no proxy) config for comparison

## Deploy

```bash
cat linkerd1-perf.yaml | kubectl apply -f -
```

## Observe

### Grafana

```bash
kubectl -n linkerd1-perf port-forward $(kubectl -n linkerd1-perf get po --selector=app=grafana -o jsonpath='{.items[*].metadata.name}') 3000:3000
open http://localhost:3000
```

<img width="1164" alt="linkerd1-perf" src="https://user-images.githubusercontent.com/236915/43617284-ff3b1ee2-9675-11e8-8877-7d3bd5127045.png">

## Hardware requirements

This test suite boots:
- 4 Linkerd's
- 6 load testers at 1000 RPS
- 6 backend servers
- 1 Prometheus
- 1 Grafana

Reommended hardware:
- 32 cores
- 8GB memory

## Tuning

Depending on hardware, set `FINAGLE_WORKERS` to twice the number of physical
cores. Also setting `JVM_HEAP_MIN` and `JVM_HEAP_MAX` to a high value (and the
same value), can help with memory fragmentation and GC pressure. For example:

```yaml
image: buoyantio/linkerd:1.4.5
env:
- name: FINAGLE_WORKERS
  value: "32"
- name: JVM_HEAP_MIN
  value: 1024M
- name: JVM_HEAP_MAX
  value: 1024M
```

More details on Linkerd performance tuning may be found at:
https://discourse.linkerd.io/t/linkerd-performance-tuning/447
