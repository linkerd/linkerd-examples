# Up & Running with Kubernetes on GKE #

[Google Container Engine (GKE)](https://cloud.google.com/container-engine/)
makes it extremely easy to get up and running with a small Kubernetes cluster.
If you haven't already done so, you can get about 2 months of compute resources
[for *free*](https://console.cloud.google.com/billing/freetrial).

Once you set up a project, you'll have to enable Container Engine for your
account and install the [`gcloud` utility](https://cloud.google.com/sdk/) on
your local system.

In this example, we use the _gobs-project_ GCP project, but you'll have to use
your own in place of it.

Once `gcloud` is configured, provision a small kubernetes cluster:
```
:; gcloud container clusters create --num-nodes=5 gobs-cluster
Creating cluster gobs-cluster...done.
Created [https://container.googleapis.com/v1/projects/gobs-project/zones/us-central1-b/clusters/gobs-cluster].
kubeconfig entry generated for gobs-cluster.
NAME          ZONE           MASTER_VERSION  MASTER_IP      MACHINE_TYPE   NODE_VERSION  NUM_NODES  STATUS
gobs-cluster  us-central1-b  1.2.0           104.154.21.97  n1-standard-1  1.2.0         5          RUNNING
```

In order to operate the cluster, we'll need the `kubectl` utility, which we can
install with:
```
:; gcloud components install kubectl
```

Then fetch adminitrative credentials with:
```
:; gcloud container clusters get-credentials gobs-cluster
Fetching cluster endpoint and auth data.
kubeconfig entry generated for gobs-cluster.
```

`kubectl` can list the GCE hosts in your cluster:
```
:; kubectl get nodes
NAME                                  STATUS    AGE
gke-gobs-cluster-03c0208c-node-4ewy   Ready     11m
gke-gobs-cluster-03c0208c-node-m2g5   Ready     11m
gke-gobs-cluster-03c0208c-node-n6m9   Ready     11m
gke-gobs-cluster-03c0208c-node-vr9c   Ready     11m
gke-gobs-cluster-03c0208c-node-w1k6   Ready     11m
```

And, while you shouldn't have to do so, you can get a terminal on any of these hosts with i.e.
```
:; gcloud compute ssh gke-gobs-cluster-03c0208c-node-4ewy
Updated [https://www.googleapis.com/compute/v1/projects/gobs-project].
Warning: Permanently added '104.197.228.216' (RSA) to the list of known hosts.
Warning: Permanently added '104.197.228.216' (RSA) to the list of known hosts.
Linux gke-gobs-cluster-03c0208c-node-4ewy 3.16.0-4-amd64 #1 SMP Debian 3.16.7-ckt20-1+deb8u4google (2016-01-26) x86_64

Welcome to Kubernetes v1.2.0!

You can find documentation for Kubernetes at:
  http://docs.kubernetes.io/

You can download the build image for this release at:
  https://storage.googleapis.com/kubernetes-release/release/v1.2.0/kubernetes-src.tar.gz

It is based on the Kubernetes source at:
  https://github.com/kubernetes/kubernetes/tree/v1.2.0

For Kubernetes copyright and licensing information, see:
  /usr/local/share/doc/kubernetes/LICENSES

ver@gke-gobs-cluster-03c0208c-node-4ewy:~$ 
```

Since we haven't deployed any applications yet, there are no Pods yet:
```
:; kubectl get pods
:;
```

At this point your kubernetes cluster is ready for prime time! Head
back to the [linkerd on Kubernetes](./README.md).
