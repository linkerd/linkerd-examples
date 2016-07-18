# Getting Started with linkerd and DC/OS

This guide will walk you through getting linkerd running in DC/OS and
routing requests to an nginx service.

For more information on linkerd on DC/OS, please see our
[blog post](https://blog.buoyant.io/2016/04/19/linkerd-dcos-microservices-in-production-made-easy/)
or visit linkerd.io.

## Downloading

Start by cloning this repo.  Make sure that the `dcos` cli is installed and
configured to talk to your DC/OS cluster.

```
git clone https://github.com/BuoyantIO/linkerd-examples.git
cd linkerd-examples/getting-started/dcos
dcos config show
```

## Starting nginx

We create a simple nginx app that simply serves a static file on port 80 using
the standard nginx docker image.

```
dcos marathon app add nginx.json
```

## Starting linkerd

We will be installing linkerd as a DC/OS Universe package.  Before we can do
that, we need to create a linker config file that controls linkerd's behavior.
For this example, we'll use the `linkerd.yaml` in this repo, hosted on github.

```yaml
# The Marathon namer (io.l5d.marathon) queries the Marathon master
# for a list of addresses for a given app.
namers:
- kind:         io.l5d.marathon
  experimental: true
  prefix:       /io.l5d.marathon
  host:         marathon.mesos
  port:         8080
routers:
- protocol: http
  # Incoming requests to linkerd with a Host header of "hello" get assigned
  # a name like /http/1.1/GET/hello.  This dtab transforms that into
  # /#/io.l5d.marathon/hello which indicates that the marathon namer should
  # query the API for addresses for the app named "hello".  linkerd will then
  # load balance over those addresses.
  baseDtab: |
    /http/1.1/* => /#/io.l5d.marathon
  servers:
  - port: 4140
    ip: 0.0.0.0
admin:
  port: 9990
```

We will reference this config file in the package options when installing the
linkerd DC/OS package.  This can be done through the DC/OS Universe UI, but for
this example we'll do it on the command line:

```
dcos package install linkerd --options=linkerd-package-options.json
```

These options specify to install linkerd on a single public node.  For a
production deployment, we recommend running linkerd on every node (see our
[linkerd on DC/OS blog post](https://blog.buoyant.io/2016/04/19/linkerd-dcos-microservices-in-production-made-easy/)).

## Send requests

Now, when we send requests to linkerd, it will look for a service with the same
name as the Host header to determine where the request should be routed.  In
this case we set the Host header to `nginx` so that the request is routed to the
nginx service.

```
curl -H "Host: nginx" <public node ip>:4140
```

## Admin dashboard

You can view the linkerd admin dashboard at `<public node ip>:9990`.
