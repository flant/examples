#!/bin/bash

# VARIABLES
. variables

OS=$(uname)

if [ $(type kubectl > /dev/null; echo $?) != '0' ]
 then {
   echo "Kubectl not found, here's how to install:"
   ./info/install-kubectl.sh
   exit 1
   }
fi

if [ $(type minikube > /dev/null; echo $?) != '0' ]
 then {
   echo "Minikube not found, here's how to install:"
   ./info/install-minikube.sh
   exit 1
   }
fi

if [ $(type werf > /dev/null; echo $?) != '0' ]
 then {
   echo "werf not found, here's how to install:"
   ./info/install-werf.sh
   exit 1
   }
fi

if [ ! -f "/etc/docker/daemon.json" ]
 then {
   echo "daemon.json missing, configure docker:"
   ./info/configure-docker.sh
   exit 1
   }
fi

if [ -z $WERF_INSECURE_REGISTRY ]
 then {
   echo "WERF_INSECURE_REGISTRY not set. :"
   ./info/registry.sh
   exit 1
   }
fi

echo \
'{
    "insecure-registries": ["'$REPO'"]
}' | sudo tee /etc/docker/daemon.json
sudo systemctl restart docker

minikube start --driver=docker --insecure-registry="$REPO"
minikube addons enable ingress
minikube addons enable registry

minikube ssh -- "echo $(minikube ip) $REPO | sudo tee -a /etc/hosts"

echo "$(minikube ip)  $REPO" | sudo tee -a /etc/hosts
sed -e "s|hostname|$REPO|g" yaml/registry-ingress.yaml | kubectl create -f -

echo '
# Environment prepared. Werf docs: https://werf.io/documentation/v1.2/using_with_ci_cd_systems.html
'
