FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/ravjotsingh9/nats-streaming-functional-service

COPY Gopkg.lock Gopkg.toml ./
COPY vendor vendor
COPY functional-service functional-service 
COPY publisher publisher
COPY nats nats 
COPY util util

RUN go install ./...

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .
