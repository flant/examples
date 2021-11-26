# Deploying Keycloak with Infinispan in Kubernetes

Here are sample configurations that can be used to deploy Keycloak in Kubernetes clusters.

This repo is intended to be used as a playground for a corresponding article:

* «[Running fault-tolerant Keycloak with Infinispan in Kubernetes](https://blog.flant.com/ha-keycloak-infinispan-kubernetes/)».
* Russian version: «[Настраиваем отказоустойчивый Keycloak с Infinispan в Kubernetes](https://habr.com/ru/company/flant/blog/567626/)».

## Contents

* `.helm` — a Helm chart with Kubernetes manifests for Keycloak (`templates/keycloak-sts.yaml`
  is for its StatefulSet, `templates/keycloak-cm.yaml` is for its ConfigMap) and Infinispan
  (`templates/infinispan-sts.yaml`, `templates/inifinispan-cm.yaml`);
* `werf.yaml` — a simple [werf](https://werf.io/) configuration to build and deploy Keycloak
  with Infinispan that will have a newer version of the PostgreSQL driver
  (required for CockroachDB);
* `jar` — JAR files used in `werf.yaml`.
