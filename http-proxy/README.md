# Example config for linkerd as an HTTP Proxy

For more information see our
[HTTP Proxy documentation](https://linkerd.io/getting-started/http-proxy/).

## Setup webapp

```bash
echo "Hello world" > hello; python3 -m http.server 8888
```

## Setup linkerd

```bash
curl -sLO https://github.com/linkerd/linkerd/releases/download/1.7.0/linkerd-1.7.0-exec
chmod +x linkerd-1.7.0-exec
./linkerd-1.7.0-exec ./linkerd.yaml
```

## Test

```bash
$ http_proxy=localhost:4140 curl -s http://webapp/hello
Hello world
```
