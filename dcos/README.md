# Example configs for DC/OS

For more information see our
[DC/OS Getting Started Guide](https://linkerd.io/getting-started/dcos/).

## Deploy webapp

```bash
dcos marathon app add webapp.json
```

## Deploy linkerd

```bash
dcos marathon app add linkerd-dcos.json
```

## Test

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```
