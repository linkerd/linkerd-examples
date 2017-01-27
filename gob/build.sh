#!/bin/sh

set -ex

for svc in gen web word ; do
  GOOS=linux GOARCH=amd64  go build -o $svc.linux_amd64 $svc/main.go
  GOOS=darwin GOARCH=amd64 go build -o $svc.darwin      $svc/main.go
done
