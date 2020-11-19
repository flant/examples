# Running Kubernetes on Raspberry Pi

Here are example YAML manifests to simplify deploying basic components for Kubernetes cluster installed on Raspberry Pi.

This repo is intended to be used as a playground for a corresponding article:

* «**[Installing fully-fledged vanilla Kubernetes on Raspberry Pi](https://medium.com/flant-com/install-vanilla-kubernetes-raspberry-pi-b21895a51550)**».
* Russian version: «[Полноценный Kubernetes с нуля на Raspberry Pi](https://habr.com/ru/company/flant/blog/513908/)».

## Contents

* `prometheus-pv.yaml` — PVs (persistent volumes) for a simple local storage (based on hostpath) for Prometheus
  (including AlertManager data);
* `cert-manager-cluster-issuer.yaml` & `cert-manager-grafana-certificate.yaml` — cert-manager's cluster issuer & SSL
  certificate for Grafana.
