# Docker

This directory contains files and scripts for building custom Docker images that
are used in our demos.

## jenkins-plus

Builds the [buoyantio/jenkins-plus](https://hub.docker.com/r/buoyantio/jenkins-plus/)
image, which provides the base jenkins image, along with the kubectl and
namerctl binaries that we need, as well as additional plugins and a
pre-configured pipeline job that we can use to run blue-green deployments.
For an example of how to use this image in Kuberenetes, see our [blog post](
https://blog.buoyant.io/2016/11/04/a-service-mesh-for-kubernetes-part-iv-continuous-deployment-via-traffic-shifting/).

## NGINX

Builds the [buoyatio/nginx](https://hub.docker.com/r/buoyantio/jenkins-plus/)
image used in the [Ingress blog post](). This image runs nginx configured as an
ingress.  For an example of how to use this image in Kuberenetes, see our [blog post](
https://blog.buoyant.io/2016/11/04/a-service-mesh-for-kubernetes-part-iv-continuous-deployment-via-traffic-shifting/).

## NGINX with nginx-extras installed

Builds the [buoyatio/nginx](https://hub.docker.com/r/buoyantio/jenkins-plus/)
mentioned in the [Ingress blog post](). This provides nginx with 3rd party
modules installed.