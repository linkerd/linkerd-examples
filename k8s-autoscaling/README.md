
Setup minikube

```bash
git clone https://github.com/kubernetes/minikube
cd minikube
ISO_VERSION=v0.19.0 make

out/minikube start \
  --memory=8192 \
  --cpus=2 \
  --iso-url=https://storage.googleapis.com/minikube/iso/minikube-v0.19.0.iso \
  --kubernetes-version=v1.7.0-alpha.2 \
  --feature-gates=AllAlpha=true \
  --v=5 \
  --alsologtostderr \
  --log_dir=/Users/sig/minilogs \
  --extra-config=kubelet.EnableCustomMetrics=true \
  --extra-config=apiserver.RuntimeConfig=apis/apiregistration.k8s.io/v1alpha1=true \
  --extra-config=apiserver.RuntimeConfig=apis/apiregistration.k8s.io=true \
  --extra-config=apiserver.RuntimeConfig=api/apiregistration.k8s.io/v1alpha1=true \
  --extra-config=apiserver.RuntimeConfig=api/apiregistration.k8s.io=true \
  --extra-config=apiserver.RuntimeConfig=apiregistration.k8s.io/v1alpha1=true \
  --extra-config=apiserver.RuntimeConfig=apiregistration.k8s.io=true \
  --extra-config=apiserver.RuntimeConfig=api/all=true \
  --extra-config=controller-manager.HorizontalPodAutoscalerUseRESTClients=true
minikube addons enable heapster
```

Monitoring

```bash
hack/cluster-monitoring/minikube-deploy
hack/example-service-monitoring/deploy
```

```bash
export CM_API=$(kubectl -n custom-metrics get svc api -o template --template {{.spec.clusterIP}})
curl -sSLk https://${CM_API}/apis/custom-metrics.metrics.k8s.io/v1alpha1/namespaces/default/pods/*/http_requests_total
```

Install linkerd and sample app

```bash
kubectl apply -f ../k8s-daemonset/k8s/linkerd.yml
kubectl apply -f ../k8s-daemonset/k8s/hello-world-legacy.yml

OUTGOING_PORT=$(kubectl get svc l5d -o jsonpath='{.spec.ports[?(@.name=="outgoing")].nodePort}')
L5D_ROUTING=http://$(minikube ip):$OUTGOING_PORT

http_proxy=$L5D_ROUTING curl -s http://hello
> Hello (172.17.0.6) world (172.17.0.10)!!

http_proxy=$L5D_ROUTING curl -s http://world
> world (172.17.0.8)!
```

Test metrics

```bash
ADMIN_PORT=$(kubectl get svc l5d -o jsonpath='{.spec.ports[?(@.name=="admin")].nodePort}')
L5D_METRICS=http://$(minikube ip):$ADMIN_PORT
curl -s $L5D_METRICS/admin/metrics/prometheus
```

Deploy HPA

```bash
kubectl apply -f l5d-hpa.yml
```
