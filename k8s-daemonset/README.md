# Hello World

## linkerd-to-linkerd using Kubernetes daemonsets

This is a sample application to demonstrate how to deploy a linkerd-to-linkerd
configuration on Kubernetes using daemonsets.  The application consists of
two python services: hello and world.  The hello service calls the world
service.

```
hello -> linkerd (outgoing) -> linkerd (incoming) -> world

```

## Building

The Docker image for the hello and world serivces can be found at
`buoyantio/helloworld:latest`.  You can also build the image yourself by running

```
cd helloworld
docker build -t <helloworld image name> .
docker push <helloworld image name>
```

## Deploying

Deploy the hello and world python services, as well as the linkerd daemonset
to the helloworld Kubernetes namespace.

```
kubectl create ns helloworld
kubectl --namespace=helloworld apply -f k8s/
```
