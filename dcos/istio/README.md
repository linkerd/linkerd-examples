# Istio on Kubernetes on DC/OS on AWS

This guide assumes you have AWS CLI installed. For more information go here: https://aws.amazon.com/cli/

## Deploy DC/OS

```bash
curl -s -o /tmp/single-master.cloudformation.json \
  https://s3-us-west-2.amazonaws.com/downloads.dcos.io/dcos/stable/commit/e38ab2aa282077c8eb7bf103c6fff7b0f08db1a4/cloudformation/single-master.cloudformation.json

AWS_KEYPAIR=<YOUR AWS KEY PAIR>

aws cloudformation deploy \
  --template-file /tmp/single-master.cloudformation.json \
  --stack-name istio-k8s-dcos-aws \
  --parameter-overrides KeyName=$AWS_KEYPAIR SlaveInstanceCount=7 \
  --capabilities CAPABILITY_IAM

DCOS_URL=$(aws cloudformation describe-stacks --stack-name istio-k8s-dcos-aws | jq -r '.Stacks[0].Outputs[] | select(.OutputKey=="DnsAddress") | .OutputValue')
open https://$DCOS_URL # osx only
# select 'Install CLI' from the DC/OS web ui
```

## Deploy Kubernetes

This section is based on https://github.com/mesosphere/dcos-kubernetes-quickstart

```bash
dcos package install --yes beta-kubernetes

# get public IP of DC/OS master
LB_NAME=$(aws elb describe-load-balancers | jq -r ".LoadBalancerDescriptions[] | select(.DNSName==\"$DCOS_URL\") | .LoadBalancerName")
INSTANCE_ID=$(aws elb describe-instance-health --load-balancer-name $LB_NAME | jq -r .InstanceStates[0].InstanceId)
PUBLIC_IP=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID | jq -r .Reservations[0].Instances[0].PublicIpAddress)

# start a tunnel to kubernetes api server
ssh -N -L 9000:apiserver-insecure.kubernetes.l4lb.thisdcos.directory:9000 core@$PUBLIC_IP

# configure kubectl
kubectl config set-cluster dcos-k8s --server=http://localhost:9000
kubectl config set-context dcos-k8s --cluster=dcos-k8s --namespace=default
kubectl config use-context dcos-k8s

# confirm cluster is operational
kubectl get nodes
kubectl --all-namespaces=true get all

# install kube-dns
kubectl create -f https://raw.githubusercontent.com/mesosphere/dcos-kubernetes-quickstart/master/add-ons/dns/kubedns-cm.yaml
kubectl create -f https://raw.githubusercontent.com/mesosphere/dcos-kubernetes-quickstart/master/add-ons/dns/kubedns-svc.yaml
kubectl create -f https://raw.githubusercontent.com/mesosphere/dcos-kubernetes-quickstart/master/add-ons/dns/kubedns-deployment.yaml

# install kubernetes-dashboard
kubectl create -f https://raw.githubusercontent.com/mesosphere/dcos-kubernetes-quickstart/master/add-ons/dashboard/kubernetes-dashboard.yaml
open http://localhost:9000/ui # osx only
```

## Deploy Istio

Note there is a known issue with the beta-kubernentes DC/OS package not
supporting service accounts. By default the Istio components use service
accounts to find and connect to the Kubernetes API Server. As a workaround,
we'll provide modified versions of the Istio and sample app Kubernetes configs
that include a kubeconfig file, instructing our components how to connect to
Kubernetes.

### `kubeconfig` file, defined in our istio `ConfigMap` object:

```yaml
kubeconfig: |-
  apiVersion: v1
  kind: Config
  preferences: {}
  current-context: dcos-k8s

  clusters:
  - cluster:
      server: http://apiserver-insecure.kubernetes.l4lb.thisdcos.directory:9000
    name: dcos-k8s

  contexts:
  - context:
      cluster: dcos-k8s
      namespace: default
      user: ""
    name: dcos-k8s
```

### Modified Kubernetes config to load a `kubeconfig` file:

```yaml
spec:
  containers:
  - name: proxy
    image: docker.io/istio/proxy_debug:0.1.6
    args: ["proxy", "egress", "-v", "2", "--kubeconfig", "/etc/istio/config/kubeconfig"]
    volumeMounts:
    - name: "istio"
      mountPath: "/etc/istio/config"
      readOnly: true
  volumes:
  - name: istio
    configMap:
      name: istio
```

These instructions are based on the istio installation instructions at
https://istio.io/docs/tasks/installing-istio.html.

```bash
curl -L https://git.io/getIstio | sh -

cd istio-0.1.6

# set up istio service account permissions
kubectl apply -f install/kubernetes/istio-rbac-beta.yaml

# deploy istio, with modified kubeconfig file
kubectl apply -f https://raw.githubusercontent.com/linkerd/linkerd-examples/master/dcos/istio/istio.yaml

# deploy metrics collection
kubectl apply -f install/kubernetes/addons/prometheus.yaml
kubectl apply -f install/kubernetes/addons/grafana.yaml
kubectl apply -f install/kubernetes/addons/servicegraph.yaml

# deploy books sample app

# Note that the Istio instructions use `istioctl kube-inject` to insert init-
# container annotations and proxy containers into this sample config file.
# Because we need to insert our own kubeconfig file, we provide an already-
# injected sample app config, and then apply our additional kubeconfig
# modifications:
kubectl apply -f https://raw.githubusercontent.com/linkerd/linkerd-examples/master/dcos/istio/bookinfo.yaml

# view books app
kubectl port-forward $(kubectl get pod -l istio=ingress -o jsonpath='{.items[0].metadata.name}') 3001:80
open http://localhost:3001/productpage

# put load on books app
while true; do curl -o /dev/null -s -w "%{http_code}\n" http://localhost:3001/productpage; done

# forward port 3000 to grafana
kubectl port-forward $(kubectl get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000

open http://localhost:3000/dashboard/db/istio-dashboard
```

## Misc teardown commands

```bash
kubectl delete -f install/kubernetes/istio.yaml

dcos package uninstall beta-kubernetes --app-id=/kubernetes

aws cloudformation delete-stack --stack-name istio-k8s-dcos-aws
```
