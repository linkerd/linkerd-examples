# Gob's Web Service on DC/OS

This documents the steps required to build and deploy Gob on DC/OS,
integrated with a custom linkerd DC/OS marathon app.

## Linkerd/Namerd Deploy

### Configuration

`dcos/linkerd-dcos-gob.yaml` is the linkerd runtime configuration, customized for this Gob example:

- linkerd routes all http traffic via port `4140` and all h2 traffic via port `4142`. Services perform RPC calls via localhost:4142.
- An internet-facing load balancer should route traffic to localhost:4140 on the DC/OS public node, from there linkerd will route to `web`.

`dcos/marathon/linkerd.json` is the the marathon configuration for running a
linkerd on each node in the cluster. Make sure you update "instances" to reflect
the total number of public and private nodes in your cluster.

### Installation

Install them via the command line with:

```bash
dcos marathon app add dcos/marathon/linkerd.json
dcos marathon app add dcos/marathon/namerd.json
```

Learn more about deploying a custom linkerd on our (DC/OS getting started guide.)[https://linkerd.io/getting-started/dcos/]

## Gob Deploy

```bash
dcos marathon app add dcos/marathon/gensvc.json
dcos marathon app add dcos/marathon/wordsvc.json
dcos marathon app add dcos/marathon/websvc.json
```

## Confirm it all works

### Set up namerctl

`namerctl` is utility for interacting with the namerd API from the commandline

```
:; go install github.com/linkerd/namerctl
:; namerctl -h
```

For the sake of this demonstration, assume $PUBLIC_URL is set to a
public-facing DC/OS node and namerd's API is serving on port 4180 of the
public-facing DC/OS node.

```bash
export PUBLIC_URL=http://example.com
export NAMERCTL_BASE_URL=$PUBLIC_URL:4180
```

### View namerd dtabs


```bash
namerctl dtab get default
```

### Update namerd dtabs

```bash
namerctl dtab update default dcos/namerd.dtab
```

### Test Gob app

```bash
curl "$PUBLIC_URL:4140/gob?text=foo&limit=2"
```
