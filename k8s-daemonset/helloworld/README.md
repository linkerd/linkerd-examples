# Hello World #

This directory contains two Python flask apps: `hello.py` and `world.py`. The
hello app makes an HTTP call to the world app before returning "Hello World".
The hello app can be configured to make its RPC call via linkerd using the
[http_proxy]( https://linkerd.io/features/http-proxy/) environment variable.

## Packaging ##

To build the apps into a docker image, run:

```
$ docker build -t buoyantio/helloworld .
```

The resulting docker image contains both of the python files, as well as a
[namerctl](https://github.com/BuoyantIO/namerctl) binary that can be used to
interact with namerd.
