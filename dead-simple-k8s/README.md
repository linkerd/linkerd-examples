# Dead Simple Kubernetes with linkerd

## Deploy linkerd

```bash
kubectl apply -f linkerd.yml
```

## Deploy Hello World

```bash
kubectl apply -f hello-world.yml
```

## Test Requests

```bash
http_proxy=$(kubectl get svc | grep l5d | awk '{ print $3 }'):4140 curl -s http://hello
http_proxy=$(kubectl get svc | grep l5d | awk '{ print $3 }'):4140 curl -s http://world
```

## Deploy linkerd-viz

```bash
kubectl apply -f linkerd-viz.yml
open http://$(kubectl get svc | grep linkerd-viz | awk '{ print $3 }')
```
