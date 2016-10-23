# Getting Started with linkerd and Kubernetes

This guide will walk you through getting linkerd running in Kubernetes and
routing requests to an nginx service.

For more information on linkerd, please visit linkerd.io

## Downloading

Start by cloning this repo.  Make sure that `kubectl` is installed and
configured to talk to your Kubernetes cluster.

```
git clone https://github.com/BuoyantIO/linkerd-examples.git
cd linkerd-examples/getting-started/k8s
kubectl cluster-info
```

## Starting nginx

We create a simple nginx app that simply serves a static file on port 80.
To do this in Kubernetes, we create a replication controller and service.  The
service is what allows linkerd to discover the nginx pods and load balance over
them.

```
kubectl apply -f nginx-rc.yml
kubectl apply -f nginx-svc.yml
```

## Starting linkerd

Before we can launch linkerd itself, we need to create a linker config file.
`linkerd.yml` creates a config file that controls linkerd's behavior and stores
it in Kubernetes as a config map.

```yaml
# The Kubernetes namer (io.l5d.k8s) queries the Kubernetes master API
# for a list of pods with a given name.
namers:
- kind: io.l5d.k8s
  experimental: true
  # kubectl proxy forwards localhost:8001 to the Kubernetes master API
  host: localhost
  port: 8001
routers:
- protocol: http
  # Incoming requests to linkerd with a Host header of "hello" get assigned
  # a name like /http/1.1/GET/hello.  This dtab transforms that into
  # /#/io.l5d.k8s/default/service/hello which indicates that the kubernetes
  # namer should query the API for ports named "service" on pods in the
  # "default" namespace named "hello".  linkerd will then load balance over
  # those pods.
  baseDtab: |
    /http/1.1/* => /#/io.l5d.k8s/default/service
  servers:
  - ip: 0.0.0.0
    port: 4140
```

Create the config map with

```
kubectl apply -f linkerd.yml
```

Now we're ready to launch linkerd into Kubernetes.

```
kubectl apply -f l5d-rc.yml
kubectl apply -f l5d-svc.yml
```

Kuberenets should create an external ip for linkerd which you can view with

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
```

## Admin dashboard

You can view the linkerd admin dashboard at `<linkerd external ip>:9990`.
