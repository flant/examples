Here are sample configurations that can be used to deploy a high-available memcached in Kubernetes using mcrouter.

This repo is intended to be used as a playground for a corresponding article:

* Russian version: «[Готовим высокодоступный memcached с mcrouter в Kubernetes](https://habr.com/ru/company/flant/blog/575656/)».

## Contents

* `.helm` — a Helm chart with Kubernetes manifests for memcached & mcrouter (`templates/memcached-ds.yaml`
  is for memcached DaemonSet, `templates/mcrouter-ds.yaml` is for mcrouter DaemonSet, `templates/mcrouter-cm.yaml` is for mcrouter ConfigMap);
* `werf.yaml` — a simple [werf](https://werf.io/) configuration to build and deploy memcached & mcrouter.
