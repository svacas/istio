#!/usr/bin/env bash

set -e

echo "go generate..."
go generate ./...

echo "go build..."
go build ./...

echo "go build linux binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/metricsd ./cmd/

echo done
