FROM golang:stretch@sha256:a9718cb62e77cc03411ea11fd71f3541e508f211169dd8f8d69ba2036545cd80 as build
WORKDIR /app
COPY app/ ./
RUN export GOBIN=$PWD/bin && \
    export PATH=$GOBIN:$PATH && \
    export GO111MODULE=on && \
    go get "github.com/google/uuid"
RUN export GOBIN=$PWD/bin && \
    export PATH=$GOBIN:$PATH && \
    export GO111MODULE=on && \
    go build ./benchmark.go

FROM nginx:1.16@sha256:d20aa6d1cae56fd17cd458f4807e0de462caf2336f0b70b5eeb69fcaaf30dd9c as final
WORKDIR /root
RUN apt update && \
    apt -y install vim unzip git curl && \
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -keyout /etc/ssl/private/nginx-selfsigned.key \
    -out /etc/ssl/certs/nginx-selfsigned.crt \
    -subj "/C=RU/ST=Moscow/L=Moscow/O=Flant/CN=benchmark.test/" \
    -addext "subjectAltName = DNS:benchmark.test,DNS:localhost" && \
    mkdir -p  /usr/local/share/ca-certificates/bench && \
    cp /etc/ssl/certs/nginx-selfsigned.crt /usr/local/share/ca-certificates/bench && \
    chmod 644 /usr/local/share/ca-certificates/bench && \
    chmod 755 /usr/local/share/ca-certificates/bench/nginx-selfsigned.crt && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /app/benchmark  ./
COPY nginx.conf /etc/nginx/nginx.conf

CMD ["nginx", "-g", "daemon off;"]