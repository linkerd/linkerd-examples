# Getting Started with linkerd and Docker Compose

This guide will walk you through getting linkerd running in a docker-compose
environment and routing requests to an nginx service using filesystem
based service discovery.

For more information on linkerd, please visit linkerd.io

## Downloading

Start by cloning this repo:

```
git clone https://github.com/linkerd/linkerd-examples.git
cd linkerd-examples/getting-started/docker
```

## Start nginx and linkerd

The included `docker-compose.yml` starts up an nginx service that serves static
content from the www directory and a linkerd. Running docker-compose will pull
the required images automatically:

```
docker-compose up -d
```

When linkerd starts, it loads the [linkerd.yaml](linkerd.yaml) to configure
a linkerd router that routes traffic to the nginx service, using file-system
backed service discovery.

## Send requests

Now, when we send requests to linkerd, it will look the Host header up in the
disco directory (which has been mounted as a volume) to determine where the
request should be routed.  In this case we set the Host header to `hello` so
that the request is routed to the nginx service.

```
curl -H "Host: hello" <docker ip>:4140
```

## Admin dashboard

You can view the linkerd admin dashboard at `<docker ip>:9990`.
