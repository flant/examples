k8s-python-operator-example
---------------------------
Kubernetes operator written in Python.

* «[Writing a Kubernetes Operator in Python without frameworks and SDK](https://blog.flant.com/writing-a-kubernetes-operator-in-python-without-frameworks-and-sdk/)».
* Russian version: «[Kubernetes Operator на Python без фреймворков и SDK](https://habr.com/ru/company/flant/blog/459320/)».


#### Launching the operator
```bash
usage: copyrator [-h] [--namespace NAMESPACE] [--rule-name RULE_NAME]

Copyrator - copy operator.

optional arguments:
  -h, --help            show this help message and exit
  --namespace NAMESPACE
                        Operator Namespace (or ${NAMESPACE}), default: default
  --rule-name RULE_NAME
                        CRD Name (or ${RULE_NAME}), default: main-rule
``` 
