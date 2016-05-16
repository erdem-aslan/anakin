#!/bin/bash
export GOPATH=/Users/erdemaslan/Projects/go
export GOSRC=$GOPATH/src
export GOOS=linux
export GOARCH=amd64

echo "Building 64bit linux..."
/usr/local/go/bin/go build -o dist/linux/64/anakin

export GOARCH=386

echo "Building 32bit linux..."
/usr/local/go/bin/go build -o dist/linux/32/anakin

export GOOS=darwin
export GOARCH=amd64

echo "Building 64bit macos..."
/usr/local/go/bin/go build -o dist/macos/anakin

GOOS=windows

echo "Building 64bit windows..."
/usr/local/go/bin/go build -o dist/windows/64/anakin.exe

export GOARCH=386

echo "Building 32bit windows..."
/usr/local/go/bin/go build -o dist/windows/32/anakin.exe

echo "Packing web ui..."

tar cfz anakin-web-ui.tgz web
mv anakin-web-ui.tgz dist/

echo "Finished building anakin!"
