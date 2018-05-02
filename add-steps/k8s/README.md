# Go faster by adding more steps (on Kubernetes)

This directory contains a Kubernetes config that runs a demo that you can use to
test linkerd's performance. The results of the demo are discussed in much more
detail in Buoyant's blog post,
[Making Things Faster by Adding More Steps](https://blog.buoyant.io/2017/01/31/making-things-faster-by-adding-more-steps/).

## Deploy

The [`add-steps.yml`](add-steps.yml) file that's included in this directory is
configured to run the demo. Deploy it to Kubernetes with:

```bash
curl https://raw.githubusercontent.com/linkerd/linkerd-examples/master/add-steps/k8s/add-steps.yml | kubectl apply -f -
```

## Demo Details

Here's an overview of the demo you just deployed:

![add-steps](../add-steps.png)

* **2 slow_cookers**: Traffic is generated using [slow\_cooker](
https://github.com/BuoyantIO/slow_cooker). The slow\_cookers are configured to
send 500 requests per second to both the linkerd cluster and the baseline
cluster.

* **linkerd**: linkerd is configured via a `linkerd.yml` config file, defined as
a `ConfigMap` in [`add-steps.yml`](add-steps.yml). This specifies a router
running on port 4140, configured to route traffic to a backend app cluster, via
Kubernetes service discovery.

* **2 backend clusters**: Each cluster runs 10 instances of the go web server
that's defined in [`server.go`](../server.go). One cluster is fronted by
linkerd, another baseline cluster simply uses a Kubernetes `Service` object.
Instance response times vary between 0 and 2 seconds. The slowest instances also
simulate a decrease in success rate.

* **prometheus**: linkerd, slow\_cooker, and the backend instances expose
metrics data in a format that can be read by
[Prometheus](https://prometheus.io/). Prometheus metrics collection is
configured via a `prometheus.yml` config file, defined as a `ConfigMap` in
[`add-steps.yml`](add-steps.yml). This file instructs Prometheus to scrape
metrics from all processes every 5 seconds.

* **grafana**: Collected metrics are displayed on a dashboard using [Grafana](
http://grafana.org/). The grafana container is preconfigured to display
dashboard comparing linkerd and the baseline cluster, defined via a
`making-things-faster-by-adding-more-steps.json` dashboard file, as a
`ConfigMap` in [`add-steps.yml`](add-steps.yml).

## View The Demo

Once up and running, open a port to the demo with:

```bash
kubectl -n add-steps port-forward $(kubectl -n add-steps get po --selector=app=grafana -o jsonpath='{.items[*].metadata.name}') 3000:3000
```

Then browse to http://localhost:3000 to view the demo. It will take about a
minute to begin reporting data:

![grafana](../screenshot-grafana.png)

## Linkerd Config Details

Linkerd's config file, `linkerd.yml`, is defined as a `ConfigMap` in the
[`add-steps.yml`](add-steps.yml) Kubernetes config.

```yaml
# Configures Linkerd’s administrative interface.
admin:
  ip: 0.0.0.0
  port: 9990

# Telemeters export metrics and tracing data about Linkerd, the services it
# connects to, and the requests it processes.
telemetry:
# Expose Prometheus style metrics on :9990/admin/metrics/prometheus
- kind: io.l5d.prometheus

# Usage is used for anonymized usage reporting. You can set the orgId to
# identify your organization or set `enabled: false` to disable entirely.
usage:
  orgId: linkerd-examples-add-steps-k8s
  enabled: true

# Namers provide Linkerd with service discovery information. To use a namer,
# you reference it in the dtab by its prefix.
namers:
# Enable the Kubernetes Service Discovery Namer, we will use `/#/io.l5d.k8s`
# in our dtab.
- kind: io.l5d.k8s
  host: localhost
  port: 8001

# Routers define how Linkerd actually handles traffic. Each router listens
# for requests, applies routing rules to those requests, and proxies them
# to the appropriate destinations. Each router is protocol specific.
routers:
# Configure one HTTP router. The application is expected to send traffic to
# this router on port 4140. Linkerd then proxies the request to the target
# application.
- protocol: http
  servers:
  - ip: 0.0.0.0
    port: 4140
  # Route requests based on Kubernetes Service Discovery information.
  # Specifically, map the HTTP host header to endpoints identified by a
  # service name in the `add-steps` namespace, on the `test-app` port.
  # For example:
  # http_proxy=localhost:4140 curl http://linkerd-app
  # routes to endpoints defined by the service object:
  # http://linkerd-app.add-steps.svc.cluster.local
  dtab: |
    /svc => /#/io.l5d.k8s/add-steps/test-app

  # This section defines the policy that Linkerd will use when talking to
  # services. The structure of this section depends on its `kind`.
  service:
    # `io.l5d.global` allows you to specifies parameters for all services.
    kind: io.l5d.global
    # A `responseClassifier` determines which HTTP responses should be
    # considered failures and which can be retried.
    responseClassifier:
      # All 5XX responses are considered to be failures. However, GET, HEAD,
      # OPTIONS, and TRACE requests may be retried automatically.
      kind: io.l5d.http.retryableRead5XX

  # This section defines how the clients that Linkerd creates will be
  # configured. The structure of this section depends on its kind.
  client:
    # `io.l5d.global` allows you to specifies parameters for all clients.
    kind: io.l5d.global
    # Specify a client-side load balancer.
    loadBalancer:
      # Specify an Exponentially Weighted Moving Average load balancer
      # algorithm.
      # More info: https://twitter.github.io/finagle/guide/Clients.html#power-of-two-choices-p2c-peak-ewma
      kind: ewma

    # Specify a Circuit Breaker, (aka Failure Accrual) policy
    # Linkerd uses failure accrual to track the number of requests that have
    # failed to a given node, and it will back off sending requests to any
    # nodes whose failures have exceeded a given threshold. Both the failure
    # threshold and the backoff behavior are configurable.
    failureAccrual:
      # Computes an exponentially-weighted moving average success rate for
      # each node, and backs off sending requests to nodes that have fallen
      # below the specified success rate. The window size for computing
      # success rate is constrained to a fixed number of requests.
      kind: io.l5d.successRate
      # Target success rate nodes must stay above to be marked alive.
      successRate: 0.9
      # Compute success rate over the last 20 requests.
      requests: 20
      # Once a node is marked dead, it will attempt to resend it traffic in
      # every 10 seconds.
      backoff:
        kind: constant
        ms: 10000
```

## Dashboards

### Prometheus

To view the Prometheus dashboard for raw metrics from all processes:

```bash
kubectl -n add-steps port-forward $(kubectl --namespace=add-steps get po --selector=app=prometheus -o jsonpath='{.items[*].metadata.name}') 9090:9090
```

Then browse to http://localhost:9090 to view Prometheus.

### Linkerd admin

To view the Linkerd admin dashboard:

```bash
kubectl -n add-steps port-forward $(kubectl --namespace=add-steps get po --selector=app=l5d -o jsonpath='{.items[*].metadata.name}') 9990:9990
```

Then browse to http://localhost:9990 to view the Linkerd Admin Dashboard.

## Teardown

```bash
kubectl delete ns add-steps
```

## Troubleshooting

If you have any issues getting the demo up and running, pop into [linkerd's
Slack]( https://slack.linkerd.io) and we'll help you get it sorted out.

Thanks! 👋

## Build

If you make changes to [`server.go`](../server.go):

```bash
docker build ../ -t buoyantio/add-steps-app:v1
docker push buoyantio/add-steps-app:v1
```
