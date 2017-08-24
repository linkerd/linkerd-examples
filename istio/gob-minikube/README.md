# Gob's microservice

Run [Gob](https://github.com/linkerd/linkerd-examples/tree/master/gob)
in minikube with Istio.

## Istio-Envoy
Follow all of the steps in the
[istio installation guide](https://istio.io/docs/tasks/installing-istio.html)
to install the istio components.

```
# Add Ingress Resource
kubectl apply -f istio/gob-minikube/gob-ingress.yaml

# Deploy Gob
kubectl apply -f <(istioctl kube-inject -f istio/gob-minikube/gen.yaml)
kubectl apply -f <(istioctl kube-inject -f istio/gob-minikube/word.yaml)
kubectl apply -f <(istioctl kube-inject -f istio/gob-minikube/web.yaml)
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
kubectl apply -f istio/mixer-pilot.yml -f istio/istio-egress.yml

# Add grpc-specific linkerds
kubectl apply -f istio/helloworld-grpc-minikube/istio-daemonset.yml
kubectl apply -f istio/istio-ingress.yml

# Add Ingress Resource
kubectl apply -f istio/gob-minikube/gob-ingress.yaml

# Deploy Gob
kubectl apply -f <(linkerd-inject -f istio/gob-minikube/word.yaml -useServiceVip)
kubectl apply -f <(linkerd-inject -f istio/gob-minikube/gen.yaml -useServiceVip)
kubectl apply -f <(linkerd-inject -f istio/gob-minikube/web.yaml -useServiceVip)
```
