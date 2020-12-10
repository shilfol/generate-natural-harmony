#!/bin/bash

GOOS=js GOARCH=wasm go build -o ./web/wasm/nh.wasm ./cmd/wasm/main.go