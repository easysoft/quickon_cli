#!/usr/bin/env sh

set -xe

export APP_DOMAIN=${DOMAIN:-k3s.local}
export TOP_DOMAIN=${APP_DOMAIN#*.}

kubectl apply -f  https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml
kubectl apply -f  https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml -n cne-system
