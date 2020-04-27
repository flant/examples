# ConfigMap demo

Helm chart with a simplistic app to demo how ConfigMaps are used/updated
in Kubernetes. It reads and prints data from config (`config.json`) as
well as watches the modifications of this config via fsnotify.

Intended to be used as a playground with corresponding article:
* Russian version: «[ConfigMaps в Kubernetes: нюансы, о которых стоит знать](https://habr.com/ru/company/flant/blog/498970/)».

## Deploying this chart

```shell
git clone https://github.com/drdeimos/habr-configmap
helm install \
  ./habr-configmap/charts/habr-configmap/ \
  --name habr-configmap \
  --namespace habr-configmap \
  --set 'name.production=Tod' \
  --set 'global.env=production'
```

## Upgrading the chart

```shell
git clone https://github.com/drdeimos/habr-configmap
helm upgrade \
  habr-configmap \
  ./habr-configmap/charts/habr-configmap/ \
  --set 'name.production=Mary' \
  --set 'global.env=production'
```
