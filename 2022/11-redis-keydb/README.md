Here are examples for testing the Redis response.

Application sets key in Redis-cluster every 1 second and export latency metrics in Prometheus format: max, min and average time of operation.

Docker:

```
docker pull trublast/redis-latency
```

Build:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -o redis-latency *.go
```