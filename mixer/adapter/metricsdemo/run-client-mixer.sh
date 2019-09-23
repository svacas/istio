#!/bin/bash

export ISTIO=$GOPATH/src/istio.io

echo build mixer client...
pushd $ISTIO/istio && make mixc

echo
echo run mixer client...
$GOPATH/out/darwin_amd64/release/mixc report -t request.time="2006-01-02T15:04:05Z",response.time="2006-01-02T15:04:05Z" -s destination.service.namespace="default",destination.service="svc.cluster.local",request.path="svc-path",request.method="get" -i request.size=1235

