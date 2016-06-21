# Getting Started with linkerd Locally

This guide will walk you through getting linkerd running locally and routing
requests to a basic python service using filesystem based service discovery.

For more information on linkerd, please visit linkerd.io

## Downloading

Start by cloning this repo and downloading linkerd:

```
git clone https://github.com/BuoyantIO/linkerd-examples.git
cd linkerd-examples/getting-started/local
curl -sLO https://github.com/BuoyantIO/linkerd/releases/download/0.6.0/linkerd-0.6.0-exec
chmod +x linkerd-0.6.0-exec
```

## Start the python service and linkerd

We start a simple python http service on port 8888 and linkerd on port 4140

```
python3 -m http.server 8888 &
./linkerd-0.6.0-exec linkerd.yaml &
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
disco directory to determine where the request should be routed.  In this case
we set the Host header to `hello` so that the request is routed to our python
service.

```
curl -H "Host: hello" 0:4140
```

## Admin dashboard

You can view the linkerd admin dashboard at `localhost:9990`.

