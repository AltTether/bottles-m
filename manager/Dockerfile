FROM golang:1.14.4-alpine3.11

RUN apk --update add \
    git \
    && rm -rf /var/cache/apk/*

#RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/gin-gonic/gin
#RUN go get -u gopkg.in/olahol/melody.v1
RUN go get -u github.com/gin-contrib/static

WORKDIR /go/src/github.com/manager
