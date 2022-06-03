A simple application to demonstrate [local development in Kubernetes with werf 1.2 and minikube](https://blog.flant.com/local-development-in-kubernetes-with-werf/).

> NB. Tested with Ubuntu 20.04 only.

To start, follow these steps:

```
cd local
./install.sh
./setup-infra.sh
./deploy-app.sh
kubectl -n local get po
```

Then visit http://test.application.local/.
