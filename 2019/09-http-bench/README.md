# MINI MANUAL

1) First of all(of course after clone this repo) build your image
```shell
docker build -t benchmark:v1 . 
```

2) Then start container
```shell
docker run --rm -d --add-host=benchmark.test:127.0.0.1 --name go_benchmark benchmark:v1
```

3) Test it
```shell
docker exec -it go_benchmark ./benchmark -u https://benchmark.test/bench -c 100 -n 100000
```
4) Ah or even easier, only in one step
 ```shell
make bench 
```
5) That`s all. Bye :-)
    
    5.1) P.S. Have a nice day!