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

### linkerd ingress configuration

To deploy a basic linkerd configuration with an ingress router running on port
4142 and an internal router running on port 4140, run:

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

```bash
dcos marathon group add webgroup.json
```

```bash
http_proxy=$PUBLIC_NODE:4140 curl webapp-a1.webapp-a.webgroup/hello
Hello world
```
