FROM golang:buster
COPY . /app
RUN cd /app && \
    go build main.go && \
    chmod +x /app/main
CMD /app/main
