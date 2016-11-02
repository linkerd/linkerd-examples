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

## Deploy namerd (for namerd-linkerd configuration)

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
