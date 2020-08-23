FROM golang:1.15.0-alpine3.12

ENV GO111MODULE=on

WORKDIR /go/src/github.com/bottles
COPY ./go.mod /go/src/github.com/bottles
RUN go mod download
