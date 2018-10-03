# linkerd on Amazon ECS

This example demonstrates deploying linkerd on an AWS ECS cluster, using Consul
for service discovery.

Full instructions can be found in the
[Getting Started: Running in ECS](https://linkerd.io/getting-started/ecs/)
Guide.

[Amazon ECS](https://aws.amazon.com/ecs/) is a container management service.
This guide will demonstrate routing and monitoring your services using Linkerd
in ECS.

## Overview

This guide will demonstrate setting up Linkerd as a service mesh, Consul for
service discovery, a hello-world sample app, and linkerd-viz for monitoring, all
on a fresh ECS cluster.

The following components make up the system:

* `ECS`: Docker container management. Every ECS instance runs the following
  Docker containers:
  * `linkerd`: proxies requests to `hello-world`
  * `consul-agent`: local service discovery agent
  * [`consul-registrator`](https://github.com/gliderlabs/registrator): bridge
  between Docker and Consul, automatically registers services with consul
* `hello-world`: example ECS task deployed separately from foundational
  `ECS`+`linkerd`+`consul-agent` configuration, composed of `hello`, `world`,
  and `world-v2` services
* [`linkerd-viz`](https://github.com/linkerd/linkerd-viz): ECS task deployed
  separately from foundational `ECS`+`linkerd`+`consul-agent` configuration,
  provides a monitoring dashboard for all service traffic
* `consul-server`: service discovery back-end, runs on a single EC2 instance

## Initial Setup

First, launch a cluster of EC2 instances orchestrated by ECS. You can use a
CloudFormation template to create a VPC, security group, and autoscaling group
for launching EC2 instances configured to connect to ECS. You can download or
clone the [linkerd-examples repo](https://github.com/linkerd/linkerd-examples/tree/master/ecs)
to get the template.

Open the CloudFormation console and deploy the `linkerd-ecs-cluster.yml` template
or if you have the AWS CLI installed and configured run:

```bash
KEY_PAIR=<MY KEY PAIR NAME>

aws cloudformation deploy \
  --stack-name linkerd-ecs-cluster \
  --template-file linkerd-ecs-cluster.yml \
  --parameter-overrides KeyName=$KEY_PAIR \
  --capabilities CAPABILITY_IAM
```

The next thing needed is a Consul server, and the `linkerd`, `consul-agent`,
and `consul-registrator` daemons. These can be deployed using the
`linkerd-daemons.yml` template:

```bash
aws cloudformation deploy \
  --stack-name linkerd-daemons \
  --template-file linkerd-daemons.yml \
  --capabilities CAPABILITY_IAM
```

### Register Task Definitions

Now we have an ECS cluster that is running all the daemons necessary for Linkerd
to function. We can run a Linkerd enabled application in the cluster. First we
need to register some task definitions that describe how to configure and boot
the application.

```bash
aws ecs register-task-definition --cli-input-json file://linkerd-viz-task-definition.json
aws ecs register-task-definition --cli-input-json file://hello-world-task-definition.json
```

### Deploy hello-world

Now that all our foundational services are deployed and the task definitions are
available, we can deploy a sample app. The `hello-world` task is composed of a
`hello` service, a `world` service, and a `world-v2` service. To demonstrate
inter-service communication, we configure the `hello` service to call the `world`
service, via `linkerd`.

```bash
aws ecs run-task --cluster l5d-demo --task-definition hello-world --count 2
```

Note that we have deployed two instances of `hello-world`, which results in two
`hello` containers, two `world` containers, and two `world-v2` containers.

## Test everything worked

First we need to create an SSH tunnel to the cluster. The following commands
will choose one of the EC2 hosts, and forward traffic on three local ports to
three remote ports on the EC2 host:

- Traffic to `localhost:9990` will go to the Linkerd dashboard on the remote
  host
- Traffic to `localhost:8500` will go to the Consul admin dashboard on the
  remote host
- Traffic to `localhost:4140` will go to the Linkerd HTTP proxy on the remote
  host

Note that if one of these four ports is already in use on your local machine
you will either have to stop whatever software if using that port, or you can
eliminate the port collision by replace the local port number with another
unused port number in the SSH tunnel command as well as all subsequent commands.

```bash
# Select an ECS node
ECS_NODE=$( \
  aws ec2 describe-instances \
    --filters Name=instance-state-name,Values=running Name=tag:Name,Values=l5d-demo-ecs \
    --query 'Reservations[0].Instances[0].PublicDnsName' \
    --output text \
)

ssh -i "~/.ssh/$KEY_PAIR.pem" \
    -L 127.0.0.1:4140:$ECS_NODE:4140 \
    -L 127.0.0.1:9990:$ECS_NODE:9990 \
    -L 127.0.0.1:8500:$ECS_NODE:8500 ec2-user@$ECS_NODE -N
```

The SSH tunnel is now launched. As long as it runs you will be able to access
the remote Linkerd HTTP proxy, Linkerd Dashboard, and Consul dashboard as if
they were on your local host:

```bash
# view Linkerd and Consul UIs (osx)
open http://localhost:9990
open http://localhost:8500
```

Lets use the tunnel to send some requests to the `helloworld` service via the
Linkerd HTTP proxy:

```bash
# test routing via Linkerd
http_proxy=localhost:4140 curl hello
Hello (172.31.20.160) World (172.31.19.35)!!
```

You will see these requests reflected in the Linkerd dashboard. The request flow
we just tested:

`curl` -> `linkerd` -> `hello` -> `linkerd` -> `world`

### Test dynamic request routing

As our `hello-world` task also included a `world-v2` service, let's test
per-request routing:

```bash
http_proxy=localhost:4140 curl -H 'l5d-dtab: /svc/world => /svc/world-v2' hello
Hello (172.31.20.160) World-V2 (172.31.19.35)!!
```

By setting the `l5d-dtab` header, we instructed Linkerd to dynamically route all
requests destined for `world` to `world-v2`.

## linkerd-viz

[`linkerd-viz`](https://github.com/linkerd/linkerd-viz) collects and displays
metrics for all `linkerd`'s running in a cluster. Prior to deploying, let's
put some load through our system:

```bash
while true; do http_proxy=localhost:4140 curl -s -o /dev/null hello; done
```

Now deploy a single `linkerd-viz` instance:

```bash
aws ecs run-task --cluster l5d-demo --task-definition linkerd-viz --count 1

# find the ECS node running linkerd-viz
TASK_ID=$(aws ecs list-tasks --cluster l5d-demo --family linkerd-viz --desired-status RUNNING --query taskArns[0] --output text)
CONTAINER_INSTANCE=$(aws ecs describe-tasks --cluster l5d-demo --tasks $TASK_ID --query tasks[0].containerInstanceArn --output text)
INSTANCE_ID=$(aws ecs describe-container-instances --cluster l5d-demo --container-instances $CONTAINER_INSTANCE --query containerInstances[0].ec2InstanceId --output text)
VIZ_NODE=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --query Reservations[*].Instances[0].PublicDnsName --output text)

ssh -i "~/.ssh/$KEY_PAIR.pem" -L 127.0.0.1:3000:$VIZ_NODE:3000 ec2-user@$VIZ_NODE -N

# view linkerd-viz (osx)
open http://localhost:3000
```

## Further reading

For more information about configuring Linkerd, see the
[Linkerd Configuration](https://api.linkerd.io/latest/linkerd) page.

For more information about linkerd-viz, see the
[linkerd-viz GitHub repo](https://github.com/linkerd/linkerd-viz).
