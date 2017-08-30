# Hello World example

Example app for istio-linkerd in minikube using grpc.
Assumes you've gone through the first four steps of the
[Istio installation guide](https://istio.io/docs/tasks/installing-istio.html).

```
# Deploy istio-pilot, istio-mixer and egress
kubectl apply -f ../mixer-pilot.yml -f ../istio-egress.yml

# Add Ingress Resource
kubectl apply -f hello-ingress.yml

# Add grpc-specific linkerds
kubectl apply -f ../istio-daemonset-grpc.yml -f istio-ingress.yml

# Use linkerd-inject to modify the hello world grpc config, and deploy it:
kubectl apply -f <(linkerd-inject -useServiceVip -f hello-world-grpc.yml)

# Verify deployment
LB=$(minikube ip):$(kubectl get svc istio-ingress -o jsonpath="{.spec.ports[0].nodePort}")
helloworld-client $LB
```
