# Consul

This directory contains a docker-compose environment that demonstrates
how to use the consul service discovery back-end with linkerd. It deploys the demonstration environment in a
[service-to-linker](https://linkerd.io/in-depth/deployment/#service-to-linker)
configuration

## Overview
The following components make up the system:
* `curl` which acts as our client application
* `linkerd` for proxying requests to our service
* `audit` example service which has a `/health` endpoint
* `consul` as our service discovery back-end
* [`consul-registrator`](https://github.com/gliderlabs/registrator)
to automatically registers services with consul

**System overview**
```
+--------+      +---------+    +-----------------+
| client +----> | linkerd +--> | service (audit) |
+--------+      +----^----+    +-------+---------+
                     |                 |
                +----+---+     +-------v------------+
                | consul <-----+ consul registrator |
                +--------+     +--------------------+
```


## Startup

The [`docker-compose.yaml`](docker-compose.yaml) file that's included
in this directory is configured to run the demo. Start everything with:

```bash
$ docker-compose build && docker-compose up -d
```

## Testing the system
To make sure everything is working properly run the following command:
```bash
$ curl localhost:4140/audit/health
```

You will get the following response:
```bash
I am healthy!
```