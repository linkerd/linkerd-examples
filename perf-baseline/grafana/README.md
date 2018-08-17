# Grafana config for Linkerd testing

This directory contains files needed to build a Grafana image with dashboards
specific to Linkerd performance testing.

```bash
docker build . -t gcr.io/linkerd-io/grafana-perf:latest
docker push gcr.io/linkerd-io/grafana-perf:latest
```
