# linkerd examples #

**Goals**

- Illustrate how to use linkerd and namerd in production-like environments
    - [DC/OS](dcos/README.md)
    - [kubernetes](k8s/README.md)
- Have some fun

**Non-goals**

- Implement a useful application
- Provide exemplary Go code

## The Application: Gob's Microservice ##

![gob](https://media.giphy.com/media/qJxFuXXWpkdEI/giphy.gif)

_Gob's program_, from the television show _Arrested Development_, was
a childish, inane program.  So, naturally, we've turned it into a
microservice web application that can be run at scale!

This application consists of several components:

- _web_ -- Gob's frontend -- serves plaintext
- _word_ -- chooses a word for _web_ when one isn't provided
- _gen_ -- given a word and a limit, generates a stream of text

The web service is fairly simple (and entirely plaintext):

```
$ curl -s 'localhost:8080/'
Gob's web service!

Send me a request like:

  localhost:8080/gob

You can tell me what to say with:

  localhost:8080/gob?text=WHAT_TO_SAY&limit=NUMBER
```

_web_ may call both _word_ and _gen_ to satisfy a request.

_word_ and _gen_ implement protobuf over gRPC.

All three services are implemented in Go with no shared code (except
the generated proto file).  They may be built and run independently.

### Running locally ###

If you want to run these programs locally, you'll need to install
[Go 1.6 or later](https://golang.org/dl). Start all three services and send a
request:

```
$ go run src/word/main.go &
$ go run src/gen/main.go &
$ go run src/web/main.go &

$ curl -s localhost:8080/gob?limit=1
banana
```

### Running with docker-compose ###

You can also use the [docker-compose](https://docs.docker.com/compose/) file to
run all three services with Docker. In the docker-compose environment, the
services communicate with each other using linkerd and namerd, configured with
the YAML files from the `config` directory of this project. To build and start
all of the services:

```
$ docker-compose build
$ docker-compose up
```

Send a request to linkerd on port 4141, which will be routed to the web service
(assuming you have the `DOCKER_IP` env variable set to your docker IP):

```
$ curl -s $DOCKER_IP:4141/gob?limit=1
illusion
```

Visit the linkerd admin dashboard at `$DOCKER_IP:9990`, and the namerd admin
dashboard at `$DOCKER_IP:9991`.

### Running remotely ###

This repo also provides configs for running the application in
[dc/os](./dcos/README.md) and [kubernetes](./k8s/README.md). These configs take
advantage of pre-built Docker images published to
[https://hub.docker.com/u/gobsvc/].
