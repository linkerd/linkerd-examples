# Enabling Ingress Traffic

This doc details how to run the
[ingress istio example](https://istio.io/docs/tasks/ingress.html) using linkerd
 daemonsets as the istio proxy.

1) Setup Istio by following the instructions in the
[Installation guide](https://istio.io/docs/tasks/installing-istio.html).

2) Deploy the linkerd daemonset:
`kubectl apply -f istio-daemonset.yml`

3) Replace the default istio ingress and egress controllers with linkerd-powered ones:
`kubectl apply -f istio-ingress.yml -f istio-egress.yml`

4) Start the httpbin sample, which will be used as the destination service to
 be exposed externally. From the istio project directory, run:
`kubectl apply -f samples/apps/httpbin/httpbin.yaml`

5) Follow istio instructions for configuring
[HTTP ingress](https://istio.io/docs/tasks/ingress.html#configuring-ingress-http) or
 [HTTPs ingress](https://istio.io/docs/tasks/ingress.html#configuring-secure-ingress-https).

More information on ingress controllers can be found on the
[Buoyant blog](https://blog.buoyant.io/2017/04/06/a-service-mesh-for-kubernetes-part-viii-linkerd-as-an-ingress-controller/).

# Enabling Egress Traffic

1) Setup Istio by following the instructions in the
[Installation guide](https://istio.io/docs/tasks/installing-istio.html).

2) Deploy the linkerd daemonset:
`kubectl apply -f istio-daemonset.yml`

3) Replace the default istio ingress and egress controllers with linkerd-powered ones:
`kubectl apply -f istio-ingress.yml -f istio-egress.yml`

4) Calls to outside the cluster are now possible (ExternalName services not
 yet supported).

# Istio on DC/OS

For an example of running Istio on DC/OS, have a look at
[this example](https://github.com/linkerd/linkerd-examples/tree/master/dcos/istio).
