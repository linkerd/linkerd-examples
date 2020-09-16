# Linkerd GitOps

This project contains scripts and instructions to manage
[Linkerd](https://linkerd.io) using
[Argo CD](https://argoproj.github.io/argo-cd/).

The scripts and commands are tested with the following software:

1. [kind](https://kind.sigs.k8s.io/) v0.8.1
1. [Linkerd](https://linkerd.io/) 2.8.1
1. [Argo CD](https://argoproj.github.io/argo-cd/) v1.6.1
1. [cert-manager](https://cert-manager.io) 0.15.0
1. [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets) 0.12.4

## Highlights

* Automate the Linkerd control plane install and upgrade lifecycle using Argo CD
* Incorporate Linkerd auto proxy injection feature into the GitOps workflow to
  auto mesh applications
* Securely store the mTLS trust anchor key/cert with offline encryption and
  runtime auto-decryption using sealed-secrets
* Use cert-manager to manage the mTLS issuer key/cert resources
* Utilize Argo CD [projects](https://argoproj.github.io/argo-cd/user-guide/projects/)
  to manage bootstrap dependencies and limit access to servers, namespaces and
  resources
* Uses Argo CD
  [_app of apps_ pattern](https://argoproj.github.io/argo-cd/operator-manual/cluster-bootstrapping/#app-of-apps-pattern)
  to declaratively manage a group of
  [application](https://argoproj.github.io/argo-cd/operator-manual/declarative-setup/#applications)

The complete guide can be found in http://linkerd.io/2/tasks/gitops/
