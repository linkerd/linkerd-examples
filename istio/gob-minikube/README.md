# Gob's microservice

Run [Gob](https://github.com/linkerd/linkerd-examples/tree/master/gob)
in minikube with Istio.

## Istio-Envoy
Follow all of the steps in the
[Istio installation guide](https://istio.io/docs/tasks/installing-istio.html)
to install the Istio components.

Then deploy the Gob application:

```
# Add Gob k8s Ingress Resource
kubectl apply -f gob-ingress.yaml

# Deploy Gob
kubectl apply -f <(istioctl kube-inject -f gen.yaml)
kubectl apply -f <(istioctl kube-inject -f word.yaml)
kubectl apply -f <(istioctl kube-inject -f web.yaml)
```

Test that it works:
```
LB=$(minikube ip):$(kubectl get svc istio-ingress -o jsonpath="{.spec.ports[0].nodePort}")
curl http://$LB/gob?limit=10
```

## Istio-Linkerd

Follow the first four steps of the
[Istio installation guide](https://istio.io/docs/tasks/installing-istio.html).

Then deploy Istio with linkerd and the Gob application:

```
# Deploy istio-pilot, istio-mixer and egress
kubectl apply -f ../mixer-pilot.yml -f ../istio-egress.yml

# Add grpc-specific linkerds
kubectl apply -f ../istio-daemonset-grpc.yml

# Add linkerd http ingress controller
kubectl apply -f ../istio-ingress.yml

# Add Gob k8s Ingress Resource
kubectl apply -f gob-ingress.yaml

# Deploy Gob
kubectl apply -f <(linkerd-inject -f word.yaml -useServiceVip)
kubectl apply -f <(linkerd-inject -f gen.yaml -useServiceVip)
kubectl apply -f <(linkerd-inject -f web.yaml -useServiceVip)
```

Test that it works:
```
LB=$(minikube ip):$(kubectl get svc istio-ingress -o jsonpath="{.spec.ports[0].nodePort}")
curl http://$LB/gob?limit=10
```
