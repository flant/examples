FROM golang:buster@sha256:eb3f9ac805435c1b2c965d63ce460988e1000058e1f67881324746362baf9572
COPY . /app
RUN cd /app && \
    go build main.go && \
    chmod +x /app/main
CMD /app/main
