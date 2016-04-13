# Linkerd Examples #

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

_word_ and _gen_ implement RPC-ish interfaces with HTTP and JSON.

All three services are implemented in Go with no shared code.  They
may be built and run independently.

If you want to run these programs locally, you'll need to install
[Go 1.6 or later](https://golang.org/dl).  However, we have already
published Docker images to [https://hub.docker.com/u/gobsvc/], so all
you'll really need to do is to set up a client for
[dc/os](./dcos/README.md) or [kubernetes](./k8s/README.md).
