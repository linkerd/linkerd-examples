# Hello World

This directory contains files for building the "hello world" microservice,
which consists of a "hello" service that calls a "world" service to complete
its request. Configs for running these services in Kubernetes live in the
[`k8s-daemonset/`](../../k8s-daemonset/) directory.

## Building

To build the [buoyantio/helloworld](https://hub.docker.com/r/buoyantio/helloworld/)
Docker image, run:

```bash
$ docker build -t buoyantio/helloworld:<tag-name> .
```

Where `<tag-name>` is the tag of the image that you want to build.

To regenerate the protobuf gRPC bindings in the proto directory, run:

```bash
$ protoc -I ./proto/ ./proto/helloworld.proto --go_out=plugins=grpc:proto
```

## Usage

The behavior of each server is controlled via command line flags and environment
variables. The available command line flags are:

```bash
$ helloworld -help
Usage of helloworld:
  -addr string
      address to serve on (default ":7777")
  -failure-rate float
      rate of error responses to return
  -json
      return JSON instead of plaintext responses (HTTP only)
  -latency duration
      time to sleep before processing request
  -protocol string
      API protocol: http or grpc (default "http")
  -target string
      target service to call before returning
  -text string
      text to serve (default "Hello")
```

The server also reads a few environment variable that are set as part of our
Kubernetes configs, as follows:

* `POD_IP`: If the service is running in a Kubernetes pod, setting this
  environment to the IP address of the pod will alter the response to include
  the IP after the text string, e.g. `world!` becomes `world (1.2.3.4)!`.

* `TARGET_WORLD`: If set, the value of this environment variable will be used
  as the text string, overriding the value of the `-text` flag, e.g. running
  `TARGET_WORLD=foo helloworld -text bar` will return `foo!`, not `bar!`.

To see a working example of the hello and world services configured to run in
Kubernetes, see [`hello-world.yml`](../../k8s-daemonset/k8s/hello-world.yml).

## Running locally

Here are brief examples to demonstrate how to run the hello world app locally.

### HTTP

Follow the instructions below to bring up the hello service and the world
service, communicating with each other via HTTP.

Start by starting the "world" service on port 7778:

```bash
$ go run main.go -addr :7778 -text world &
starting HTTP server on :7778
```

Next bring up the "hello" service on port 7777, and configure it to make an
additional call to the "world" service running on port 7778:

```bash
$ go run main.go -addr :7777 -text Hello -target localhost:7778 &
starting HTTP server on :7777
```

Send traffic to the "hello" service with:

```bash
$ curl localhost:7777
Hello world!!
```

### gRPC

Follow the instructions below to bring up the hello service and the world
service, communicating with each other via gRPC.

Start by starting the "world" service on port 7778:

```bash
$ go run main.go -addr :7778 -text world -protocol grpc &
starting gRPC server on :7778
```

Next bring up the "hello" service on port 7777, and configure it to make an
additional call to the "world" service running on port 7778:

```bash
$ go run main.go -addr :7777 -text Hello -target localhost:7778 -protocol grpc &
starting gRPC server on :7777
```

Send a unary gRPC request to the "hello" service with:

```bash
$ go run helloworld-client/main.go localhost:7777
Hello world!!
```

Or send a streaming gRPC request with:

```bash
$ go run helloworld-client/main.go -streaming localhost:7777
Hello world!!
Hello world!!
Hello world!!
Hello world!!
Hello world!!
```

Hello friend, nice to meet you.
