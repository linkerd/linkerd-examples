# Hello World

## linkerd-to-linkerd using Kubernetes DaemonSets and linkerd-viz

This is a sample application to demonstrate how to deploy a linkerd-to-linkerd
configuration on Kubernetes using DaemonSets. The application consists of two
go services: hello and world. The hello service calls the world service.

```
hello -> linkerd (outgoing) -> linkerd (incoming) -> world
```

## Building

The Docker image for the hello and world services is [buoyantio/helloworld](
https://hub.docker.com/r/buoyantio/helloworld/). There are instructions for
building the image yourself in the [`docker/helloworld`](../docker/helloworld)
directory.

## Deploying

### Hello World

Start by deploying the hello and the world apps by running:

```bash
kubectl apply -f k8s/hello-world.yml
```

Note that this app configuration does not work on [Minikube](
https://github.com/kubernetes/minikube) or versions of
Kubernetes older than 1.4. If you are on one of those platforms, then you can
instead use the legacy app configuration by running:

```bash
kubectl apply -f k8s/hello-world-legacy.yml
```

Either of these commands will create new Kubernetes Services and
ReplicationControllers for both the hello app and the world app.

More information about running linkerd on Kubernetes prior to 1.4 can be found
on our
[Flavors of Kubernetes wiki page](https://github.com/linkerd/linkerd/wiki/Flavors-of-Kubernetes#minikube).

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

#### Daemonsets + CNI

If you are using CNI such as Calico or Weave, you need a slightly modified
config as described [here](https://github.com/linkerd/linkerd/wiki/Flavors-of-Kubernetes#cnicalicoweave).

To deploy this configuration, you can run:

```bash
kubectl apply -f k8s/linkerd-cni.yml
```

#### Daemonsets CNI + TLS + namerd

If you are using CNI such as Calico or Weave, you need a slightly modified
config as described [here](https://github.com/linkerd/linkerd/wiki/Flavors-of-Kubernetes#cnicalicoweave).

To deploy this configuration, you can run:

```bash
kubectl apply -f k8s/certificates.yml
kubectl apply -f k8s/namerd.yml
kubectl apply -f k8s/linkerd-namerd-cni.yml
```

This configuration enables routing via io.l5d.namerd on port 4140, and
io.l5d.mesh on port 4142:

```bash
# via io.l5d.namerd
http_proxy=$L5D_INGRESS_LB:4140 curl -s http://hello
http_proxy=$L5D_INGRESS_LB:4140 curl -s http://world

# via io.l5d.namerd tls
http_proxy=$L5D_INGRESS_LB:4142 curl -s http://hello
http_proxy=$L5D_INGRESS_LB:4142 curl -s http://world

# via io.l5d.mesh
http_proxy=$L5D_INGRESS_LB:4144 curl -s http://hello
http_proxy=$L5D_INGRESS_LB:4144 curl -s http://world

# via io.l5d.mesh tls
http_proxy=$L5D_INGRESS_LB:4146 curl -s http://hello
http_proxy=$L5D_INGRESS_LB:4146 curl -s http://world
```

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

If using Kubernetes prior to version 1.8
```bash
kubectl apply -f k8s/certificates.yml
kubectl apply -f k8s/namerd-legacy.yml
kubectl apply -f k8s/linkerd-namerd.yml
```


If using Kubernetes version 1.8+
```bash
kubectl apply -f k8s/certificates.yml
kubectl apply -f k8s/namerd.yml
kubectl apply -f k8s/linkerd-namerd.yml
```

Note: namerd stores dtabs with the Kubernetes master via the [CustomResourceDefinitions
APIs](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/), which
requires a cluster running Kubernetes 1.8+. If using a cluster running Kubernetes < 1.8, and ThirdPartyResources are enabled in the cluster, you can use [namerd-legacy.yml](/k8s/namerd-legact.yml) to
store dtabs with ThirdPartyResources since CustomResourceDefinitions are not supported.

Those commands will create the dtab resource, create the namerd config file,
start namerd, create the "external" and "internal" namespaces in namerd, create
the linkerd config file, and start linkerd, which will use namerd for routing.
linkerd will also run an edge router for handling external web traffic.

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

This configuration is covered in more detail in:

* [A Service Mesh for Kubernetes, Part VII: Distributed tracing made easy](https://blog.buoyant.io/2017/03/14/a-service-mesh-for-kubernetes-part-vii-distributed-tracing-made-easy/)

#### Daemonsets + Ingress Controller

To route to external requests using ingress resources, deploy linkerd as an
ingress controller and create an ingress resource:

```bash
kubectl apply -f k8s/linkerd-ingress-controller.yml
kubectl apply -f k8s/hello-world-ingress.yml
```

#### Daemonsets + Per-Service Timeouts

This deployment adds 500ms of artificial latency to the `hello` service and
demonstrates linkerd's ability to configure per-service timeouts.  Deploy it
with:

```bash
kubectl apply -f k8s/hello-world-latency.yml
kubectl apply -f k8s/linkerd-latency.yml
```

#### Daemonsets + Egress

To have linkerd fall back to routing to external services via DNS, use this
configuration:

```bash
kubectl apply -f k8s/linkerd-egress.yml
```
#### Daemonsets + Linkerd2

Running Linkerd 1.x together with Linkerd2 is relatively simple and only
requires a few minor changes. These changes relate to the way both proxies
identify services for traffic routing. Linkerd 1.x relies on dtabs and
service names in the `/svc/service name>` format while Linkerd2
uses Kubernetes's DNS name format with port number e.g.
`<service>.<namespace>.svc.<cluster>.local:7777`.

Linkerd 1.x and Linkerd2 work well when Linkerd 1.x is running as a daemonset
in Kubernetes. This example walks you through how to set this up using
the helloworld sample app.

##### Installing Linkerd 1.x
In your k8s environment, install a Linkerd 1.x daemonset configuration. In this
example, we use the `servicemesh.yml` available in the Linkerd examples repo.

First, create a `linkerd` namespace in your environment.
```bash
kubectl create ns linkerd
```

Then, apply RBAC permissions for Linkerd 1.x's k8s service discovery. We first
make changes to the `linkerd-rbac.yml` file in the k8s-daemonset directory
We change lines 33 and 46 to `linkerd` so that the default account gets RBAC
permissions for the `linkerd` namespace. We then apply the file:

```bash
sed 's/namespace: default/namespace: linkerd/g' linkerd-rbac.yml |
  kubectl -n linkerd apply -f -
```

Next, install Linkerd 1.x servicemesh config to the `linkerd` namespace.

```bash
$ kubectl apply -f k8s/servicemesh.yml
```

Confirm that Linkerd 1.x daemonset pods are up and running.

```bash
$ kubectl get pods -n linkerd
NAME        READY     STATUS    RESTARTS   AGE
l5d-dswzp   2/2       Running   0          14h
```

Next, save the YAML configuration below as `hello-world.yml` and install the
`hello-world` using `kubectl apply -f hello-world.yml`

Notice that we are applying this to the `default` namespace instead of
`linkerd`. We want to seperate our sample app in its own data plane and not
add it to the `linkerd` namespace, which is our control plane.

```yaml
---
apiVersion: v1
kind: ReplicationController
metadata:
  name: hello
spec:
  replicas: 3
  selector:
    app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      dnsPolicy: ClusterFirst
      containers:
      - name: service
        image: buoyantio/helloworld:0.1.6
        env:
        - name: NODE_IP
          valueFrom:
            fieldRef:
              fieldPath: spec.hostIP
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: http_proxy
          value: $(NODE_IP):4140
        args:
        - "-addr=:7777"
        - "-text=Hello"
        - "-target=world"
        ports:
        - name: service
          containerPort: 7777
---
apiVersion: v1
kind: Service
metadata:
  name: hello
spec:
  selector:
    app: hello
  clusterIP: None
  ports:
  - name: http
    port: 7777
---
apiVersion: v1
kind: ReplicationController
metadata:
  name: world-v1
spec:
  replicas: 3
  selector:
    app: world-v1
  template:
    metadata:
      labels:
        app: world-v1
    spec:
      dnsPolicy: ClusterFirst
      containers:
      - name: service
        image: buoyantio/helloworld:0.1.6
        env:
        - name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        - name: TARGET_WORLD
          value: world
        args:
        - "-addr=:7778"
        ports:
        - name: service
          containerPort: 7778
---
apiVersion: v1
kind: Service
metadata:
  name: world-v1
spec:
  selector:
    app: world-v1
  clusterIP: None
  ports:
  - name: http
    port: 7778
```

To make things simple, we've provided a slow cooker app that generates requests
to the `hello-world` endpoint periodically. This will help us see Linkerd 1.x
stats without having to manually generate requests.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: slow-cooker1
spec:
  template:
    metadata:
      name: slow-cooker1
      labels:
        testrun: test1
    spec:
      containers:
      - name: slow-cooker
        image: buoyantio/slow_cooker:1.2.0
        env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: status.hostIP
        - name: http_proxy
          value: http://l5d.linkerd.svc.cluster.local:4140
        command:
        - "/bin/bash"
        args:
        - "-c"
        - |
          slow_cooker \
            -qps 10 \
            -concurrency 1 \
            -interval 30s \
            http://hello
        ports:
        - name: slow-cooker
          containerPort: 9991
      restartPolicy: OnFailure
```

Install slow cooker in the default namespace like so:

```bash
kubectl apply -f slow-cooker.yml
```
Confirm that traffic is routing successfully by accessing the slow cooker
container logs and observe the results. You should see a line similar to:

```bash
                      good/b/f t   good%   min [p50 p95 p99  p999]  max change
2018-07-31T21:37:31Z    300/0/0 300 100% 30s   6 [  8  25  50   81 ]   81
...
```
The "good" column indicates the number of requests that were successful.
We also see that our success rate is 100%. So everything is looking great.

##### Installing Linkerd2

Now that we've confirmed that Linkerd 1.x is working as expected, let's install
Linkerd2. Use the [Getting Started](https://github.com/linkerd/linkerd2#getting-started-with-linkerd2)
guide to quickly get Linkerd2 running in your Kubernetes environment.

After running the `linkerd dashboard` command you should see the default
namespace with 0/7 meshed pods. Clicking the `default` link takes you to
a detailed page showing pods that are installed. Before we inject the
hello world sample app with the Linkerd2 proxy, we need to make a few
modifications to our `hello-world.yml`

In our Linkerd 1.x setup, our hello app sent traffic to `http://world-v1`
using an `http_proxy` environment variable. With Linkerd2, we no longer
need to have this environment variable set due to Linkerd2's transparent
traffic routing. All we have to do in our `hello-world.yml` is to remove the
`http_proxy` env variable and change the `-target` URL to
`world-v1.default.svc.cluster.local:7778` on line 33 of our `hello-world.yml`.
The new `hello` pod section config yaml should look like this:

```YAML
...
apiVersion: v1
kind: ReplicationController
metadata:
  name: hello
spec:
  replicas: 3
  selector:
    app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      dnsPolicy: ClusterFirst
      containers:
      - name: service
        image: buoyantio/helloworld:0.1.6
        args:
        - "-addr=:7777"
        - "-text=Hello"
        - "-target=world-v1.default.svc.cluster.local:7778"
        ports:
        - name: service
          containerPort: 7777
....
# world-v1 configuration
```
Once we've made the necessary changes, we can now inject our apps with the
Linkerd2 proxy.

```bash
kubectl delete -f hello-world.yml
linkerd inject hello-world.yml | kubectl apply -f -
```

If you have the Linkerd2 dashboard up, you should see traffic statistics showing
up for both hello and world-v1 pods.

### linkerd-viz

And lastly, once you have linkerd deployed, you can deploy
[linkerd-viz](https://github.com/linkerd/linkerd-viz) as well, by running:

```bash
kubectl apply -f https://raw.githubusercontent.com/linkerd/linkerd-viz/master/k8s/linkerd-viz.yml
```

## Verifying

Use the commands below to test out the app configurations that you deployed in
the previous section. Note that if you're running on Minikube, the verification
commands are different, and are covered in the next section,
[Verifying Minikube](#verifying-minikube).

More information about running linkerd on Minikube can be found on our
[Flavors of Kubernetes wiki page](https://github.com/linkerd/linkerd/wiki/Flavors-of-Kubernetes#minikube).

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

If you deployed namerd, view the namerd admin dashboard:

```bash
NAMERD_INGRESS_LB=$(kubectl get svc namerd -o jsonpath="{.status.loadBalancer.ingress[0].*}")
open http://$NAMERD_INGRESS_LB:9991 # on OS X
```

### Zipkin

If you deployed Zipkin, load the Zipkin UI:

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

If you configured linkerd for egress, then you can also send requests to external services:

```bash
http_proxy=$L5D_INGRESS_LB:4140 curl -s http://linkerd.io/index.html
http_proxy=$L5D_INGRESS_LB:4140 curl -s http://linkerd.io/index.html:443
```

### linkerd-viz dashboard

View the linkerd-viz dashboard:

```bash
L5D_VIZ_LB=$(kubectl get svc linkerd-viz -o jsonpath="{.status.loadBalancer.ingress[0].*}")
open http://$L5D_VIZ_LB # on OS X
```

## Verifying Minikube

If you're not running on Minikube, see the previous [Verifying](#verifying)
section.

### linkerd admin page

To view the linkerd admin dashboard:

```bash
minikube service l5d --url | tail -n1 | xargs open # on OS X
```

### namerd admin page

If you deployed namerd, view the namerd admin dashboard:

```bash
minikube service namerd --url | tail -n1 | xargs open # on OS X
```

### Zipkin

If you deployed Zipkin, view the Zipkin UI:

```bash
minikube service zipkin
```

### linkerd-viz dashboard

If you deployed linkerd-viz, view the linkerd-viz dashboard:

```bash
minikube service linkerd-viz --url | head -n1 | xargs open # on OS X
```

### Test Requests

Send some test requests:

```bash
L5D_INGRESS_LB=$(minikube service l5d --url | head -n1)
http_proxy=$L5D_INGRESS_LB curl -s http://hello
http_proxy=$L5D_INGRESS_LB curl -s http://world
```

## gRPC

There's also a gRPC version of the hello world example available. Start by
deploying linkerd configured to route gRPC requests:

```bash
kubectl apply -f k8s/linkerd-grpc.yml
```

Then deploy the hello and world services configured to communicate using gRPC:

```bash
kubectl apply -f k8s/hello-world-grpc.yml
```

View the linkerd admin dashboard (may take a few minutes until external IP is
available):

```bash
L5D_INGRESS_LB=$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].*}")
open http://$L5D_INGRESS_LB:9990 # on OS X
```

Send a unary gRPC request using the `helloworld-client` script provided by the
buoyantio/helloworld docker image:

```bash
docker run --rm --entrypoint=helloworld-client buoyantio/helloworld:0.1.6 $L5D_INGRESS_LB:4140
```

Add the `-streaming` flag to send a streaming gRPC request:

```bash
docker run --rm --entrypoint=helloworld-client buoyantio/helloworld:0.1.6 -streaming $L5D_INGRESS_LB:4140
```

## RBAC

As of version 1.6, Kubernetes supports
[Role Based Access Control](http://blog.kubernetes.io/2017/04/rbac-support-in-kubernetes.html)
which allows for dynamic control of access to k8s resources.

If you're trying out these examples in an RBAC enabled cluster, you'll need to
add [RBAC rules](https://kubernetes.io/docs/admin/authorization/rbac/) so that
linkerd can access list services using the kubernetes API (which it needs to do
service discovery).

You can find out if your cluster supports RBAC by running:
```
kubectl api-versions | grep rbac
```

If you do, you'll need to grant linkerd/namerd access:
```
kubectl apply -f k8s/linkerd-rbac.yml
```

Note that this is a beta RBAC config. If your cluster only supports alpha RBAC
(`rbac.authorization.k8s.io/v1alpha1`), you'll need to modify the apiVersion in
this config.

The `linkerd-endpoints-reader` role grants access to the `endpoints` and
`services` resources. The `namerd-dtab-storage` role grants access to the
`dtabs.l5d.io` custom resource definition. These roles are applied to the `default`
`ServiceAccount`. In cases where you don't want to grant these perimissions on
the default service account, you should create a service account for linkerd
(and namerd, if applicable), and use that in the `RoleBinding`.

For example, you can create the following ServiceAccount:

```
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: linkerd-daemonset-account
```

And in `linkerd.yml`, specify this service account in the pod spec
with `serviceAccount: linkerd-daemonset-account`.

And in the RoleBinding change the subject of the binding from `default` to
the service account you've created:

```
subjects:
  - kind: ServiceAccount
    name: linkerd-daemonset-account
```

You may also want to use a `Role` in a specified namespace rather than
applying the rules cluster-wide with `ClusterRole`.

If you're using `namerd.yml`, note that namerd needs access to the
`dtabs.l5d.io` `CustomResourceDefinition`. We've already configured this under the
`namerd-dtab-storage` role. If you've modified the name of this resource, be
sure to update the role.
