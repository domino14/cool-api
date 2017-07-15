FROM golang:alpine
ENV GOPATH=/go
ADD . /go/src/github.com/domino14/cool-api
WORKDIR /go/src/github.com/domino14/cool-api

EXPOSE 8086