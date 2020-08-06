# Running Kubernetes on Raspberry Pi

Here are example YAML manifests to simplify deploying basic components for Kubernetes cluster installed on Raspberry Pi.

This repo is intended to be used as a playground for a corresponding article:

* «**[TBA](TBA)**».
* Russian version: «[TBA](TBA)».

## Contents

* `prometheus-pvc.yaml` — PV (persistent volume) and PVC (persistent volume claim) for a simple local storage
  (based on hostpath) for Prometheus (including AlertManager data);
* `grafana-cluster-issuer.yaml` & `grafana-certificate.yaml` — cluster issuer & SSL certificate itself
  for Grafana (using cert-manager).
