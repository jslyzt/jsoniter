#!/bin/bash
set -e
set -x

if [ ! -d /tmp/build-golang/src/github.com/json-iterator ]; then
    mkdir -p /tmp/build-golang/src/github.com/json-iterator
    ln -s $PWD /tmp/build-golang/src/github.com/jslyzt/jsoniter
fi
export GOPATH=/tmp/build-golang
go get -u github.com/golang/dep/cmd/dep
cd /tmp/build-golang/src/github.com/jslyzt/jsoniter
exec $GOPATH/bin/dep ensure -update
