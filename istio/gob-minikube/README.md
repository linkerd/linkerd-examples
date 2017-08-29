# Gob's microservice

Run [Gob](https://github.com/linkerd/linkerd-examples/tree/master/gob)
in minikube with Istio.

## Istio-Envoy
Follow all of the steps in the
[istio installation guide](https://istio.io/docs/tasks/installing-istio.html)
to install the istio components.

```
# Add Ingress Resource
kubectl apply -f gob-ingress.yaml

# Deploy Gob
kubectl apply -f <(istioctl kube-inject -f gen.yaml)
kubectl apply -f <(istioctl kube-inject -f word.yaml)
kubectl apply -f <(istioctl kube-inject -f web.yaml)
```

Test that it works:
```
LB=$(minikube ip):$(kubectl get svc istio-ingress -o jsonpath="{.spec.ports[0].nodePort}")
open http://$LB/gob
```

## Istio-Linkerd

Assumes you've gone through the first four steps of the
[Istio installation guide](https://istio.io/docs/tasks/installing-istio.html).

```
# Deploy istio-pilot, istio-mixer and egress
kubectl apply -f ../mixer-pilot.yml -f ../istio-egress.yml

# Add grpc-specific linkerds
kubectl apply -f ../istio-daemonset-grpc.yml

# Add http ingress linkerd
kubectl apply -f ../istio-ingress.yml

# Add Ingress Resource
kubectl apply -f gob-ingress.yaml

# Deploy Gob
kubectl apply -f <(linkerd-inject -f word.yaml -useServiceVip)
kubectl apply -f <(linkerd-inject -f gen.yaml -useServiceVip)
kubectl apply -f <(linkerd-inject -f web.yaml -useServiceVip)
```
