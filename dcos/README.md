# Example configs for DC/OS

For more information see our
[DC/OS Getting Started Guide](https://linkerd.io/getting-started/dcos/).

## Deploy webapp

```bash
dcos marathon app add webapp.json
```

## Deploy linkerd

Note the `linkerd-dcos.json` files assume 4 nodes. Modify this to equal the
total number of public+private nodes in your cluster.

Note that the [simple-proxy/linkerd-dcos.json](simple-proxy/linkerd-dcos.json)
is one of many ways you can setup Linkerd to run in DC/OS. Unfortunately using
`linkerd-dcos.json` might make it difficult to send operating system signals to
Linkerd e.g. `SIGTERM` for graceful shutdown. An
[alternative](simple-proxy/linkerd-dcos-with-fetch.json) is one way to set up
Linkerd so that it can catch os signals. This application definition uses the
`fetch` API available in DC/OS 1.10. You can use a top-level `uri`
 list for DC/OS <= 1.9.

Multiple linkerd configurations are described below. Pick the one that's most
appropriate for your setup. When testing configurations, be sure to set the
`PUBLIC_NODE` env variable to the external address of the public node in your
cluster.

### linkerd simple proxy

To deploy the most basic configuration, with linkerd as a proxy running on port
4140 for inbound requests, run:

```bash
dcos marathon app add simple-proxy/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### linkerd ingress configuration

To deploy linkerd with an ingress router running on port 4142 and an internal
router running on port 4140, run:

```bash
dcos marathon app add ingress/linkerd-dcos.json
```

Test this configuration with:

```bash
$ curl $PUBLIC_NODE:4242/hello
Hello world
```

### linkerd in linker-to-linker mode

To deploy linkerd in linker-to-linker mode, with outgoing traffic served on a
router running on port 4140, and incoming traffic served on a router running on
port 4141, run:

```bash
dcos marathon app add linker-to-linker/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### linkerd with Strict Mode and Marathon Authentication

DC/OS Supports increased security modes. Specifically, with Strict mode, all
Marathon access is via TLS. Additionally, Mesosphere Enterprise DC/OS defaults
to requiring authenticated access to Marathon. This example demonstrates
configuring a linkerd's Marathon Namer for both Strict mode and Marathon
Authentication.

#### Strict mode

To enable linkerd's Marathon namer for Strict mode,
[`linkerd-marathon-auth/linkerd-config.yml`](linkerd-marathon-auth/linkerd-config.yml)
includes the following `io.l5d.marathon` config block:

```yaml
namers:
- kind: io.l5d.marathon
  host: leader.mesos
  port: 443
  prefix: "/io.l5d.marathon"
  uriPrefix: "/marathon"
  tls:
    disableValidation: false
    commonName: master.mesos
    trustCerts:
      - /mnt/mesos/sandbox/.ssl/ca.crt
```

Browse to the [DC/OS Security page](https://docs.mesosphere.com/1.9/security/)
for more information on Strict mode.

#### Marathon Authentication

To configure linkerd to make authenticated requests to Marathon, create a
keypair, service account, and DC/OS secret. Full instructions are documented in
the [DC/OS examples repo](https://github.com/dcos/examples/tree/master/linkerd/1.9#mesosphere-enterprise-dcos).
Follow those instructions, stop at the `Install linkerd` step, as we're going to
install our own linkerd rather than use the DC/OS Universe package.

We now specify the location of these credentials in
[`linkerd-marathon-auth/linkerd-dcos.json`](linkerd-marathon-auth/linkerd-dcos.json):

```json
"env": {
  "DCOS_SERVICE_ACCOUNT_CREDENTIAL": { "secret": "serviceCredential" }
},
"secrets": {
  "serviceCredential": {
    "source": "linkerd-secret"
  }
},
```

#### Deploying

With linkerd configured for Strict mode and Marathon Authentication, we're ready
to deploy:

```bash
dcos marathon app add linkerd-marathon-auth/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### linkerd with namerd

#### namerd

Start by deploying namerd:

```bash
dcos marathon app add namerd/namerd-dcos.json
```

Test the namerd configuration with:

```bash
$ curl $PUBLIC_NODE:4180/api/1/dtabs/default
[{"prefix":"/marathonId","dst":"/#/io.l5d.marathon"},{"prefix":"/svc","dst":"/$/io.buoyant.http.domainToPathPfx/marathonId"}]
```

#### linkerd

Next deploy linkerd configured to talk to namerd when routing requests:

```bash
dcos marathon app add linkerd-with-namerd/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### linkerd with namerd in linker-to-linker mode

Deploy namerd as described in the previous section. Then deploy linkerd in
linker-to-linker mode, configured to talk to namerd when routing requests:

```bash
dcos marathon app add linker-to-linker-with-namerd/linkerd-dcos.json
```

Test this configuration with:

```bash
$ http_proxy=$PUBLIC_NODE:4140 curl -s http://webapp/hello
Hello world
```

### Application Groups

Marathon supports an "Application Group" concept, where applications are
deployed and named using a hierarchical path-based naming structure. Because the
linkerd config examples documented here all use the `domainToPathPfx` rewriting
namer, marathon applications within a group are routed by reversing the group
name into a domain-like name. For example, `webgroup/webapp-a/webapp-a1` becomes `webapp-a1.webapp-a.webgroup`:

#### Webgroup

This example demonstrates linkerd routing requests to a Marathon app in an application group.

```bash
dcos marathon group add webgroup.json
```

```bash
http_proxy=$PUBLIC_NODE:4140 curl webapp-a1.webapp-a.webgroup/hello
Hello world
```

#### Hello World

This example demonstrates inter-service routing, along with a routing override.

Deploy 3 services: `hello`, `world-v1`, `world-v2`:

```bash
dcos marathon group add hello-world.json
```

Route requests `linkerd` -> `hello` -> `linkerd` -> `world-v1`:

```bash
http_proxy=$PUBLIC_NODE:4140 curl hello.hw.buoyant
Hello (10.0.3.80) world (10.0.1.148)!
```

Routing override from `world-v1` to `world-v2`:

```bash
# 25% to world-v2
http_proxy=$PUBLIC_NODE:4140 curl -H 'l5d-dtab: /svc/world-v1.hw.buoyant => 3 * /marathonId/buoyant/hw/world-v1 & /marathonId/buoyant/hw/world-v2' hello.hw.buoyant
Hello (10.0.1.56) world (10.0.1.56)!!

# 75% to world-v2
http_proxy=$PUBLIC_NODE:4140 curl -H 'l5d-dtab: /svc/world-v1.hw.buoyant => /marathonId/buoyant/hw/world-v1 & 3 * /marathonId/buoyant/hw/world-v2' hello.hw.buoyant
Hello (10.0.1.56) earth (10.0.1.56)!!

# 100% to world-v2
http_proxy=$PUBLIC_NODE:4140 curl -H 'l5d-dtab: /svc/world-v1.hw.buoyant => /svc/world-v2.hw.buoyant' hello.hw.buoyant
Hello (10.0.1.56) earth (10.0.1.56)!!
```

### Istio

This example demonstrates deploying Istio on Kubernetes on DC/OS on AWS. For
more details have a look at that example's [`README.md`](istio/README.md).
