# Hello World

This directory contains files for building an the "hello world" microservice,
which consists of a "hello" service that calls a "world" service to complete
it's request.

## Building

To build the [buoyantio/helloworld](https://hub.docker.com/r/buoyantio/helloworld/)
Docker image, run:

```bash
$ ./dockerize <tag-name>
```

Where `<tag-name>` is the tag of the image that you want to build.

## Usage

See all available configuration options with the `-help` flag:

```bash
$ go run main.go -help
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
exit status 2
```

To see a working example in Kubernetes, checkout
[`hello-world.yml`](../../k8s-daemonset/k8s/hello-world.yml).
