FROM golang:1.14.4-alpine3.11

RUN apk --update add \
    git \
    mysql-client \
    && rm -rf /var/cache/apk/*

#RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/gin-gonic/gin \
    && go get github.com/go-sql-driver/mysql

RUN go get -u github.com/jinzhu/gorm

RUN go get -u github.com/stretchr/testify

RUN apk add alpine-sdk

WORKDIR /go/src/github.com/message
