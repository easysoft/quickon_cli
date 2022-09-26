#!/usr/bin/env sh

helm repo add install https://hub.qucheng.com/chartrepo/stable

helm repo update

kubectl create ns cne-system

export TOP_DOMAIN=${APP_DOMAIN#*.}

kubectl apply -f  https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml

cat > /tmp/qucheng.yaml <<EOF
cloud:
  defaultChannel: stable
env:
  APP_DOMAIN: ${APP_DOMAIN:-k3s.local}
  CNE_API_TOKEN: ${APP_TOKEN:-XZdrjxhAhq5pDjpEU3kR4djsvJ3rfj0M}
ingress:
  host: ${APP_DOMAIN:-k3s.local}
EOF

helm upgrade -i ingress install/nginx-ingress-controller -n cne-system
helm upgrade -i cne-operator install/cne-operator -n cne-system
helm upgrade -i qucheng install/qucheng -f /tmp/qucheng.yaml -n cne-system

[ -d "/qcli/root/.kube" ] || mkdir -pv /qcli/root/.kube
[ -f "/qcli/k3syaml/k3s.yaml" ] && cp -a /qcli/k3syaml/k3s.yaml /qcli/root/.kube/config

