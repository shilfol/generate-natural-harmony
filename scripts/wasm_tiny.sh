#!/bin/bash

docker run -v $(pwd):/go/src/github.com/shilfol/generate-natural-harmony -v $GOPATH:/go \
-w /go/src/github.com/shilfol/generate-natural-harmony \
-e "GOPATH=/go" tinygo/tinygo:0.16.0 tinygo \
build -o /go/src/github.com/shilfol/generate-natural-harmony/web/wasm/nh_tiny.wasm -target wasm \
--no-debug /go/src/github.com/shilfol/generate-natural-harmony/cmd/wasm/main.go