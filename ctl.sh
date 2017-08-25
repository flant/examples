#!/usr/bin/env bash

set -e

! read -rd '' HELP_STRING <<"EOF"
Usage: ctl.sh [OPTION]... --gitlab-url GITLAB_URL --oauth2-id ID --oauth2-secret SECRET --dashboard-url DASHBOARD_URL
Install kubernetes-dashboard to Kubernetes cluster.
Mandatory arguments:
  -i, --install                install into 'kube-nginx-ingress' namespace
  -u, --upgrade                upgrade existing installation, will reuse password and host names
  -d, --delete                 remove everything, including the namespace
Optional arguments:
  -h, --help                   output this message
      --gitlab-url             set gitlab url with schema (https://gitlab.example.com)
      --oauth2-id              set OAUTH2_PROXY_CLIENT_ID from gitlab
      --oauth2-secret          set OAUTH2_PROXY_CLIENT_SECRET from gitlab
      --dashboard-url          set dashboard url without schema (dashboard.example.com)
EOF

RANDOM_NUMBER=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 4 | head -n 1)
TMP_DIR="/tmp/kubernetes-dashboard-$RANDOM_NUMBER"
WORKDIR="$TMP_DIR/kubernetes-dashboard"

if [ -z "${KUBECONFIG}" ]; then
    export KUBECONFIG=~/.kube/config
fi

FIRST_INSTALL="true"

TEMP=$(getopt -o i,u,d,h --long install,upgrade,delete,gitlab-url:,oauth2-id:,oauth2-secret:,dashboard-url:,help \
             -n 'ctl.sh' -- "$@")

eval set -- "$TEMP"

while true; do
  case "$1" in
    -i | --install )
      MODE=install; shift ;;
    -u | --upgrade )
      MODE=upgrade; shift ;;
    -d | --delete )
      MODE=delete; shift ;;
    --gitlab-url )
      GITLAB_URL="$2"; shift 2;;
    --oauth2-id )
      OAUTH2_PROXY_CLIENT_ID="$2"; shift 2;;
    --oauth2-secret )
      OAUTH2_PROXY_CLIENT_SECRET="$2"; shift 2;;
    --dashboard-url )
      DASHBORD_URL="$2"; shift 2;;
    -h | --help )
      echo "$HELP_STRING"; exit 0 ;;
    -- )
      shift; break ;;
    * )
      break ;;
  esac
done
if [ -z "${MODE}" ]; then
  echo "$HELP_STRING"; exit 1
fi
if [ "$MODE" == "install" ]; then
  if [ -z "${GITLAB_URL}" ]; then
    echo "$HELP_STRING"; exit 0
  fi
  if [ -z "${OAUTH2_PROXY_CLIENT_ID}" ]; then
    echo "$HELP_STRING"; exit 0
  fi
  if [ -z "${OAUTH2_PROXY_CLIENT_SECRET}" ]; then
    echo "$HELP_STRING"; exit 0
  fi
  if [ -z "${DASHBORD_URL}" ]; then
    echo "$HELP_STRING"; exit 0
  fi
fi
type git >/dev/null 2>&1 || { echo >&2 "I require git but it's not installed.  Aborting."; exit 1; }
type kubectl >/dev/null 2>&1 || { echo >&2 "I require kubectl but it's not installed.  Aborting."; exit 1; }
type jq >/dev/null 2>&1 || { echo >&2 "I require jq but it's not installed.  Aborting."; exit 1; }


mkdir -p "$TMP_DIR"
cd "$TMP_DIR"
git clone --depth 1 https://github.com/AndrewKoryakin/kubernetes-dashboard.git
cd "$WORKDIR"

LOGIN_URL="${GITLAB_URL}/oauth/authorize"
REDEEM_URL="${GITLAB_URL}/oauth/token"
VALIDATE_URL="${GITLAB_URL}/api/v3/user"
OAUTH2_PROXY_COOKIE_SECRET="$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 15 | head -n 1| base64)"

function install {
  sed -i -e "s%##LOGIN_URL##%$LOGIN_URL%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##REDEEM_URL##%$REDEEM_URL%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##VALIDATE_URL##%$VALIDATE_URL%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##OAUTH2_PROXY_CLIENT_ID##%$OAUTH2_PROXY_CLIENT_ID%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##OAUTH2_PROXY_CLIENT_SECRET##%$OAUTH2_PROXY_CLIENT_SECRET%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##OAUTH2_PROXY_COOKIE_SECRET##%$OAUTH2_PROXY_COOKIE_SECRET%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##DASHBORD_URL##%$DASHBORD_URL%g" manifests/kube-dashboard-ingress.yaml
  kubectl apply -Rf manifests/
}

