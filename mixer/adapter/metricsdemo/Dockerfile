# #hadolint ignore=DL3026
# FROM golang:1.11 as BUILD

# RUN apt-get update && apt -y install unzip

# RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip
# RUN unzip -o protoc-3.7.0-linux-x86_64.zip -d /usr/local bin/protoc
# RUN go get -u github.com/gogo/protobuf/types github.com/golang/dep/cmd/dep

# WORKDIR /go/src/github.com/mulesoft/service-mesh-grpc-adapter/grpcmule
# COPY grpcmule/ .
# COPY Gopkg.* ../
# RUN mkdir ../bin
# COPY bin/codegen.sh ../bin

# RUN dep ensure
# RUN go generate ./... && go build ./... && go test ./...

# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/grpcmule ./cmd/

###############
# PRODUCTION
###############

# prerequisites
#   go generate ./... && go build ./...
#   CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/grpcmule ./cmd/

# hadolint ignore=DL3026
FROM alpine:3.8 as PRODUCTION

# hadolint ignore=DL3018
RUN apk --no-cache add ca-certificates
WORKDIR /bin/
#COPY --from=BUILD /go/src/github.com/mulesoft/service-mesh-grpc-adapter/grpcmule/bin/grpcmule .
COPY bin/metricsd .
ENTRYPOINT [ "/bin/metricsd" ]
CMD [ "8005" ]
EXPOSE 8005
