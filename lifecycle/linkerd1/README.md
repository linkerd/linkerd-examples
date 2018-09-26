# Linkerd1 lifecycle test configuration

Production testing Linkerd1's discovery & caching.

This test is similar to the [Linkerd2](../) lifecycle environment. Differences
are documented here.

## First time setup

[`lifecycle.yml`](lifecycle.yml) creates a `ClusterRole`, which requires your
user to have this ability.

```bash
kubectl create clusterrolebinding cluster-admin-binding-$USER \
  --clusterrole=cluster-admin --user=$(gcloud config get-value account)
```

## Batch Deploy / Scale / Teardown

Deploy a Linkerd daemonset to the `linkerd` namespace:

```bash
kubectl apply -f linkerd.yml
```

Deploy the lifecycle environment in 3 namespaces:

```bash
bin/deploy 3
```

Scale 3 lifecycle environments to 3 replicas of `bb-broadcast`, `bb-p2p`, and
`bb-terminus` each:

```bash
bin/scale 3 3
```

Teardown 3 lifecycle environments:

```bash
bin/teardown 3
kubectl delete ns linkerd
```

## Observe

Browse to Grafana:

```bash
kubectl -n linkerd port-forward $(
  kubectl -n linkerd get po --selector=name=linkerd-viz -o jsonpath='{.items[*].metadata.name}'
) 3000:3000

open http://localhost:3000
```

Browse to Prometheus:

```bash
kubectl -n linkerd port-forward $(
  kubectl -n linkerd get po --selector=name=linkerd-viz -o jsonpath='{.items[*].metadata.name}'
) 9191:9191

open http://localhost:9191

# view slow-cooker success rates
open "http://localhost:9191/graph?g0.range_input=5m&g0.stacked=0&g0.expr=irate(successes%7Bjob%3D%22slow-cooker%22%7D%5B1m%5D)%20%2F%20irate(requests%7Bjob%3D%22slow-cooker%22%7D%5B30s%5D)&g0.tab=0"
```

Tail slow-cooker logs:

```bash
LIFECYCLE_NS=lifecycle1

kubectl -n $LIFECYCLE_NS logs -f $(
  kubectl -n $LIFECYCLE_NS get po --selector=app=slow-cooker -o name
)
```
