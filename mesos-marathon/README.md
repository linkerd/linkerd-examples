# Example configs for Mesos + Marathon

This directory provides a [`docker-compose.yml`](docker-compose.yml) file that
can be used to run [Mesos](http://mesos.apache.org/) and
[Marathon](https://mesosphere.github.io/marathon/). It also provides two
Marathon config files for running linkerd and a sample web app, once the docker
cluster is initialized.

## Get Docker IP address

The `$DOCKER_IP` environment variable is required by the docker-compose file.
The command to obtain this variable may vary depending on your docker setup.

```bash
export DOCKER_IP=$(docker-machine ip)
```

## Boot Mesos + Marathon

```bash
docker-compose up -d

open http://$DOCKER_IP:5050
open http://$DOCKER_IP:8080
```

## Deploy webapp, linkerd, and linkerd-viz

```bash
# deploy webapp and linkerd
curl -H "Content-type: application/json" -X POST http://$DOCKER_IP:8080/v2/apps -d @webapp.json
curl -H "Content-type: application/json" -X POST http://$DOCKER_IP:8080/v2/apps -d @linkerd-marathon.json
curl -H "Content-type: application/json" -X POST http://$DOCKER_IP:8080/v2/apps -d @linkerd-viz.json

# test linkerd
open http://$DOCKER_IP:9990
```

Note that the [`linkerd-marathon.json`](linkerd-marathon.json) file inlines a
copy of [`linkerd-config.yml`](linkerd-config.yml) for convenience. If you are
interested in inlining a custom version of this config file, have a look at
(Deploying A Custom linkerd)[https://linkerd.io/getting-started/dcos/#deploying-a-custom-linkerd]
on the linkerd docs site.

## Test linkerd + webapp

```bash
http_proxy=$DOCKER_IP:4140 curl webapp/hello
```

## Test linkerd-viz

Generate some traffic:

```bash
while true; do http_proxy=$DOCKER_IP:4140 curl -so /dev/null webapp/hello; done
```

Open linkerd-viz:

```bash
open http://$DOCKER_IP:3000
```
