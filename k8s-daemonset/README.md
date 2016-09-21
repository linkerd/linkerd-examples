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

Start by building the Docker image for the hello world application and pushing
it to a repository that your Kubernetes cluster can access.

```
cd helloworld
docker build -t <helloworld image name> .
docker push <helloworld image name>
```

## Deploying

Deploy the hello and world python services, as well as the linkerd daemonset
to the helloworld Kubernetes namespace.  You will need to substitute the name
of your helloworld Docker image into `hello-rc.yml` and `world-rc.yml`.

```
kubectl create ns helloworld
kubectl --namespace=helloworld apply -f k8s/
```
