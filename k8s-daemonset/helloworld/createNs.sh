#!/bin/bash

set -e

for dtab in `ls dtab`; do
  namerctl --base-url http://namerd.${NS:-default}.svc.cluster.local:4180 dtab \
    create `echo $dtab | cut -d '.' -f1` dtab/$dtab
done
