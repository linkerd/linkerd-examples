# Thanos Demo

This environment demonstrates deploying
[Linkerd2](https://github.com/linkerd/linkerd2) and sample apps across 4 cluster
providers, aggregating metrics into a single
[Thanos](https://github.com/improbable-eng/thanos) Querier, and autoscaling the
sample apps when global latency increases.

These configs assume the following:
- Working Kubernetes clusters in 4 cloud providers:
  - Amazon EKS
  - Azure Kubernetes Service (AKS)
  - DigitalOcean
  - Google Kubernetes Engine (GKE)
- Block storage configured and accessible from each Kubernetes cluster

This demo was given at Kubecon Barcelona in 2019:
- https://sched.co/MPbU
- https://youtu.be/qTxunwzYO0g

## Config files

- `cluster-*.yaml`: `Namespace`, `Service`, and `ConfigMap` objects
  to enable persistent IP addresses and provide access to object stores,
  specific to each cluster provider. For more information, see
  [Object Storage](https://github.com/improbable-eng/thanos/blob/master/docs/storage.md)
  in the [Thanos repo](https://github.com/improbable-eng/thanos).
- `linkerd-install-*.yaml`: modified `linkerd install` configs to support Thanos
  integration.
- [`thanos-querier.yaml`](thanos-querier.yaml): The Thanos Querier and Grafana,
  to aggregate metrics from all clusters.
- [`k8s-prometheus-adapter.yaml`](k8s-prometheus-adapter.yaml):
  [k8s-prometheus-adapter](https://github.com/DirectXMan12/k8s-prometheus-adapter),
  to query metrics from Thanos and serve them to the HorizontalPodAutoscaler.
- [`strest.yaml`](strest.yaml):
  [strest-grpc](https://github.com/BuoyantIO/strest-grpc) sample application.

### Linkerd / Thanos integration

To enable Linkerd integration with Thanos, 4 changes are required to the default
`linkerd install` output. These changes have already been made in the
`linkerd-install-*.yaml` files:

1. In `linkerd-prometheus` Deployment, introduce a `thanos-config` volume. This
   references the `ConfigMap` defined in `cluster-*.yaml`, enabling object
   storage:
    ```yaml
    kind: Deployment
      name: linkerd-prometheus
    spec:
      template:
        spec:
          volumes:
          - configMap:
              name: thanos-config
            name: thanos-config
    ```

2. In `ConfigMap/linkerd-prometheus-config`, introduce an `external_labels`
    field to the Prometheus config file, indicating the cluster:
    ```yaml
    kind: ConfigMap
    metadata:
      name: linkerd-prometheus-config
    data:
      prometheus.yml: |-
        global:
          external_labels:
            cluster: aks # or do, eks, gke
    ```

3. In the `linkerd-prometheus` container, set
    `--storage.tsdb.max-block-duration=2h` and
    `--storage.tsdb.min-block-duration=2h`:
    ```yaml
    kind: Deployment
      name: linkerd-prometheus
    spec:
      template:
        spec:
          containers:
          - args:
            - --storage.tsdb.max-block-duration=2h
            - --storage.tsdb.min-block-duration=2h
    ```

4. In `linkerd-prometheus` pod, introduce `thanos-sidecar` and `thanos-store`:
    ```yaml
    kind: Deployment
      name: linkerd-prometheus
    spec:
      template:
        spec:
          containers:
          - name: thanos-sidecar
            image: improbable/thanos:v0.3.2
            args:
            - sidecar
            - --tsdb.path=/data
            - --prometheus.url=http://localhost:9090
            - --cluster.disable
            - --objstore.config-file=/etc/thanos/bucket.yml
            - --grpc-address=0.0.0.0:10901
            - --http-address=0.0.0.0:10902
            ports:
            - name: http-sidecar
              containerPort: 10902
            - name: grpc
              containerPort: 10901
            volumeMounts:
            - mountPath: /data
              name: data
            - mountPath: /etc/prometheus
              name: prometheus-config
              readOnly: true
            - mountPath: /etc/thanos
              name: thanos-config
              readOnly: true
          - name: thanos-store
            image: improbable/thanos:v0.3.2
            args:
            - store
            - --data-dir=/data
            - --cluster.disable
            - --objstore.config-file=/etc/thanos/bucket.yml
            - --grpc-address=0.0.0.0:10911
            - --http-address=0.0.0.0:10912
            - --index-cache-size=500MB
            - --chunk-pool-size=500MB
            ports:
            - name: store-http
              containerPort: 10912
            - name: store-grpc
              containerPort: 10911
            volumeMounts:
            - mountPath: /etc/thanos
              name: thanos-config
              readOnly: true
    ```

## Install Linkerd + Thanos sidecars

```bash
# define kubectl contexts
export AKS=[AKS CONTEXT]
export DO=[DIGITAL OCEAN CONTEXT]
export EKS=[AMAZON EKS CONTEXT]
export GKE=[GKE CONTEXT]

# namespaces, services, and block storage access
cat cluster-aks.yaml | kubectl --context $AKS apply -f -
cat cluster-do.yaml |  kubectl --context $DO apply -f -
cat cluster-eks.yaml | kubectl --context $EKS apply -f -
cat cluster-gke.yaml | kubectl --context $GKE apply -f -

# linkerd and thanos sidecars
cat linkerd-install-aks.yaml | kubectl --context $AKS apply -f -
cat linkerd-install-do.yaml |  kubectl --context $DO apply -f -
cat linkerd-install-eks.yaml | kubectl --context $EKS apply -f -
cat linkerd-install-gke.yaml | kubectl --context $GKE apply -f -
```

### Validate Linkerd is installed

```bash
for CLUSTER in $AKS $DO $EKS $GKE
do
  printf "\n$CLUSTER\n"
  linkerd --context $CLUSTER version
done
```

```bash
linkerd --context $AKS check
```

## Install sample apps, inject Linkerd

```bash
for CLUSTER in $AKS $DO $EKS $GKE
do
  cat strest.yaml |
    linkerd --context $CLUSTER inject - |
    kubectl --context $CLUSTER apply -f -
done
```

### Validate proxy injected

```bash
for CLUSTER in $AKS $DO $EKS $GKE
do
  printf "\n$CLUSTER\n"
  linkerd --context $CLUSTER -n strest stat deploy
done
```

```bash
linkerd --context $AKS dashboard
```

## Install Thanos Querier + Grafana on AKS

The [`thanos-querier.yaml`](thanos-querier.yaml) deployment reads from all 4
Thanos Sidecars. Addresses from each sidecar must be provided to the Thanos
Querier via command line flags.

### Static Addresses

Obtain static addresses for all 4 sidecars:

```bash
for CLUSTER in AKS DO GKE
do
  echo "$CLUSTER"_ADDRESS=$(
    kubectl --context ${!CLUSTER} -n linkerd get svc/thanos-sidecar \
      -o jsonpath='{.spec.loadBalancerIP}'
    )
done

# EKS Load Balancers are configured slightly differently
echo EKS_ADDRESS=$(
  kubectl --context $EKS \
    -n linkerd get svc/thanos-sidecar \
    -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
  )
```

#### Update `store` values in `thanos-querier.yaml`

Each cluster has a Thanos Sidecar and a Thanos Store. Thanos Sidecar reads from
Prometheus and writes into the object store. Thanos Store reads from the object
store are provides historical data to Thanos Querier. Configure Thanos Querier
to read from both sidecar and store on each cluster:

```yaml
kind: Deployment
metadata:
  name: thanos-querier
spec:
  template:
    spec:
      containers:
      - name: thanos
        args:
        - query
        - --store=[AKS_ADDRESS]:10901 # AKS sidecar
        - --store=[AKS_ADDRESS]:10911 # AKS store
        - --store=[DO_ADDRESS]:10901  # DO sidecar
        - --store=[DO_ADDRESS]:10911  # DO store
        - --store=[EKS_ADDRESS]:10901 # EKS sidecar
        - --store=[EKS_ADDRESS]:10911 # EKS store
        - --store=[GKE_ADDRESS]:10901 # GKE sidecar
        - --store=[GKE_ADDRESS]:10911 # GKE store
```

### Deploy

In this example, we deploy Thanos Querier into AKS:

```bash
cat thanos-querier.yaml | kubectl --context $AKS apply -f -
```

### View Thanos Querier Dashboard

```bash
kubectl --context $AKS -n thanos-demo port-forward svc/thanos-querier 10902
```

### View Thanos Querier Grafana

```bash
kubectl --context $AKS -n thanos-demo port-forward svc/grafana 3000
```

#### Autoscaling

[`k8s-prometheus-adapter.yaml`](k8s-prometheus-adapter.yaml) runs in each
cluster, queries Thanos for global metrics, and provides those metrics to a
HorizontalPodAutoscaler. The adapter must be provided a Thanos IP address, which
was created via a LoadBalancer in [`thanos-querier.yaml`](thanos-querier.yaml).

Obtain static addresses for Thanos Querier:

```bash
kubectl --context $AKS -n thanos-demo get svc/thanos-querier -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

Substitute IP address in
[`k8s-prometheus-adapter.yaml`](k8s-prometheus-adapter.yaml):

```yaml
containers:
- name: custom-metrics-apiserver
  image: directxman12/k8s-prometheus-adapter-amd64:v0.5.0
  args:
  ...
  - --prometheus-url=http://[THANOS_QUERIER_IP]:10902/
```

```bash
for CLUSTER in $AKS $DO $EKS $GKE
do
  cat k8s-prometheus-adapter.yaml | kubectl --context $CLUSTER apply -f -
done

# verify
kubectl --context $AKS get --raw /apis/custom.metrics.k8s.io/v1beta1
kubectl --context $AKS get --raw /apis/custom.metrics.k8s.io/v1beta1/namespaces/strest/deployments/*/response_latency_ms_p99
kubectl --context $AKS -n strest describe hpa/strest
```

##### Scale up one strest-client

```bash
kubectl --context $GKE -n strest scale --replicas=20 deploy/strest-client
```
