#!/bin/sh

set -e

sleep 10
curl -s "${K8S_API:-localhost:8001}/api/v1/namespaces/$NS/pods/$POD_NAME" | jq '.status.hostIP' | sed 's/"//g'
