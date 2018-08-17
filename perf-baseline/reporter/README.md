# Reporter

This tool gathers performance metrics from Prometheus and prints a report to
stdout. Intended to be run as part of the `perf-baseline` environment.

```bash
docker build . -t gcr.io/linkerd-io/reporter:latest
docker push gcr.io/linkerd-io/reporter:latest
```
