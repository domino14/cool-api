FROM alpine:latest
RUN apk add --no-cache go libc-dev git
EXPOSE 8086
ENV GOPATH /go

CMD ['go', 'run', 'main.go']

