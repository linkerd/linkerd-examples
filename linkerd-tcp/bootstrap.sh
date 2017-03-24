#!/bin/bash

set -e

/run.sh &

prometheus_data_source=$(cat <<EOF
{
  "access": "proxy",
  "isDefault": true,
  "jsonData": {},
  "name": "prometheus",
  "type": "prometheus",
  "url": "http://127.0.0.1:9090"
}
EOF
)

until $(curl -sfo /dev/null http://127.0.0.1:3000/api/datasources); do
  echo "waiting for grafana to start"
  sleep 1
done

curl -sX POST -d "${prometheus_data_source}" -H "Content-Type: application/json" \
  http://127.0.0.1:3000/api/datasources

wait
