# linkerd on Kubernetes #

This doc explains how to get Gob's Microservice up and running.  It describes:
- deploying the application in kubernetes
- staging with per-request routing
- canary & blue-green deployments with traffic management
- fallback

We assume you already have some familiarity with
[linkerd](https://linkerd.io/doc/introduction/) and
[kubernetes](http://kubernetes.io/docs/user-guide/).

This example works in the default namespace of any kubernetes cluster.
If you don't have a cluster at your disposal, though, it's easy (and free)
to try out on [Google Container Engine](./GKE.md).

Before we start running's Gob's app, let's start by setting up _namerd_.

### `namerd` ###

The _namerd_ service isn't required to run _linkerd_, but it allows us to
manage routing rules centrally so that they may be updated without redeploying
linkerd.

Before deploying namerd, install the `namerctl` utility which lets us
interact with the namerd API from the commandline. If you have
[go](https://golang.org) configured on your system, this is as simple
as:

```
:; go get github.com/linkerd/namerctl
:; go install github.com/linkerd/namerctl
:; namerctl -h
Find more information at https://linkerd.io

Usage:
  namerctl [command]

Available Commands:
  dtab        Control namer's dtab interface

Flags:
      --base-url string   namer location (e.g. http://namerd.example.com:4080) (default "http://104.197.215.51")

Use "namerctl [command] --help" for more information about a command.
```

The [`k8s/namerd`](./namerd) directory contains all of the
configuration files needed to boot namerd:

- [`k8s/namerd/config.yml`](./namerd/config.yml) describes a ConfigMap
  volume containing a namerd configuration.
- [`k8s/namerd/rc.yml`](./namerd/rc.yml) describes a
  ReplicationController that runs a single namerd instance.
- [`k8s/namerd/svc-ctl.yml`](./namerd/svc-ctl.yml) describes an
  external Service (with an external IP) so that `namerctl` may
  configure namerd.
- [`k8s/namerd/svc-sync.yml`](./namerd/svc-sync.yml) describes an
  internal Service (with an internal IP) so that _linkerd_ may receive
  configuration updates from namerd.

Running the following will create all of these objects in our
Kubernetes cluster:
```
:; kubectl apply -f k8s/namerd
configmap "namerd-config" created
replicationcontroller "namerd" created
service "namerd-ctl" created
service "namerd-sync" created
```

A namerd pod should be up and running very quickly:
```
:; kubectl get pods
NAME           READY     STATUS    RESTARTS   AGE
namerd-0j88e   1/1       Running   0          2s
```

It may take some time to provision external IPs, so if we list our
services immediately the EXTERNAL-IP is not set for _namerd-ctl_:
```
:; kubectl get svc
NAME          CLUSTER-IP     EXTERNAL-IP   PORT(S)           AGE
kubernetes    10.3.240.1     <none>        443/TCP           1m
namerd-ctl    10.3.255.61                  80/TCP,9990/TCP   7s
namerd-sync   10.3.252.184   <none>        4100/TCP          7s
```

After about a minute the external IP is available:
```
:; kubectl get svc
NAME          CLUSTER-IP     EXTERNAL-IP      PORT(S)           AGE
kubernetes    10.3.240.1     <none>           443/TCP           2m
namerd-ctl    10.3.255.61    104.197.215.51   80/TCP,9990/TCP   1m
namerd-sync   10.3.252.184   <none>           4100/TCP          1m
```

Once this is available, we can configure namerd with a Delegation
Table (Dtab) that describes our default routing policy:

```
:; export NAMERCTL_BASE_URL=http://104.197.215.51
:; namerctl dtab create default k8s/namerd/default.dtab
:; namerctl dtab get default
# version AAAAAAAAAAI=
/srv                => /#/io.l5d.k8s/default/grpc ;
/srv/proto.GenSvc   => /srv/gen ;
/srv/proto.WordSvc  => /srv/word ;
/grpc               => /srv ;
```

A delegation table describes how named requests,
e.g. _/svc/myService/myMethod_, are routed onto a service discovery
backend.  In kubernetes, the _io.l5d.k8s_ namer is used to resolve
names against the kubernetes Endpoints API (which describes Service
objects).  In this dtab, we discover endpoints in the _default_
kubernetes namespace with the _grpc_ port, and map the name of a protobuf
service to a kubernetes service. The linkerd documentation contains a
richer description of [Dtabs](https://linkerd.io/doc/dtabs/).

Now, we're ready to deploy Gob's service.

### Deploying Gob's Microservice ###

Initially, we create a ConfigMap volume with linkerd's configuration:
```
:; kubectl apply -f k8s/linkerd.yml
configmap "linkerd-config" configured
```

Each of Gob's microservices is deployed as a separate
ReplicaController and Service:
```
:; kubectl apply -f k8s/gen
replicationcontroller "gen" created
service "gen" created
:; kubectl apply -f k8s/word
replicationcontroller "word" created
service "word" created
:; kubectl apply -f k8s/web
replicationcontroller "web" created
service "web" created
```

The _web_ Service is created with an external load balancer so that it
may be reached from the internet.
```
:; kubectl get svc/web
NAME      CLUSTER-IP     EXTERNAL-IP       PORT(S)           AGE
web       10.3.251.167   146.148.102.218   80/TCP,9990/TCP   1m
:; export GOB_HOST=146.148.102.218
```

Now, we're able to curl the service!

```
:; curl "$GOB_HOST"
Gob's web service!

Send me a request like:

  146.148.102.218/gob

You can tell me what to say with:

  146.148.102.218/gob?text=WHAT_TO_SAY&limit=NUMBER
```

We have exposed a linkerd admin page for our _web_ frontend as a service!
```
:; open "http://$GOB_HOST:9990"
```

We can curl the site, and it works, using all 3 of gob's services:
```
:; curl -s "$GOB_HOST/gob?text=gob&limit=10"
gob gob gob gob gob gob gob gob gob gob
```

### Staging a new version of _gen_ with linkerd and namerd ###

How do we deploy a new version of a service?  We could replace the
running version with a new one, but this makes it hard for us to get
confidence in the new version before exposing it to users.  In order
to get this sort of confidence, we want to shift _traffic_ from the
old version to the new version. This approach affords the operator
much greater control, and allows instantaneous roll-back.

So, I branch and fixup my code:
```
diff --git a/gob/src/gen/main.go b/gob/src/gen/main.go
index 3b9b762..13c6cd6 100644
--- a/gob/src/gen/main.go
+++ b/gob/src/gen/main.go
@@ -14,11 +14,12 @@ import (
 type genSvc struct{}

 func (s *genSvc) Gen(req *pb.GenRequest, stream pb.GenSvc_GenServer) error {
-       if err := stream.Send(&pb.GenResponse{req.Text}); err != nil {
+       line := req.Text + " <3 k8s\n"
+       if err := stream.Send(&pb.GenResponse{line}); err != nil {
                return err
        }
        doWrite := func() bool {
-               err := stream.Send(&pb.GenResponse{" " + req.Text})
+               err := stream.Send(&pb.GenResponse{line})
                return err == nil
        }
        if req.Limit == 0 {
```

_A docker image with these changes is already published to
[gobsvc/gob:0.8.6-growthhack](https://hub.docker.com/r/gobsvc/gob/tags/)._

#### Staging ####

First, we'll stage the service so that it's not yet receiving
production traffic.

We'll deploy the _gen-growthhack_ pod and service (so that it can be
distinguished from the prior version):

```
:; kubectl apply -f k8s/gen-growthhack
replicationcontroller "gen-growthhack" created
service "gen-growthhack" created
```

This doesn't change what users see---they still see the prior version:
```
:; curl -s "$GOB_HOST/gob?text=gob&limit=10"
gob gob gob gob gob gob gob gob gob gob
```

We can test out our staged service (without altering the web service
at all), by adding a delegation (routing rule) to a request.  For example:

```
:; curl -H 'l5d-dtab: /srv/gen => /srv/gen-growthhack' "$GOB_HOST/gob?text=gob&limit=10"
gob <3 k8s
gob <3 k8s
gob <3 k8s
gob <3 k8s
gob <3 k8s
gob <3 k8s
gob <3 k8s
gob <3 k8s
gob <3 k8s
gob <3 k8s
```

This override says "whenever you refer to the _gen_ service, use
_gen-growthhack_ instead."

With this mechanism, we can enable _per-request_ staging of services
in production!

#### Canary ####

For the sake of the demo, we need to generate some load on the site.
To do so, [slow_cooker](https://github.com/BuoyantIO/slow_cooker) can easily
be launched as a one-off task with:

```
:; kubectl run --image=buoyantio/slow_cooker:1.1.0 slow-cooker -- "-qps 200" "http://$GOB_HOST/gob?limit=10&text=buoyant"
```

Going to the admin page (at `$GOB_HOST:9990/`), we should see
a few hundred requests per second on the site.

Now we can update our dtab to send a controlled 1/20th of requests to
canary the new service:

```
:; cat k8s/namerd/default.dtab
/srv                => /#/io.l5d.k8s/default/grpc;
/srv/proto.GenSvc   => /srv/gen;
/srv/proto.WordSvc  => /srv/word;
/grpc               => /srv;
/srv/proto.GenSvc   => 1 * /srv/gen-growthhack & 19 * /srv/gen;
```

```
:; namerctl dtab update default k8s/namerd/default.dtab
```

And we see a small 5% of traffic go to the new service.  If we don't
like what we see -- latency is too high, success rate drops, users
complain, we can simply remove the growthhack service from the Dtab
and traffic will go back to the original version.

#### Blue-green deploy ####

If we're happy with the canary's performance, we can slowly shift more
traffic onto the new service by updating the Dtab.

For example, the following will send 20% of requests through the new service

```
/srv/proto.GenSvc   => 1 * /srv/gen-growthhack & 4 * /srv/gen;
```
```
:; namerctl dtab update default k8s/namerd/default.dtab
```

Next, we bring it to equal 50% with:
```
/srv/proto.GenSvc   => /srv/gen-growthhack & /srv/gen;
```

When we have a sufficient confidence in the new service, we can give
it 100% of the traffic.  Furthermore, we can enable fallback to the
original service, so we still have a safety net:

```
/srv/proto.GenSvc   => /srv/gen-growthhack | /srv/gen;
```
```
:; namerctl dtab update default k8s/namerd/default.dtab
```

Should the _gen-growthhack_ service suddenly disappear--because we
want to roll-back or due to operator error--we can seamlessly revert
to the old version.

You can simulate this with:

```
:; kubectl delete svc/gen-growthhack
service "gen-growthhack" deleted
```

## Summary ##

So, we find that Delegations give us a uniform tool to instrument:

- Canaries
- Staging
- Blue-green deploys
- Migrations
- Failover
- etc

