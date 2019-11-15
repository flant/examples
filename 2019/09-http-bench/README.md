This repo demonstrates how a simple benchmark (written in Golang and built into Docker container) works against your HTTP server. More details are described in our [6 recent cases from our SRE workaday routine](https://medium.com/flant-com/6-sre-troubleshooting-cases-faf72ed36d6b?source=friends_link&sk=a72ed571f28327138c76edbbf5aab168) article (case #1).

# How to use it

Clone this repo and execute:

```shell
make bench 
```
… or do everything step by step:

1) Build your image:
```shell
docker build -t benchmark:v1 . 
```

2) Start a container:
```shell
docker run --rm -d --add-host=benchmark.test:127.0.0.1 --name go_benchmark benchmark:v1
```

3) Run the benchmark for it:
```shell
docker exec -it go_benchmark ./benchmark -u https://benchmark.test/bench -c 100 -n 100000
```
