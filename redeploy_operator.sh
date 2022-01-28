#!/bin/bash
set -x #echo on


rm -rf bundle
kubectl delete -f config/rbac
operator-sdk cleanup egbtmagent --delete-all=true
export GO111MODULE=on
go mod tidy
make generate
make manifests
docker logout
echo "$DOCKER_PASSWORD" | docker login --username egapm --password-stdin
make docker-build docker-push
make bundle
make bundle-build bundle-push
kubectl create -f config/rbac
operator-sdk run bundle 172.16.8.78:5000/egbtmagent-bundle:v0.0.1 --skip-tls --timeout 15m0s
