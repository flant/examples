# Доклад DUMP SPB 2025

Оператор для Kubernetes: backend для людей и автоматики


## Инструменты для разработки

- Библиотека client-go

  - go get k8s.io/client-go

  - https://github.com/kubernetes/sample-controller 

- Библиотека controller-runtime 

  - https://github.com/kubernetes-sigs/controller-runtime/

- kube-builder генератор

  - https://book.kubebuilder.io/quick-start.html

- operator framework от RedHat

  - https://sdk.operatorframework.io/

- Kopf для Python

  - https://github.com/nolar/kopf

- shell-operator для скриптовых языков и не только

  - https://github.com/flant/shell-operator

  - Примеры операторов

    - https://github.com/flant/examples/tree/master/2020/08-kubecon

    - https://habr.com/ru/companies/flant/articles/447442/


## Репозитории операторов

- https://operatorhub.io

- https://artifacthub.io/packages/search


## Дополнительные ресурсы для дальнейшего погружения

- KubernetesAPI conventions от участников SIG Architecture
  - https://github.com/kubernetes/community/blob/8a99192b3780b656f9dd53c0c37d9372a1c975f9/contributors/devel/sig-architecture/api-conventions.md

- Размышления о том, что учесть в разработке операторов
  - https://ahmet.im/blog/controller-pitfalls/

- Немного примеров от CNCF
  - https://www.cncf.io/blog/2022/06/15/kubernetes-operators-what-are-they-some-examples/

- Концепции оператора и контроллера из документации
  - https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
  - https://kubernetes.io/docs/concepts/architecture/controller/

- Почему не стоит писать оператор?
  - https://rm-rf.ca/posts/2020/when-not-to-write-kubernetes-operator/

