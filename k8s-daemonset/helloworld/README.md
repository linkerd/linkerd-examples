# Hello World (legacy) #

NOTE: This is the _old_ helloworld app. You shouldn't need to use it unless
you're using the [`hello-world-legacy.yml`](../k8s/hello-world-legacy.yml)
configuration. For all other use cases, checkout the
[new helloworld app](../../docker/helloworld/).

EXTRA NOTE: Do not move/remove/rename the [`world.txt`](world.txt) file that's
in this directory. It is required for the continuous deployment demo.

This directory contains three Python flask apps: `hello.py`, `world.py` and
`api.py`. The hello app makes an HTTP call to the world app before returning
"Hello World". The api app makes a call to the hello app and returns a json
response. The hello app can be configured to make its RPC call via linkerd using
the [http_proxy](https://linkerd.io/features/http-proxy/) environment variable.

## Packaging ##

To build the apps into a docker image, run:

```
$ docker build -t buoyantio/helloworld .
```

The resulting docker image contains the python files, as well as a
[namerctl](https://github.com/BuoyantIO/namerctl) binary that can be used to
interact with namerd.
