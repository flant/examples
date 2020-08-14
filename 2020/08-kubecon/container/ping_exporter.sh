#!/bin/bash

source /shell_lib.sh

function __config__() {
  cat << EOF
    configVersion: v1
    kubernetes:
    - name: nodes
      apiVersion: v1
      kind: Node
      jqFilter: |
        {
          name: .metadata.name,
          ip: (
          .status.addresses[] |
            select(.type == "InternalIP") |
            .address
          )
        }
      group: main
      keepFullObjectsInMemory: false
      executeHookOnEvent: []
    schedule:
    - name: every_minute
      group: main
      crontab: "* * * * *"
EOF
}

function __main__() {
  for i in $(seq 0 "$(context::jq -r '(.snapshots.nodes | length) - 1')"); do
    node_name="$(context::jq -r '.snapshots.nodes['"$i"'].filterResult.name')"
    node_ip="$(context::jq -r '.snapshots.nodes['"$i"'].filterResult.ip')"

    packets_lost=0
    if ! ping -c 1 "$node_ip" -t 1 ; then
      packets_lost=1
    fi

    cat > "$METRICS_PATH" <<END
      {
        "name": "node_packets_lost",
        "add": $packets_lost,
        "labels": {
          "node": $node_name
        }
      }
END

  done
}

hook::run "$@"
