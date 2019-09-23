#!/bin/bash

export ISTIO=${GOPATH}/src/istio.io
export ADAPTER_NAME=metricsdemo

echo build mixer...
pushd $ISTIO/istio && make mixs

echo
echo run mixer...
${GOPATH}/out/darwin_amd64/release/mixs server --configStoreURL=fs://$(pwd)/mixer/adapter/${ADAPTER_NAME}/testdata-local
