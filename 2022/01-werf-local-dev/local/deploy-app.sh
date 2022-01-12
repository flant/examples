#!/bin/bash

# VARIABLES
. variables

RELEASE=application

export WERF_INSECURE_REGISTRY=1

cd ../
if [[ $# -gt 0 ]]; then
    for argument in "$@"; do
        case $argument in
            -f|--follow)
                DEFAULT="--dev --follow"
                werf converge --repo $APP_REPO --release $RELEASE --env $ENV --namespace $NAMESPACE --set app.ci_url=$URL --set app.enabled=true $MODE --ignore-secret-key=true
                cd ./local
                ;;
            *)
                DEFAULT="--dev"
                werf converge --repo $APP_REPO --release $RELEASE --env $ENV --namespace $NAMESPACE --set app.ci_url=$URL --set app.enabled=true $MODE --ignore-secret-key=true
                cd ./local
                ;;
        esac
    done
else
    werf converge --repo $APP_REPO --release $RELEASE --env $ENV --namespace $NAMESPACE --set app.ci_url=$URL --set app.enabled=true $MODE --ignore-secret-key=true
fi