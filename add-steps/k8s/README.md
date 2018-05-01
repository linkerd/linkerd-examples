# Add Steps on Kubernetes

## Deploy

```bash
curl https://raw.githubusercontent.com/linkerd/linkerd-examples/master/add-steps/k8s/add-steps.yml | kubectl apply -f -
```

## View Dashboard

```bash
kubectl -n add-steps port-forward $(kubectl -n add-steps get po --selector=app=grafana -o jsonpath='{.items[*].metadata.name}') 3000:3000
open http://localhost:3000
```

## Appendix

### Linkerd Dashboard

```bash
kubectl -n add-steps port-forward $(kubectl --namespace=add-steps get po --selector=app=l5d -o jsonpath='{.items[*].metadata.name}') 9990:9990
open http://localhost:9990
```

### Building

```bash
docker build ../ -t buoyantio/add-steps-app:v1
docker push buoyantio/add-steps-app:v1
```
