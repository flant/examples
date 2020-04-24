# habr-configmap

## Install

```shell
git clone https://github.com/drdeimos/habr-configmap
helm install \
  ./habr-configmap/charts/habr-configmap/ \
  --name habr-configmap \
  --namespace habr-configmap \
  --set 'name.production=Tod' \
  --set 'global.env=production'
```

## Upgrade

```shell
git clone https://github.com/drdeimos/habr-configmap
helm upgrade \
  habr-configmap \
  ./habr-configmap/charts/habr-configmap/ \
  --set 'name.production=Mary' \
  --set 'global.env=production'
```
