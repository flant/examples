#!/bin/bash

# VARIABLES
. variables

RELEASE=infra

export WERF_INSECURE_REGISTRY=1

sudo echo
./prepare.sh
while ! [[ "$(curl -w %{http_code} -sI http://$REPO/v2/ -o /dev/null)" -eq 200 ]] ; do echo "Waiting for registry..."; sleep 2; done;
cd ..
werf converge --repo $APP_REPO --release $RELEASE --env $ENV --namespace $NAMESPACE --dev --set mysql.enabled=true --ignore-secret-key=true
cd ./local

echo "$(minikube ip)  $URL" | sudo tee -a /etc/hosts