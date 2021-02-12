# Using Prometheus exporters with AWS CloudWatch

Here are Kubernetes manifests to deploy cloudwatch_exporter and prometheus_aws_cost_exporter in Kubernetes for getting AWS CloudWatch metrics in Prometheus.

This repo is intended to be used as a playground for a corresponding article:

* Russian version: «[Мониторим основные сервисы в AWS с Prometheus и exporter’ами для CloudWatch](https://habr.com/ru/company/flant/blog/542082/)».

## Contents

* `.helm` — Helm chart with Kubernetes manifests for cloudwatch_exporter & prometheus_aws_cost_exporter; 
  it includes sample Prometheus rules for alerts (`.helm/templates/60-rules.yaml`);
* `role.json` — AWS IAM permissions for prometheus_aws_cost_exporter;
* `policy.json` & `terraform_user_policy.tf` — AWS IAM policy for API (in JSON) and Terraform.
