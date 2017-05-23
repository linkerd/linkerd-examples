# Docker

This directory contains files and scripts for building custom Docker images that
are used in our demos.

## helloworld

* [`helloworld/`](helloworld/)

Builds the [buoyantio/helloworld](https://hub.docker.com/r/buoyantio/helloworld/)
image, which provides a go script for running the hello and world services.
These two services function together to make a highly scalable, "hello world"
microservice (where the hello service, naturally, calls the world service to
complete its request). For more information, see:

* [A Service Mesh for Kubernetes, Part I: Top-line service metrics](https://blog.buoyant.io/2016/10/04/a-service-mesh-for-kubernetes-part-i-top-line-service-metrics/)

## jenkins-plus

* [`jenkins-plus/`](jenkins-plus/)

Builds the [buoyantio/jenkins-plus](https://hub.docker.com/r/buoyantio/jenkins-plus/)
image, which provides the base jenkins image, along with the kubectl and
namerctl binaries that we need, as well as additional plugins and a
pre-configured pipeline job that we can use to run blue-green deployments. For
more information, see:

* [A Service Mesh for Kubernetes, Part IV: Continuous deployment via traffic shifting](https://blog.buoyant.io/2016/11/04/a-service-mesh-for-kubernetes-part-iv-continuous-deployment-via-traffic-shifting/)

## nginx

* [`nginx/`](nginx/)

Builds the [buoyantio/nginx](https://hub.docker.com/r/buoyantio/nginx/) image,
which provides the base nginx image, along with pre-installed modules that are
useful when using nginx to proxy requests to linkerd. For more information, see:

* [A Service Mesh for Kubernetes, Part V: Dogfood environments, ingress and edge routing](https://blog.buoyant.io/2016/11/18/a-service-mesh-for-kubernetes-part-v-dogfood-environments-ingress-and-edge-routing/)

## kubectl

The [`k8s-daemonset`](../k8s-daemonset/) examples make use of a
[kubectl Docker container](https://hub.docker.com/r/buoyantio/kubectl/), hosted
in the [buoyantio DockerHub repo](https://hub.docker.com/r/buoyantio/). To build
this container yourself, follow the instructions in the
[Kubernetes Repo](https://github.com/kubernetes/kubernetes/tree/master/examples/kubectl-container).
