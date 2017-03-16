# Linkerd as Kubernetes Ingress Controller

This example application demonstrates the basics of using linkerd as an
ingress controller. It consists of:
* `ingress-controller.yaml`: A linkerd config and linkerd deployment listening on `hostPort: 80`.
* `services.yaml`: Deployment and service definitions for a number of seal-related example services.
* `ingress.yaml`: An ingress resource that defines routing to those services.

## Deploy

### Ingress Controller

```bash
$ kubectl apply -f ingress-controller.yaml
```

Verify linkerd pods:

```bash
$ kubectl get po --namespace=l5d-system
NAME                  READY     STATUS    RESTARTS   AGE
l5d-932856984-v906t   2/2       Running   0          34m
```

Open admin ui (minikube specific):

```bash
$ minikube service l5d -n l5d-system --url | tail -n1 | xargs open
```

## Services

```bash
$ kubectl apply -f services.yaml
```

Verify pods:

```bash
$ kubectl get po
NAME                             READY     STATUS    RESTARTS   AGE
elephant-seal-4102539096-4kg5k   1/1       Running   0          4s
fur-seal-3543844180-v6dhv        1/1       Running   0          4s
sea-lion-2828322125-wzl3p        1/1       Running   0          4s
seal-2460861402-7rqg7            1/1       Running   0          5s
walrus-2639512787-nzlmh          1/1       Running   0          5s
```

### Ingress Resource

```bash
$ kubectl apply -f ingress.yaml
```

Verify resource:

```bash
$ kubectl get ingress
NAME       HOSTS                                                                                    ADDRESS   PORTS     AGE
pinniped   carnivora.pinniped.odobenidae,carnivora.pinniped.otariidae,carnivora.pinniped.phocidae             80        1h
```

## Ingress Routing

The pinniped ingress defines a `.spec.backend.serviceName` and
`.spec.backend.servicePort`, which linkerd routes to if the request doesn't
match any of the `.spec.rules`. You can exercise that functionality by
specifying neither host header nor path:

```bash
# minikube specific cmd, but this can be any k8s cluster ip
CLUSTER_IP=$(minikube ip)
curl -v $CLUSTER_IP
```

By specifying a host header, only `.spec.rules` that match that header are used
for routing.

```bash
curl -v --header "Host:carnivora.pinniped.odobenidae" $CLUSTER_IP
```

By specifying both host header and path, you can choose a specific
`.spec.rules.paths` to route to.

```bash
curl -v --header "Host:carnivora.pinniped.phocidae" $CLUSTER_IP/elephant-seal
curl -v --header "Host:carnivora.pinniped.phocidae" $CLUSTER_IP/true-seal
```

And here's more complex `.spec.rules.paths.path` regex matching:

```bash
curl -v --header "Host:carnivora.pinniped.otariidae" $CLUSTER_IP/fur/seal
curl -v --header "Host:carnivora.pinniped.otariidae" $CLUSTER_IP/water-lion
```
