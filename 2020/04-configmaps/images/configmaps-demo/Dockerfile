FROM golang:1.13-alpine@sha256:e74b83b94d499cb34c7edf55fbdff9d3cfabd238a35f9cb1b59d0e74f50b1ddd

RUN adduser -u 40004 -h /app -D app

RUN mkdir /app/configfiles /app/secretfiles && \
    chown -R app:app /app/configfiles /app/secretfiles

RUN apk update && apk add git

RUN go get github.com/fsnotify/fsnotify

RUN cd /go/src/github.com && \
    git clone https://github.com/flant/examples.git flant-examples

RUN cd /go/src/github.com/flant-examples/2020/04-configmaps/src/ && \
    go build -o /bin/configmaps-demo main.go
