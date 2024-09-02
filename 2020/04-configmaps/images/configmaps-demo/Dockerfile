FROM golang:1.23-alpine@sha256:d0b31558e6b3e4cc59f6011d79905835108c919143ebecc58f35965bf79948f4

RUN adduser -u 40004 -h /app -D app

RUN mkdir /app/configfiles /app/secretfiles && \
    chown -R app:app /app/configfiles /app/secretfiles

RUN apk update && apk add git

RUN go get github.com/fsnotify/fsnotify

RUN cd /go/src/github.com && \
    git clone https://github.com/flant/examples.git flant-examples

RUN cd /go/src/github.com/flant-examples/2020/04-configmaps/src/ && \
    go build -o /bin/configmaps-demo main.go
