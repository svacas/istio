#!/usr/bin/env bash

set -e

echo "go generate..."
go generate ./...

echo "go build..."
go build ./...

echo "go build linux binary..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/metricsd ./cmd/

echo done

## post

# docker build -t 302195489440.dkr.ecr.us-east-1.amazonaws.com/mulesoft/metricsd-adapter:0.0.1_3 .
# docker push 302195489440.dkr.ecr.us-east-1.amazonaws.com/mulesoft/metricsd-adapter:0.0.1_3
