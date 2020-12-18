# Enhanced version of etcdhelper

This is a modified version of [etcdhelper](https://github.com/openshift/origin/tree/master/tools/etcdhelper) from OpenShift (modifications are made by [Flant](https://flant.com/)).
Two new functions are introduced: `changeServiceCIDR` and `changePodCIDR`.

This repo is intended to be used as a playground for a corresponding article:

* «**[How to modify etcd data of your Kubernetes directly (without K8s API)](https://medium.com/flant-com/modifying-kubernetes-etcd-data-ed3d4bb42379)**».
* Russian version: «[Наш опыт работы с данными в etcd Kubernetes-кластера напрямую (без K8s API)](https://habr.com/ru/company/flant/blog/501956/)».

# Using etcdhelper

## Building

### Install golang

```shell
GOPATH=/root/golang
mkdir -p $GOPATH/local
curl -sSL https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz | tar -xzvC $GOPATH/local
echo "export GOPATH=\"$GOPATH\"" >> ~/.bashrc
echo 'export GOROOT="$GOPATH/local/go"' >> ~/.bashrc
echo 'export PATH="$PATH:$GOPATH/local/go/bin"' >> ~/.bashrc
```

### Install dependencies

```shell
go get go.etcd.io/etcd/clientv3 k8s.io/kubectl/pkg/scheme k8s.io/apimachinery/pkg/runtime
```

### Build

```shell
go build -o etcdhelper etcdhelper.go
```

## Using

### Change service CIDR

```shell
./etcdhelper -cacert /etc/kubernetes/pki/etcd/ca.crt -cert /etc/kubernetes/pki/etcd/server.crt -key /etc/kubernetes/pki/etcd/server.key -endpoint https://127.0.0.1:2379 change-service-cidr 172.30.0.0/16
```

### Change pod CIDR

```shell
./etcdhelper -cacert /etc/kubernetes/pki/etcd/ca.crt -cert /etc/kubernetes/pki/etcd/server.crt -key /etc/kubernetes/pki/etcd/server.key -endpoint https://127.0.0.1:2379 change-pod-cidr 10.55.0.0/16
```

### Change Ceph PV monitors list

This feature was added later (#15), thus it was not mentioned in the original article.

```shell
./etcdhelper -cacert /etc/kubernetes/pki/etcd/ca.crt -cert /etc/kubernetes/pki/etcd/server.crt -key /etc/kubernetes/pki/etcd/server.key -endpoint https://127.0.0.1:2379 change-monitors-list pvc-d748b019-52be-4c0c-a928-44503ccd94ac 10.0.1.1:6789,10.0.1.2:6789,10.0.1.3:6789
```

# Status

This enhanced version of etcdhelper is **PoC (proof of concept)**. Use it on your own risk.
