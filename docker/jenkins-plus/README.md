# jenkins-plus

This directory contains files for building an augmented version of the jenkins
Docker container that is suitable for use in Kubernetes, with namerctl and
kubectl pre-installed.

To build the image, you first have to copy a current kubectl binary to this
directory so that it can be included in the image. Pull it out of the latest
buoyantio/kubectl docker image:

```bash
$ docker run --name kubectl -d buoyantio/kubectl:v1.8.5
$ docker cp kubectl:/kubectl .
$ docker rm -vf kubectl
```

Then build the bouyantio/jenkins-plus image:

```bash
$ docker build -t buoyantio/jenkins-plus:2.60.1 .
```

The version number should match the version number in line 1 of the Dockerfile.
