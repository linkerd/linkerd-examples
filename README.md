![linkerd][l5d-logo]

[![GitHub license][license-badge]](LICENSE)
[![Circle CI][l5d-ci-badge]][l5d-ci]

# linkerd examples

🎈 Welcome to linkerd-examples! 👋

This repo contains subdirectories with various examples for how to use linkerd
and namerd. Each of the subdirectories is described below.

## Getting started

* [`getting-started/`](getting-started/)

Provides guides for getting linkerd up and running in multiple different
environments, including local development, docker-compose, DC/OS, and
Kubernetes. More information:

* [Getting Started: Running with Docker](https://linkerd.io/getting-started/docker/)
* [Getting Started: Running in Kubernetes](https://linkerd.io/getting-started/k8s/)

## Add steps demo

* [`add-steps/`](add-steps/)

Provides a self-contained docker-compose environment that can be used to test
linkerd's performance. More information:

* [Making Things Faster by Adding More Steps](https://blog.buoyant.io/2017/01/31/making-things-faster-by-adding-more-steps/)

## Consul

* [`consul/`](consul/)

Provides a self-contained docker-compose environment which shows how to
use consul as a service discovery back-end with linkerd.

## DC/OS

* [`dcos/`](dcos/)

Provides common configurations for deploying linkerd and namerd to DC/OS. More
information:

* [Getting Started: Running in DC/OS](https://linkerd.io/getting-started/dcos/)
* [linkerd on DC/OS for Service Discovery and Visibility](https://blog.buoyant.io/2016/10/10/linkerd-on-dcos-for-service-discovery-and-visibility/)

## Docker

* [`docker/`](docker/)

Contains files and scripts for building custom Docker images that are used in
some of the examples in this repo.

## Amazon ECS

* [`ecs/`](ecs/)

Provides common configurations for deploying linkerd to Amazon ECS. More
information:

* [Getting Started: Running in ECS](https://linkerd.io/getting-started/ecs/)

## Failure accrual demo

* [`failure-accrual/`](failure-accrual/)

Provides a self-contained docker-compose environment that can be used to test
various failure accrual settings. More information:

* [Making microservices more resilient with advanced circuit breaking](https://blog.buoyant.io/2017/01/13/making-microservices-more-resilient-with-circuit-breaking/)

## Gob's microservice

* [`gob/`](gob/)

Defines an example microservice application that uses linkerd and namerd to do
staging, canary, and blue-green deploy.

## http_proxy

* [`http-proxy/`](http-proxy/)

Contains a linkerd configuration file that demonstrates how to make requests
through linkerd using the http_proxy environment variable. More information:

* [Features: HTTP proxy integration](https://linkerd.io/features/http-proxy/)

## influxdb

* [`influxdb/`](influxdb/)

Sets up a demo environment that configures a
[linkerd-viz](https://github.com/linkerd/linkerd-viz)-like dashboard using
linkerd, Telegraf, InfluxDB, and Grafana. Provides helpful configuration files
and dashboards for all components.

## Istio

* [`istio/`](istio/)

Provides common configurations for deploying linkerd to Istio. More information:

* [Getting Started: Running with Istio](https://linkerd.io/getting-started/istio/)
* [Linkerd and Istio: Like Peanut Butter and Jelly](https://buoyant.io/2017/07/11/linkerd-istio/)

## A Service Mesh for Kubernetes

* [`k8s-daemonset/`](k8s-daemonset/)

Defines a sample hello world app and multiple configs for deploying the app
to Kubernetes in various configurations, in support for Buoyant's "A Service
Mesh for Kubernetes" series of blog posts. More information:

* [A Service Mesh for Kubernetes, Part I: Top-Line Service Metrics](https://blog.buoyant.io/2016/10/04/a-service-mesh-for-kubernetes-part-i-top-line-service-metrics/)
* [Getting Started: Running in Kubernetes with DaemonSets](https://linkerd.io/getting-started/k8s-daemonset/)

## linkerd-tcp

* [`linkerd-tcp/`](linkerd-tcp/)

Sets up a demo environment that uses linkerd to route HTTP traffic and
[linkerd-tcp](https://github.com/linkerd/linkerd-tcp) to route Redis traffic.
Provides helpful configuration files and dashboards for all components.

## Mesos + Marathon

* [`mesos-marathon/`](mesos-marathon/)

Demonstrates running linkerd and a sample hello world app on Mesos and Marathon,
without DC/OS.

## Plugins

* [`plugins/`](plugins/)

Contains sample code for building linkerd plugins. More information:

* [In Depth: Plugins](https://linkerd.io/in-depth/plugin/)

## Testing

```bash
docker run -v `pwd`:/root/linkerd-examples --entrypoint=/root/linkerd-examples/.circleci/ci.sh buoyantio/linkerd:1.3.6
```

<!-- references -->
[l5d-ci]: https://circleci.com/gh/linkerd/linkerd-examples
[l5d-ci-badge]: https://circleci.com/gh/linkerd/linkerd-examples.svg?style=shield
[l5d-logo]: https://user-images.githubusercontent.com/9226/33582867-3e646e02-d90c-11e7-85a2-2e238737e859.png
[license-badge]: https://img.shields.io/github/license/linkerd/linkerd-examples.svg
