# Hello World

## linkerd-to-linkerd using Kubernetes DaemonSets and linkerd-viz

This is a sample application to demonstrate how to deploy a linkerd-to-linkerd
configuration on Kubernetes using DaemonSets. The application consists of two
python services: hello and world. The hello service calls the world service.

```
hello -> linkerd (outgoing) -> linkerd (incoming) -> world
```

## Building

The Docker image for the hello and world services can be found at
`buoyantio/helloworld:0.0.1`. You can also build the image yourself by running:

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

Alternatively, if you are running on Kubernetes 1.4 or later, you can take
advantage of some newer Kubernetes features, and run:

```bash
kubectl apply -f k8s/hello-world-1_4.yml
```

Either of these commands will create new Kubernetes Services and
ReplicationControllers for both the hello app and the world app.

### linkerd

Next deploy linkerd, to route requests from the hello app to the world app.

For the most basic linkerd DaemonSets configuration, you can run:

```bash
kubectl apply -f k8s/linkerd.yml
```

That command will create linkerd's config file as a Kubernetes ConfigMap, and it
will use that config file to start linkerd as part of a Kubernetes DaemonSet. It
will also add a Kuberentes Service for granting external access to linkerd.

For a linkerd configuration that adds TLS to all service-to-service calls, run:

```bash
kubectl apply -f k8s/linkerd-tls.yml
```

Or finally, to run linkerd with routing resolutions handled by namerd, run:

```bash
kubectl apply -f k8s/namerd.yml
kubectl apply -f k8s/linkerd-namerd.yml
```

That command will start namerd in your cluster, and configure linkerd to use
namerd when applying routing rules.

### linkerd-viz

And lastly, deploy linkerd-viz by running:

```bash
kubectl apply -f k8s/linkerd-viz.yml
```

## Verifying

### linkerd admin page

View the linkerd admin dashboard:

```bash
open http://$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):9990
```

If you deployed namerd, visit the namerd admin dashboard:

```bash
open http://$(kubectl get svc namerd -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):9990
```

### Test Requests

Send some test requests:

```bash
http_proxy=$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):4140 curl -s http://hello
http_proxy=$(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}"):4140 curl -s http://world
```

If you deployed namerd, then linkerd is also setup to proxy edge requests:

```bash
curl $(kubectl get svc l5d -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
```

### linkerd-viz dashboard

View the linkerd-viz dashboard:

```bash
open http://$(kubectl get svc linkerd-viz -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
```
