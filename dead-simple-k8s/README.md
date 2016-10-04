# Dead Simple Kubernetes with linkerd

## Deploy linkerd

```bash
kubectl apply -f linkerd.yml
```

### View linkerd admin page

```bash
open http://$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):9990
```

## Deploy Hello World

```bash
kubectl apply -f hello-world.yml
```

## Test Requests

```bash
http_proxy=$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):4140 curl -s http://hello
http_proxy=$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):4140 curl -s http://world
```

## Deploy linkerd-viz

```bash
kubectl apply -f linkerd-viz.yml
open http://$(kubectl get svc linkerd-viz -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
```
