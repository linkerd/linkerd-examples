# Getting Started with linkerd and Kubernetes

This guide will walk you through getting linkerd running in Kubernetes and
routing requests to an nginx service.

For more information on linkerd, please visit linkerd.io

## Downloading

Start by cloning this repo.  Make sure that `kubectl` is installed and
configured to talk to your Kubernetes cluster.

```
git clone https://github.com/linkerd/linkerd-examples.git
cd linkerd-examples/getting-started/k8s
kubectl cluster-info
```

## Starting nginx

We create a simple nginx app that simply serves a static file on port 80. To do
this in Kubernetes, we create a replication controller and service, defined in
[nginx.yml](nginx.yml). The service is what allows linkerd to discover the nginx
pods and load balance over them. To create nginx in the default namespace, run:

```
kubectl apply -f nginx.yml
```

## Starting linkerd

linkerd stores its config file in a Kubernetes
[ConfigMap](http://kubernetes.io/docs/user-guide/configmap/). The config map,
replication controller, and service for running linkerd are defined in
[linkerd.yml](linkerd.yml). To create linkerd in the default namespace, run:

```
kubectl apply -f linkerd.yml
```

Kubernetes will create an external ip for linkerd which you can view with:

```
kubectl get svc/linkerd
```

## Send requests

Now, when we send requests to linkerd, it will look for a service with the same
name as the Host header to determine where the request should be routed.  In
this case we set the Host header to `nginx` so that the request is routed to the
nginx service.

```
curl -H "Host: nginx" <linkerd external ip>:4140
Hello, linkerd!
```

## Admin dashboard

You can view the linkerd admin dashboard at `<linkerd external ip>:9990`.
