#!/bin/bash
export GOPATH=/Users/erdemaslan/Projects/go
export GOSRC=$GOPATH/src
export GOOS=linux
export GOARCH=amd64

echo "Building 64bit linux..."
/usr/local/go/bin/go fmt
/usr/local/go/bin/go build -o dist/linux/64/anakin
echo "Finished building anakin!"
