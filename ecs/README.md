# linkerd on Amazon ECS

This example demonstrates deploying linkerd on an AWS ECS cluster, using Consul
for service discovery.

Full instructions can be found in the
[Getting Started: Running in ECS](https://linkerd.io/getting-started/ecs/)
Guide.

## Initial Setup

```bash
KEY_PAIR=<MY KEY PAIR NAME>
```

## Security Group

```bash
GROUP_ID=$(aws ec2 create-security-group --group-name l5d-demo-sg --description "Linkerd Demo" | jq -r .GroupId)
aws ec2 authorize-security-group-ingress --group-id $GROUP_ID \
  --ip-permissions \
  FromPort=22,IpProtocol=tcp,ToPort=22,IpRanges=[{CidrIp="0.0.0.0/0"}] \
  FromPort=4140,IpProtocol=tcp,ToPort=4140,IpRanges=[{CidrIp="0.0.0.0/0"}] \
  FromPort=9990,IpProtocol=tcp,ToPort=9990,IpRanges=[{CidrIp="0.0.0.0/0"}] \
  FromPort=3000,IpProtocol=tcp,ToPort=3000,IpRanges=[{CidrIp="0.0.0.0/0"}] \
  FromPort=8500,IpProtocol=tcp,ToPort=8500,IpRanges=[{CidrIp="0.0.0.0/0"}] \
  IpProtocol=-1,UserIdGroupPairs=[{GroupId=$GROUP_ID}]
```

## `consul-server`

```bash
aws ec2 run-instances --image-id ami-62e0d802 \
  --instance-type m4.xlarge \
  --user-data file://consul-server-user-data.txt \
  --placement AvailabilityZone=us-west-1a \
  --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=l5d-demo-consul-server}]" \
  --key-name $KEY_PAIR --security-group-ids $GROUP_ID
```

## ECS Cluster

```bash
aws ecs create-cluster --cluster-name l5d-demo
```

### Role Policy

```bash
aws iam put-role-policy --role-name ecsInstanceRole --policy-name l5dDemoPolicy --policy-document file://ecs-role-policy.json
```

### Register Task Definitions

```bash
aws ecs register-task-definition --cli-input-json file://linkerd-task-definition.json
aws ecs register-task-definition --cli-input-json file://linkerd-viz-task-definition.json
aws ecs register-task-definition --cli-input-json file://consul-agent-task-definition.json
aws ecs register-task-definition --cli-input-json file://consul-registrator-task-definition.json
aws ecs register-task-definition --cli-input-json file://hello-world-task-definition.json
```

### Create Launch Configuration

```bash
aws autoscaling create-launch-configuration \
  --launch-configuration-name l5d-demo-lc \
  --image-id ami-62e0d802 \
  --instance-type m4.xlarge \
  --user-data file://ecs-user-data.txt \
  --iam-instance-profile ecsInstanceRole \
  --security-groups $GROUP_ID \
  --key-name $KEY_PAIR
```

### Create Auto Scaling Group

```bash
aws autoscaling create-auto-scaling-group \
  --auto-scaling-group-name l5d-demo-asg \
  --launch-configuration-name l5d-demo-lc \
  --min-size 1 --max-size 3 --desired-capacity 2 \
  --tags ResourceId=l5d-demo-asg,ResourceType=auto-scaling-group,Key=Name,Value=l5d-demo-ecs,PropagateAtLaunch=true \
  --availability-zones us-west-1a
```

### Deploy `hello-world`

```bash
aws ecs run-task --cluster l5d-demo --task-definition hello-world --count 2
```

## Test everything worked

```bash
# Select an ECS node
ECS_NODE=$( \
  aws ec2 describe-instances \
    --filters Name=instance-state-name,Values=running Name=tag:Name,Values=l5d-demo-ecs \
    --query Reservations[*].Instances[0].PublicDnsName --output text \
)

# test routing via linkerd
http_proxy=$ECS_NODE:4140 curl hello
Hello (172.31.20.160) World (172.31.19.35)!!

# test dynamic routing to world-v2
http_proxy=$ECS_NODE:4140 curl -H 'l5d-dtab: /svc/world => /svc/world-v2' hello
Hello (172.31.20.160) World-V2 (172.31.19.35)!!

# view linkerd and Consul UIs (osx)
open http://$ECS_NODE:9990
open http://$ECS_NODE:8500
```

## linkerd-viz

Put some load on our example app

```bash
while true; do http_proxy=$ECS_NODE:4140 curl -s -o /dev/null hello; done
```

Deploy linkerd-viz

```bash
aws ecs run-task --cluster l5d-demo --task-definition linkerd-viz --count 1

# find the ECS node running linkerd-viz
TASK_ID=$(aws ecs list-tasks --cluster l5d-demo --family linkerd-viz --desired-status RUNNING --query taskArns[0] --output text)
CONTAINER_INSTANCE=$(aws ecs describe-tasks --cluster l5d-demo --tasks $TASK_ID --query tasks[0].containerInstanceArn --output text)
INSTANCE_ID=$(aws ecs describe-container-instances --cluster l5d-demo --container-instances $CONTAINER_INSTANCE --query containerInstances[0].ec2InstanceId --output text)
ECS_NODE=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --query Reservations[*].Instances[0].PublicDnsName --output text)

# view linkerd-viz (osx)
open http://$ECS_NODE:3000
```

## Credits

This example setup is based on these excellent blog posts:

  - https://kevinholditch.co.uk/2017/06/28/running-linkerd-in-a-docker-container-on-aws-ecs/
  - https://blog.unif.io/deploying-consul-with-ecs-2c4ca7ab2981
  - https://medium.com/attest-engineering/linkerd-a-service-mesh-for-aws-ecs-937f201f847a

## TODO

- rolling a new linkerd version

## Teardown

```bash
CONSUL_SERVER_INSTANCE=$(aws ec2 describe-instances --filters Name=instance-state-name,Values=running Name=tag:Name,Values=l5d-demo-consul-server --query Reservations[*].Instances[*].[InstanceId] --output text)
aws ec2 terminate-instances --instance-id $CONSUL_SERVER_INSTANCE
aws autoscaling update-auto-scaling-group --auto-scaling-group-name l5d-demo-asg --min-size 0 --max-size 0 --desired-capacity 0
until (aws autoscaling delete-auto-scaling-group --auto-scaling-group-name l5d-demo-asg); do
  echo "waiting to delete auto scaling group"
  sleep 1
done
aws autoscaling delete-launch-configuration --launch-configuration-name l5d-demo-lc
aws ecs delete-cluster --cluster l5d-demo
aws ec2 delete-security-group --group-name l5d-demo-sg
```
