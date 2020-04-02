## Install golang

```shell
GOPATH=/root/golang
mkdir -p $GOPATH/local
curl -sSL https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz | tar -xzvC $GOPATH/local
echo "export GOPATH=\"$GOPATH\"" >> ~/.bashrc
echo 'export GOROOT="$GOPATH/local/go"' >> ~/.bashrc
echo 'export PATH="$PATH:$GOPATH/local/go/bin"' >> ~/.bashrc
```

## Install dependencies

```shell
go get go.etcd.io/etcd/clientv3 k8s.io/kubectl/pkg/scheme k8s.io/apimachinery/pkg/runtime
```

## Build

```shell
go build -o etcdhelper etcdhelper.go
```

## Sample Usage

### Change service CIDR

```shell
./etcdhelper -cacert /etc/kubernetes/pki/etcd/ca.crt -cert /etc/kubernetes/pki/etcd/server.crt -key /etc/kubernetes/pki/etcd/server.key -endpoint https://127.0.0.1:2379 change-service-cidr 172.30.0.0/16
```

### Change pod CIDR

```shell
./etcdhelper -cacert /etc/kubernetes/pki/etcd/ca.crt -cert /etc/kubernetes/pki/etcd/server.crt -key /etc/kubernetes/pki/etcd/server.key -endpoint https://127.0.0.1:2379 change-pod-cidr 10.55.0.0/16
```
