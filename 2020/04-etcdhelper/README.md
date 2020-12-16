# Enhanced version of etcdhelper

This is a modified version of [etcdhelper](https://github.com/openshift/origin/tree/master/tools/etcdhelper) from OpenShift (modifications are made by [Flant](https://flant.com/)).
Two new functions are introduced: `changeServiceCIDR` and `changePodCIDR`.

This repo is intended to be used as a playground for a corresponding article:

* «**[How to modify etcd data of your Kubernetes directly (without K8s API)](https://medium.com/flant-com/modifying-kubernetes-etcd-data-ed3d4bb42379)**».
* Russian version: «[Наш опыт работы с данными в etcd Kubernetes-кластера напрямую (без K8s API)](https://habr.com/ru/company/flant/blog/501956/)».

# Using etcdhelper

## Build

The fastest way is to use official golang image:

```shell
docker run --rm -v $(pwd):/app -w /app -e CGO_ENABLED=0 -e GOOS=linux golang:1.15-alpine go build etcdhelper.go
```

etcdhelper binary will be created in the current directory. Use `GOOS=darwin` to build MacOS executable. 

## Using

### Change service CIDR

```shell
./etcdhelper -cacert /etc/kubernetes/pki/etcd/ca.crt -cert /etc/kubernetes/pki/etcd/server.crt -key /etc/kubernetes/pki/etcd/server.key -endpoint https://127.0.0.1:2379 change-service-cidr 172.30.0.0/16
```

### Change pod CIDR

```shell
./etcdhelper -cacert /etc/kubernetes/pki/etcd/ca.crt -cert /etc/kubernetes/pki/etcd/server.crt -key /etc/kubernetes/pki/etcd/server.key -endpoint https://127.0.0.1:2379 change-pod-cidr 10.55.0.0/16
```

# Status

This enhanced version of etcdhelper is **PoC (proof of concept)**. Use it on your own risk.
