# Hello World example

Example app for istio-linkerd in minikube.
Assumes you've gone through the first four steps of the
[Istio installation guide](https://istio.io/docs/tasks/installing-istio.html).

```
# Deploy istio-linkerd
kubectl apply -f ../istio-linkerd.yml

# Add Ingress Resource
kubectl apply -f hello-ingress.yml

# Use linkerd-inject to modify the hello world config, and deploy it:
kubectl apply -f <(linkerd-inject -useServiceVip -f hello-world.yml)

# Verify deployment
curl -v $(minikube ip):$(kubectl get svc istio-ingress -o jsonpath="{.spec.ports[0].nodePort}") -H "Host: hello"
```
