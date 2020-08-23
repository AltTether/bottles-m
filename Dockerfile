FROM golang:1.15.0-alpine3.12

RUN apk --update add alpine-sdk

ENV GO111MODULE=on

WORKDIR /go/src/github.com/bottles
COPY ./go.mod /go/src/github.com/bottles
COPY ./go.sum /go/src/github.com/bottles
RUN go mod download
