# Hello World

## linkerd-to-linkerd using Kubernetes DaemonSets and linkerd-viz

This is a sample application to demonstrate how to deploy a linkerd-to-linkerd
configuration on Kubernetes using DaemonSets. The application consists of two
python services: hello and world. The hello service calls the world service.

```
hello -> linkerd (outgoing) -> linkerd (incoming) -> world
```

## Building

The Docker image for the hello and world services is [buoyantio/helloword](
https://hub.docker.com/r/buoyantio/helloworld/). You can also build the image
yourself by running:

```bash
cd helloworld
docker build -t <helloworld image name> .
docker push <helloworld image name>
```

## Deploying

### Hello World

Start by deploying the hello and the world apps by running:

```bash
kubectl apply -f k8s/hello-world.yml
```

The default hello-world.yml config requires Kubernetes 1.4 or later. If you're
running an older version of Kubernetes, then you can instead use the legacy
config by running:

```bash
kubectl apply -f k8s/hello-world-legacy.yml
```

Either of these commands will create new Kubernetes Services and
ReplicationControllers for both the hello app and the world app.

### linkerd

Next deploy linkerd. There are three different linkerd deployment configurations
listed below, from simplest to most advanced. Pick which one is best for your
use case.

#### Daemonsets

For the most basic linkerd DaemonSets configuration, you can run:

```bash
kubectl apply -f k8s/linkerd.yml
```

That command will create linkerd's config file as a Kubernetes ConfigMap, and
use the config file to start linkerd as part of a Kubernetes DaemonSet. It will
also add a Kuberentes Service for granting external access to linkerd.

#### Daemonsets + TLS

For a linkerd configuration that adds TLS to all service-to-service calls, run:

```bash
kubectl apply -f k8s/certificates.yml
kubectl apply -f k8s/linkerd-tls.yml
```

Those commands will create the TLS certificates as a Kubernetes Secret, and use
the certificates and the linkerd config to start linkerd as part of a Kubernetes
DaemonSet and encrypt all traffic between linkerd instances using TLS.

#### Daemonsets + namerd + edge routing

To run linkerd and [namerd](https://linkerd.io/in-depth/namerd/) together, with
linkerd running in DaemonSets and serving edge traffic, run:

```bash
kubectl apply -f k8s/namerd.yml
kubectl apply -f k8s/linkerd-namerd.yml
```

Note: namerd stores dtab with the Kubernetes master via the [ThirdPartyResource
APIs](https://kubernetes.io/docs/user-guide/thirdpartyresources/), which
requires a cluster running Kubernetes 1.2+ with the ThirdPartyResource feature
enabled.

Those commands will create the dtab resource, create the namerd config file,
start namerd, create the linkerd config file, and start linkerd, which will use
namerd for routing. linkerd will also run an edge router for handling external
web traffic.

If this is your first time running namerd in your k8s cluster, then you also
need to create the namerd namespaces that are required to run the hello world
app, by running:

```bash
kubectl run namerctl --image=buoyantio/helloworld:0.0.5 --restart=Never -- "./createNs.sh"
```

You can verify the namespaces were created with:

```bash
$ kubectl logs namerctl
Created external
Created internal
```

### linkerd-viz

And lastly, once you have linkerd deployed, you can deploy
[linkerd-viz](https://github.com/BuoyantIO/linkerd-viz) as well, by running:

```bash
curl -s https://raw.githubusercontent.com/BuoyantIO/linkerd-viz/master/k8s/linkerd-viz.yml | kubectl apply -f -
```

## Verifying

### linkerd admin page

View the linkerd admin dashboard:

```bash
L5D_INGRESS_IP=$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
open http://$L5D_INGRESS_IP:9990 # on OS X
```

Note: Kubernetes deploys loadbalancers asynchronously, which means that there
can be a lag time from when a service is created to when its external IP is
available. If the command above fails, make sure that the `EXTERNAL-IP` field in
the output of `kubectl get svc l5d` is not still in the `<pending>` state. If it
is, wait until the external IP is available, and then re-run the command.

If you deployed namerd, visit the namerd admin dashboard:

```bash
NAMERD_INGRESS_IP=$(kubectl get svc namerd -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
open http://$NAMERD_INGRESS_IP:9990 # on OS X
```

### Test Requests

Send some test requests:

```bash
http_proxy=$L5D_INGRESS_IP:4140 curl -s http://hello
http_proxy=$L5D_INGRESS_IP:4140 curl -s http://world
```

If you deployed namerd, then linkerd is also setup to proxy edge requests:

```bash
curl http://$L5D_INGRESS_IP
```

### linkerd-viz dashboard

View the linkerd-viz dashboard:

```bash
L5D_VIZ_IP=$(kubectl get svc linkerd-viz -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
open http://$L5D_VIZ_IP # on OS X
```

## Ingress with linkerd

This folder also contains files demonstrating the use of linkerd (or linkerd in
conjunction with nginx) as an ingress controller.