function upgrade {
  if $(kubectl get deployment oauth2-proxy -n kube-system > /dev/null 2>/dev/null); then
    LOGIN_URL=$(kubectl get deployment oauth2-proxy -n kube-system -o json | jq -r '.spec.template.spec.containers[0]' |grep '\-login\-url'| sed -e 's/^[[:space:]]*//'| sed -e 's/,$//')
    REDEEM_URL=$(kubectl get deployment oauth2-proxy -n kube-system -o json | jq -r '.spec.template.spec.containers[0]' |grep '\-redeem\-url'| sed -e 's/^[[:space:]]*//'| sed -e 's/,$//')
    VALIDATE_URL=$(kubectl get deployment oauth2-proxy -n kube-system -o json | jq -r '.spec.template.spec.containers[0]' |grep '\-validate\-url'| sed -e 's/^[[:space:]]*//'| sed -e 's/,$//')
    OAUTH2_PROXY_COOKIE_SECRET=$(kubectl get deployment oauth2-proxy -n kube-system -o json | jq -r '.spec.template.spec.containers[0]' |grep 'OAUTH2_PROXY_COOKIE_SECRET' -A1  |grep value |awk -F ': ' '{print $2}')
    OAUTH2_PROXY_CLIENT_SECRET=$(kubectl get deployment oauth2-proxy -n kube-system -o json | jq -r '.spec.template.spec.containers[0]' |grep 'OAUTH2_PROXY_CLIENT_SECRET' -A1  |grep value |awk -F ': ' '{print $2}')
    OAUTH2_PROXY_CLIENT_ID=$(kubectl get deployment oauth2-proxy -n kube-system -o json | jq -r '.spec.template.spec.containers[0]' |grep 'OAUTH2_PROXY_CLIENT_ID' -A1  |grep value |awk -F ': ' '{print $2}')
  else
    echo "Can't upgrade. Deployment kubernetes-dashboard does not exists. " && exit 1
  fi
  if $(kubectl get deployment oauth2-proxy -n kube-system > /dev/null 2>/dev/null); then
    DASHBORD_URL=$(kubectl get ing oauth2-proxy -n kube-system -o json | jq -r '.spec.tls[0].hosts[0]')
  else
    echo "Can't upgrade. Ingress oauth2-proxy does not exists. " && exit 1
  fi
  sed -i -e "s%-login-url=##LOGIN_URL##%$LOGIN_URL%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%-redeem-url=##REDEEM_URL##%$REDEEM_URL%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%-validate-url=##VALIDATE_URL##%$VALIDATE_URL%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##OAUTH2_PROXY_CLIENT_ID##%$OAUTH2_PROXY_CLIENT_ID%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##OAUTH2_PROXY_CLIENT_SECRET##%$OAUTH2_PROXY_CLIENT_SECRET%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##OAUTH2_PROXY_COOKIE_SECRET##%$OAUTH2_PROXY_COOKIE_SECRET%g" manifests/kube-dashboard-oauth2-proxy.yaml
  sed -i -e "s%##DASHBORD_URL##%$DASHBORD_URL%g" manifests/kube-dashboard-ingress.yaml
  kubectl apply -Rf manifests/
}

if [ "$MODE" == "install" ]
then
  kubectl get deployment kubernetes-dashboard -n kube-system >/dev/null 2>&1 && FIRST_INSTALL="false"
  if [ "$FIRST_INSTALL" == "true" ]
  then
    install
  else
    echo "Deployment kubernetes-dashboard exists. Please, delete or run with the --upgrade option it to avoid shooting yourself in the foot."
  fi
elif [ "$MODE" == "upgrade" ]
then
  upgrade
elif [ "$MODE" == "delete" ]
then
  kubectl delete ing external-auth-oauth2 -n kube-system ||true
  kubectl delete ing oauth2-proxy -n kube-system ||true
  kubectl delete deployment oauth2-proxy -n kube-system ||true
  kubectl delete svc oauth2-proxy -n kube-system ||true
  kubectl delete clusterrole dashboard ||true
  kubectl delete sa kubernetes-dashboard -n kube-system || true
  kubectl delete clusterrolebinding kubernetes-dashboard ||true
  kubectl delete deployment kubernetes-dashboard -n kube-system ||true
  kubectl delete svc kubernetes-dashboard -n kube-system ||true
fi

function cleanup {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT
