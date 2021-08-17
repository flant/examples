# High available Memcached variant in Kubernetes (using mcrouter)

Here are sample configurations that can be used to deploy Memcached+mcrouter DaemonSet in Kubernetes clusters.

## Contents

* `.helm` — a Helm chart with Kubernetes manifests for Memcached+mcrouter (`templates/memcached-ds.yaml`
  is for Memcached DaemonSet, `templates/mcrouter-ds.yaml` is for mcrouter DaemonSet, `templates/mcrouter-cm.yaml` is for mcrouter ConfigMap);
* `werf.yaml` — a simple [werf](https://werf.io/) configuration to build and deploy Memcached+mcrouter
