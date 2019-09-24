# hadolint ignore=DL3026
FROM golang:1.11 as BUILD

RUN apt-get update && apt -y install unzip

RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip
RUN unzip -o protoc-3.7.0-linux-x86_64.zip -d /usr/local bin/protoc
RUN go get github.com/gogo/protobuf/types

WORKDIR /go/src/istio.io/
RUN git clone --branch 1.0.6 --depth 1 https://github.com/istio/istio 

WORKDIR /go/src/istio.io/istio/mixer/adapter/metricsdemo
COPY cmd ./cmd
COPY config ./config
COPY metricsdemo.go ./

RUN go generate ./... && go build ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/metricsd ./cmd/

###############
# PRODUCTION
###############

# hadolint ignore=DL3026
FROM alpine:3.8 as PRODUCTION

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates
WORKDIR /bin/
COPY --from=BUILD /go/src/istio.io/istio/mixer/adapter/metricsdemo/bin/metricsd .
#COPY bin/metricsd .
ENTRYPOINT [ "/bin/metricsd" ]
CMD [ "8005" ]
EXPOSE 8005
