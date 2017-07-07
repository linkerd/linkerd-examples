# Example configs for DC/OS

For more information see our
[DC/OS Getting Started Guide](https://linkerd.io/getting-started/dcos/).

## Deploy webapp

```bash
dcos marathon app add webapp.json
```

## Deploy linkerd

Note the `linkerd-dcos.json` files assume 4 nodes. Modify this to equal the
total number of public+private nodes in your cluster.

Multiple linkerd configurations are described below. Pick the one that's most
appropriate for your setup. When testing configurations, be sure to set the
`PUBLIC_NODE` env variable to the external address of the public node in your
cluster.

### linkerd simple proxy

To deploy the most basic configuration, with linkerd as a proxy running on port
4140 for inbound requests, run:

```bash
dcos marathon app add simple-proxy/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### linkerd ingress configuration

To deploy linkerd with an ingress router running on port 4142 and an internal
router running on port 4140, run:

```bash
dcos marathon app add ingress/linkerd-dcos.json
```

Test this configuration with:

```bash
$ curl $PUBLIC_NODE:4242/hello
Hello world
```

### linkerd in linker-to-linker mode

To deploy linkerd in linker-to-linker mode, with outgoing traffic served on a
router running on port 4140, and incoming traffic served on a router running on
port 4141, run:

```bash
dcos marathon app add linker-to-linker/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### linkerd with namerd

#### namerd

Start by deploying namerd:

```bash
dcos marathon app add namerd/namerd-dcos.json
```

Test the namerd configuration with:

```bash
$ curl $PUBLIC_NODE:4180/api/1/dtabs/default
[{"prefix":"/marathonId","dst":"/#/io.l5d.marathon"},{"prefix":"/svc","dst":"/$/io.buoyant.http.domainToPathPfx/marathonId"}]
```

#### linkerd

Next deploy linkerd configured to talk to namerd when routing requests:

```bash
dcos marathon app add linkerd-with-namerd/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### linkerd with namerd in linker-to-linker mode

Deploy namerd as described in the previous section. Then deploy linkerd in
linker-to-linker mode, configured to talk to namerd when routing requests:

```bash
dcos marathon app add linker-to-linker-with-namerd/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### Application Groups

Marathon supports an "Application Group" concept, where applications are
deployed and named using a hierarchical path-based naming structure. Because the
linkerd config examples documented here all use the `domainToPathPfx` rewriting
namer, marathon applications within a group are routed by reversing the group
name into a domain-like name. For example, `webgroup/webapp-a/webapp-a1` becomes `webapp-a1.webapp-a.webgroup`:

#### Webgroup

This example demonstrates linkerd routing requests to a Marathon app in an application group.

```bash
dcos marathon group add webgroup.json
```

```bash
http_proxy=$PUBLIC_NODE:4140 curl webapp-a1.webapp-a.webgroup/hello
Hello world
```

#### Hello World

This example demonstrates inter-service routing, along with a routing override.

Deploy 3 services: `hello`, `world-v1`, `world-v2`:

```bash
dcos marathon group add hello-world.json
```

Route requests `linkerd` -> `hello` -> `linkerd` -> `world-v1`:

```bash
http_proxy=$PUBLIC_NODE:4140 curl hello.hw.buoyant
Hello (10.0.3.80) world (10.0.1.148)!
```

Routing override from `world-v1` to `world-v2`:

```bash
# 25% to world-v2
http_proxy=$PUBLIC_NODE:4140 curl -H 'l5d-dtab: /svc/world-v1.hw.buoyant => 3 * /marathonId/buoyant/hw/world-v1 & /marathonId/buoyant/hw/world-v2' hello.hw.buoyant
Hello (10.0.1.56) world (10.0.1.56)!!

# 75% to world-v2
http_proxy=$PUBLIC_NODE:4140 curl -H 'l5d-dtab: /svc/world-v1.hw.buoyant => /marathonId/buoyant/hw/world-v1 & 3 * /marathonId/buoyant/hw/world-v2' hello.hw.buoyant
Hello (10.0.1.56) earth (10.0.1.56)!!

# 100% to world-v2
http_proxy=$PUBLIC_NODE:4140 curl -H 'l5d-dtab: /svc/world-v1.hw.buoyant => /svc/world-v2.hw.buoyant' hello.hw.buoyant
Hello (10.0.1.56) earth (10.0.1.56)!!
```
