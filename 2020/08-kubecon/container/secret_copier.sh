#!/bin/bash

source /shell_lib.sh

function __config__() {
  cat << EOF
    configVersion: v1
    kubernetes:
    - name: src_secret
      apiVersion: v1
      kind: Secret
      nameSelector:
        matchNames:
        - mysecret
      namespace:
        nameSelector:
          matchNames: ["default"]
      group: main
    - name: dst_secrets
      apiVersion: v1
      kind: Secret
      labelSelector:
        matchLabels:
          managed-secret: "yes"
      jqFilter: |
        {
          "namespace": .metadata.namespace,
          "resourceVersion": .metadata.annotations.resourceVersion
        }
      group: main
      keepFullObjectsInMemory: false
    - name: namespaces
      group: main
      apiVersion: v1
      kind: Namespace
      jqFilter: |
        {
          name: .metadata.name,
          hasLabel: (.metadata.labels // {} | contains({"secret": "yes"}))
        }
      group: main
      keepFullObjectsInMemory: false
EOF
}

function sync_secret() {
  src_resource_version="$(context::jq -r '.snapshots.src_secret[0].object.metadata.resourceVersion')"
  dst_resource_version="$(context::jq -r '.snapshots.dst_secrets[].filterResult | select(.namespace == "'"$1"'") | .resourceVersion')"

  if [ "$src_resource_version" != "$dst_resource_version" ] ; then
    new_secret="$(context::jq -r '.snapshots.src_secret[0].object | . |= (.metadata =
                                                                            {"name": .metadata.name, "namespace": "'"$1"'",
                                                                            "labels": {"managed-secret": "yes"},
                                                                            "annotations": {"resourceVersion": .metadata.resourceVersion}})')"

    kubectl -n "$1" replace -f <(echo "$new_secret") || kubectl -n "$1" create -f  <(echo "$new_secret")
  fi
}

function delete_secret() {
  if context::jq -e --arg ns "$1" '.snapshots.dst_secrets[].filterResult | select(.namespace == $ns)' ; then
    if [[ "$1" == "default" ]]; then
      return 0
    fi

    kubectl -n "$1" delete secret "mysecret"
  fi
}

function __main__() {
  for i in $(seq 0 "$(context::jq -r '(.snapshots.namespaces | length) - 1')"); do
    ns_name="$(context::jq -r '.snapshots.namespaces['"$i"'].filterResult.name')"
    if context::jq -e '.snapshots.namespaces['"$i"'].filterResult.hasLabel' ; then
      sync_secret "$ns_name"
    else
      delete_secret "$ns_name"
    fi
  done
}

hook::run "$@"
