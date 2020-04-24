FROM golang:1.13-alpine

RUN adduser -u 40004 -h /app -D app

RUN mkdir /app/configfiles /app/secretfiles && \
    chown -R app:app /app/configfiles /app/secretfiles

RUN apk update && apk add git

RUN go get github.com/fsnotify/fsnotify

RUN cd /go/src/github.com && \
    git clone https://github.com/drdeimos/habr-configmap

RUN cd /go/src/github.com/habr-configmap/src && \
    go build -o /bin/habr-configmap main.go
