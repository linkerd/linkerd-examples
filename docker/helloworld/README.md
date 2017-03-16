# Hello World

This directory contains files for building the "hello world" microservice,
which consists of a "hello" service that calls a "world" service to complete
its request. Configs for running these services in Kubernetes live in the
[`k8s-daemonset/`](../../k8s-daemonset/) directory.

## Building

To build the [buoyantio/helloworld](https://hub.docker.com/r/buoyantio/helloworld/)
Docker image, run:

```bash
$ ./dockerize <tag-name>
```

Where `<tag-name>` is the tag of the image that you want to build.

## Usage

The behavior of each server is controlled via command line flags and environment
variables. The available command line flags are:

```bash
$ helloworld -help
Usage of helloworld:
  -addr string
      address to serve on (default ":7777")
  -failure-rate float
      rate of 500 responses to return
  -json
      return json instead of plaintext responses
  -latency duration
      time to sleep before processing request
  -target string
      target service to call before returning
  -text string
      text to serve (default "Hello")
```

The server also reads a few environment variable that are set as part of our
Kubernetes configs, as follows:

* `POD_IP`---If the service is running in a Kubernetes pod, setting this
  environment to the IP address of the pod will alter the response to include
  the IP after the text string, e.g. `world!` becomes `world (1.2.3.4)!`.

* `TARGET_WORLD`---If set, the value of this environment variable will be used
  as the text string, overriding the value of the `-text` flag, e.g. running
  `TARGET_WORLD=foo helloworld -text bar` will return `foo!`, not `bar!`.

To see a working example of the hello and world services configured to run in
Kubernetes, see [`hello-world.yml`](../../k8s-daemonset/k8s/hello-world.yml).
