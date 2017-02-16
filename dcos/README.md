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

```bash
dcos marathon app add linkerd-dcos.json
```

## Deploy namerd (for configurations that include namerd)

```bash
dcos marathon app add namerd-dcos.json
```

### Test namerd dtab interface

```bash
curl $PUBLIC_NODE:4180/api/1/dtabs/default
```

## Test

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### Test ingress configuration

```bash
$ curl $PUBLIC_NODE:4242/hello
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
