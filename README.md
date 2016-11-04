# linkerd examples

The repo contains various examples for how to use linkerd and namerd.

* [Gob's microservice](gob/) is an example microservice application that uses
  linkerd and namerd to do staging, canary, and blue-green deploy
* The [plugins](plugins/) directory contains sample code for building linkerd
  plugins
* Several [getting started](getting-started/) guides for different environments
  including local development, docker-compose, Kubernetes, and Mesos.
* The [k8s-daemonset](k8s-daemonset/) directory contains a sample hello world
  app and multiple configs for deploying the app to Kubernetes in various
  configurations.
* [Dockerfiles](docker/) and configs for custom-built example images are
  provided in the `docker` directory.
