#!/bin/sh

LINKERD_PID=/tmp/linkerd.pid
WAIT=10

start_linkerd() {
  echo -e "\nWaiting for linkerd to start: $1"
  $L5D_EXEC $1 &
  echo $! > $LINKERD_PID
  echo -e "Starting linkerd on pid $(cat $LINKERD_PID)"

  for i in $(seq $WAIT); do

    pong=$(curl -s localhost:9990/admin/ping)
    if [ "$pong" = "pong" ]; then
      echo "Startup complete on pid $(cat $LINKERD_PID)"
      return 0
    fi

    printf "."
    sleep 1
  done

  echo "Could not start linkerd on pid $(cat $LINKERD_PID)"
  return 1
}

shutdown_linkerd() {
  echo "Waiting for linkerd to shutdown on pid $(cat $LINKERD_PID)"
  curl -s -X POST localhost:9990/admin/shutdown

  for i in $(seq $WAIT); do
    pong=$(curl -s localhost:9990/admin/ping)
    if [ "$pong" = "pong" ]; then
      printf "."
      sleep 1
    else
      echo "Shutdown complete on pid $(cat $LINKERD_PID)"
      rm -f $LINKERD_PID
      return 0
    fi
  done

  echo "Could not shutdown linkerd, forcing on pid $(cat $LINKERD_PID) and failing"
  kill -9 "$(cat $LINKERD_PID)"
  for i in $(seq $WAIT); do
    if ps -o pid|grep $(cat $LINKERD_PID) > /dev/null; then
      printf "."
      sleep 1
    else
      echo "Forced shutdown complete on pid $(cat $LINKERD_PID)"
      rm -f $LINKERD_PID
      return 1
    fi
  done

  echo "Could not force shutdown linkerd on pid $(cat $LINKERD_PID)"
  return 1
}

test_config() {
  if ! start_linkerd $1 ; then
    echo "Could not start linkerd with $1"
    exit 1
  fi

  if ! shutdown_linkerd ; then
    echo "Could not shutdown linkerd with $1"
    exit 1
  fi
}

cp -a ~/linkerd-examples/add-steps/disco /
test_config ~/linkerd-examples/add-steps/linkerd.yml

test_config ~/linkerd-examples/consul/linkerd.yml

test_config ~/linkerd-examples/dcos/ingress/linkerd-config.yml
test_config ~/linkerd-examples/dcos/linker-to-linker/linkerd-config.yml
test_config ~/linkerd-examples/dcos/linker-to-linker-with-namerd/linkerd-config.yml
test_config ~/linkerd-examples/dcos/linkerd-with-namerd/linkerd-config.yml
test_config ~/linkerd-examples/dcos/simple-proxy/linkerd-config.yml

cp -a ~/linkerd-examples/failure-accrual/disco /
test_config ~/linkerd-examples/failure-accrual/linkerd.yml

cp -a ~/linkerd-examples/getting-started/docker/disco /io.buoyant/
test_config ~/linkerd-examples/getting-started/docker/linkerd.yaml

cp -a ~/linkerd-examples/getting-started/docker/disco .
test_config ~/linkerd-examples/getting-started/local/linkerd.yaml

test_config ~/linkerd-examples/gob/config/linkerd.yml
test_config ~/linkerd-examples/gob/dcos/linkerd-dcos-gob.yaml

test_config ~/linkerd-examples/http-proxy/linkerd.yaml

cp -a ~/linkerd-examples/influxdb/disco /
test_config ~/linkerd-examples/influxdb/linkerd1.yml
test_config ~/linkerd-examples/influxdb/linkerd2.yml

test_config ~/linkerd-examples/linkerd-tcp/linkerd.yml

test_config ~/linkerd-examples/mesos-marathon/linkerd-config.yml

echo "All test completed successfully"
