# Header Classifier

This is an HTTP response classifier plugin for [linkerd](https://linkerd.io)
to serve as an example of how to build and install linkerd plugins.  This
classifer inspects a response header (named "status" by default) to determine if
the response should be classified as a success or failure and if the request
should be retried.  A value of "success" means success, a value of "retry" means
retryable failure, and any other value means non-retryable failure.

# HelloWorld Identifier

HelloWorld Identifier injects a header with the value set from config, and pass
the updated request to next identifier

## Building

This plugin is built with sbt.  Run sbt from the plugins directory.

```
./sbt headerClassifier/assembly
```

This will produce the plugin jar at
`header-classifier/target/scala-2.12/header-classifier-assembly-0.1-SNAPSHOT.jar`.

## Installing

To install this plugin with linkerd, simply move the plugin jar into linkerd's
plugin directory (`$L5D_HOME/plugins`).  Then add a classifier block to the
router in your linkerd config:

```
routers:
- protocol: http
  dtab: /svc => /$/inet/localhost/8888
  service:
    responseClassifier:
      kind: io.buoyant.headerClassifier
      headerName: status

  identifier:
    - kind: io.buoyant.helloWorldIdentifier
      name: foobar
    - kind: io.l5d.methodAndHost
      httpUriInDst: true

  servers:
  - ip: 0.0.0.0
    port: 4140
```
