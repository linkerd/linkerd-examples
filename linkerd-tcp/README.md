# linkerd-tcp

This directory contains a docker-compose environment that runs a demo that you
can use to try out [linkerd-tcp](https://github.com/linkerd/linkerd-tcp).

## Startup

The [`docker-compose.yml`](docker-compose.yml) file that's included in this
directory is configured to run the demo. Start everything with:

```bash
$ docker-compose build && docker-compose up -d
```

That command will build and run all of the following containers:

* **namerd**: namerd is configured via the provided [`namerd.yml`](namerd.yml)
file, which defines a `default` namespace that can be used to route traffic
using the `io.l5d.fs` namer in conjunction with the files in the [`disco`](
disco/) directory.

* **linkerd**: linkerd is configured via the provided [`linkerd.yml`](
linkerd.yml), which specifies 1 HTTP router running on port 4140, routing HTTP
traffic via namerd.

* **linkerd-tcp**: linkerd-tcp is configured via the provided
[`linkerd-tcp.yml`](linkerd-tcp.yml), which specifies 1 TCP proxy running on
port 7474, routing Redis traffic via namerd.

* **2 Redis instances**: Two [Redis](https://redis.io/) instances are configured
to run on ports 6379 and 6380. The `default` namerd namespace is setup to send
all Redis traffic to the first redis instance, but that routing decision can be
changed by modifying namerd's dtab.

* **1 HTTP cluster**: The HTTP cluster consists of 10 instances of the HTTP
web server defined in [`server.go`](server.go). The web server is configured
to respond with the string "hello", and to cache its responses in Redis. If the
response is found in cache, it returns immediately. On cache miss, it sleeps
for 300 milliseconds before returning.

* **slow_cooker**: Traffic to linkerd is generated using [slow\_cooker](
https://github.com/BuoyantIO/slow_cooker). slow\_cooker is configured to send
500 requests per second to the linkerd HTTP router, which load balances the
requests over the 10 HTTP web server instances.

* **prometheus**: linkerd, linkerd-tcp, and the Redis instances expose metrics
data in a format that can be read by [Prometheus](https://prometheus.io/).
Prometheus metrics collection is configured in [`prometheus.yml`](
prometheus.yml), which scrapes all metrics from all processes every 10 seconds.

* **grafana**: Collected metrics are displayed on a dashboard using [Grafana](
http://grafana.org/). The grafana dasbhard is running on port 3000, and is
defined in [`grafana.json`](grafana.json).

* **linkerd-viz**: The setup is also running a [linkerd-viz](
https://github.com/BuoyantIO/linkerd-viz) on port 3000, which gives an overall
picture of the performance of slow_cooker's HTTP requests.

## Dashboards

### Grafana

Grafana is running on port 3000 in your docker-compose environment. To see a
dashboard that displays linkerd-tcp metrics alongside linkerd and redis metrics,
load the dashboard by going to port 3000 on your docker host. It should look
like this:

![grafana](screenshot.png)

### linkerd admin

The linkerd admin server is running on port 9990 in your docker-compose
environment, and displays instance-specific metrics data.

### linkerd-viz

The linkerd-viz server is running on port 3001 in your docker-compose
environment, and displays high-level performance data.

## Troubleshooting

If you have any issues getting the demo up and running, pop into [linkerd's
Slack]( https://slack.linkerd.io) and we'll help you get it sorted out.

Thanks! 👋
