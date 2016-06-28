# Getting Started with linkerd and Docker Compose

This guide will walk you through getting linkerd running in a docker-compose
environment and routing requests to an nginx service using filesystem
based service discovery.

For more information on linkerd, please visit linkerd.io

## Downloading

Start by cloning this repo and pulling the nginx and linkerd docker images:

```
git clone https://github.com/BuoyantIO/linkerd-examples.git
cd linkerd-examples/getting-started/docker
docker pull nginx
docker pull buoyantio/linkerd:0.6.0
```

## Start nginx and linkerd

The included `docker-compose.yml` starts up an nginx service that serves static
content from the www directory and a linkerd.

```
docker-compose up -d
```

`linkerd.yaml` is a config file that controls linkerd's behavior:

```
# The filesystem namer (io.l5d.fs) watches the disco directory for changes.
# Each file in this directory represents a concrete name and contains a list
# of hostname/port pairs.
namers:
- kind: io.l5d.fs
  rootDir: disco

routers:
- protocol: http
  # Incoming requests to linkerd with a Host header of "hello" get assigned a
  # name like /http/1.1/GET/hello.  This dtab transforms that into
  # /#/io.l5d.fs/hello which indicates that the filesystem namer should be used
  # and should look for a file named "hello".  linkerd will then load balance
  # over the entries in that file.
  baseDtab: |
    /http/1.1/* => /#/io.l5d.fs
  servers:
  - ip: 0.0.0.0
    port: 4140

```

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
