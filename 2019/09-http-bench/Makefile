build:
	docker build -t benchmark:v1 .
build-apache:
	docker build -f Dockerfile.apache -t benchmark-apache:v1 .

run: build
	docker run --rm -d --add-host=benchmark.test:127.0.0.1 --name go_benchmark benchmark:v1

run-apache: build-apache
	docker run --rm -d --add-host=benchmark.test:127.0.0.1 --name go_benchmark benchmark-apache:v1

bench: run
	docker exec -it go_benchmark ./benchmark -u https://benchmark.test/bench -c 1000 -n 10000
	docker rm -f go_benchmark 2>&1 1>/dev/null

bench-apache: run-apache
	docker exec -it go_benchmark ./benchmark -u https://benchmark.test/bench -c 1000 -n 10000
	docker rm -f go_benchmark 2>&1 1>/dev/null
