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

Next deploy linkerd. There are multiple different linkerd deployment
configurations listed below, showcasing different linkerd features. Pick which
one is best for your use case.

#### Daemonsets

For the most basic linkerd DaemonSets configuration, you can run:

```bash
kubectl apply -f k8s/linkerd.yml
```

That command will create linkerd's config file as a Kubernetes ConfigMap, and
use the config file to start linkerd as part of a Kubernetes DaemonSet. It will
also add a Kuberentes Service for granting external access to linkerd.

This configuration is covered in more detail in:

* [A Service Mesh for Kubernetes, Part II: Pods Are Great Until They're Not](https://blog.buoyant.io/2016/10/14/a-service-mesh-for-kubernetes-part-ii-pods-are-great-until-theyre-not/)

#### Daemonsets + TLS

For a linkerd configuration that adds TLS to all service-to-service calls, run:

```bash
kubectl apply -f k8s/certificates.yml
kubectl apply -f k8s/linkerd-tls.yml
```

Those commands will create the TLS certificates as a Kubernetes Secret, and use
the certificates and the linkerd config to start linkerd as part of a Kubernetes
DaemonSet and encrypt all traffic between linkerd instances using TLS.

This configuration is covered in more detail in:

* [A Service Mesh for Kubernetes, Part III: Encrypting all the things](https://blog.buoyant.io/2016/10/24/a-service-mesh-for-kubernetes-part-iii-encrypting-all-the-things/)

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
kubectl run namerctl --image=buoyantio/helloworld:0.0.6 --restart=Never -- "./createNs.sh"
```

You can verify the namespaces were created with:

```bash
$ kubectl logs namerctl
Created external
Created internal
```

This configuration is covered in more detail in:

* [A Service Mesh for Kubernetes, Part IV: Continuous deployment via traffic shifting](https://blog.buoyant.io/2016/11/04/a-service-mesh-for-kubernetes-part-iv-continuous-deployment-via-traffic-shifting/)

#### Daemonsets + ingress + NGINX

To deploy a version of linkerd that uses NGINX for edge routing, and routing
traffic to multiple different backends based on domain, run:

```bash
kubectl apply -f k8s/nginx.yml
kubectl apply -f k8s/linkerd-ingress.yml
kubectl apply -f k8s/api.yml
```

Those commands will deploy a version of NGINX that's configured to route traffic
to linkerd. linkerd is configured to route requests on the `www.hello.world`
domain to the hello service, and requests on the `api.hello.world` to the API
service.

This configuration is covered in more detail in:

* [A Service Mesh for Kubernetes, Part V: Dogfood environments, ingress and edge routing](https://blog.buoyant.io/2016/11/18/a-service-mesh-for-kubernetes-part-v-dogfood-environments-ingress-and-edge-routing/)

#### Daemonsets + Zipkin

For a linkerd configuration that exports tracing data to Zipkin, first start
Zipkin, and then start linkerd with a Zipkin tracer configured:

```bash
kubectl apply -f k8s/zipkin.yml
kubectl apply -f k8s/linkerd-zipkin.yml
```

Those commands will start a Zipkin process in your cluster, running with an
in-memory span store and scribe collector, receiving linkerd tracing data via
the `io.l5d.zipkin` tracer.

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
L5D_INGRESS_LB=$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].*}")
open http://$L5D_INGRESS_LB:9990 # on OS X
```

Note: Kubernetes deploys loadbalancers asynchronously, which means that there
can be a lag time from when a service is created to when its external IP is
available. If the command above fails, make sure that the `EXTERNAL-IP` field in
the output of `kubectl get svc l5d` is not still in the `<pending>` state. If it
is, wait until the external IP is available, and then re-run the command.

### namerd admin page

If you deployed namerd, visit the namerd admin dashboard:

```bash
NAMERD_INGRESS_LB=$(kubectl get svc namerd -o jsonpath="{.status.loadBalancer.ingress[0].*}")
open http://$NAMERD_INGRESS_LB:9990 # on OS X
```

### Zipkin

If you deployed zipkin, load the Zipkin UI:

```bash
ZIPKIN_LB=$(kubectl get svc zipkin -o jsonpath="{.status.loadBalancer.ingress[0].*}")
open http://$ZIPKIN_LB # on OS X
```

### Test Requests

Send some test requests:

```bash
http_proxy=$L5D_INGRESS_LB:4140 curl -s http://hello
http_proxy=$L5D_INGRESS_LB:4140 curl -s http://world
```

If you deployed namerd, then linkerd is also setup to proxy edge requests:

```bash
curl http://$L5D_INGRESS_LB
```

If you deployed NGINX, then you can also use that to initiate requests:

```bash
NGINX_LB=$(kubectl get svc nginx -o jsonpath="{.status.loadBalancer.ingress[0].*}")
curl -H 'Host: www.hello.world' http://$NGINX_LB
curl -H 'Host: api.hello.world' http://$NGINX_LB
```

### linkerd-viz dashboard

View the linkerd-viz dashboard:

```bash
L5D_VIZ_LB=$(kubectl get svc linkerd-viz -o jsonpath="{.status.loadBalancer.ingress[0].*}")
open http://$L5D_VIZ_LB # on OS X
```
