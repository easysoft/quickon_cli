#!/usr/bin/env bash
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.


[ -n "${DEBUG:+1}" ] && set -x

helm repo add install https://hub.zentao.net/chartrepo/stable

helm repo update

kubectl create ns quickon-system
kubectl create ns quickon-app
kubectl create ns quickon-ci

export APP_DOMAIN=${APP_DOMAIN:-k3s.local}
export APP_DOMAIN=$(echo -e $APP_DOMAIN | sed 's/[[:space:]]//g' | xargs -I {} echo {})
export APP_TOKEN=${APP_TOKEN:-XZdrjxhAhq5pDjpEU3kR4djsvJ3rfj0M}
export TOP_DOMAIN=${APP_DOMAIN#*.}

# add debug

echo "\"$APP_DOMAIN\" \"$TOP_DOMAIN\""

cat > /tmp/qucheng.yaml <<EOF
cloud:
  defaultChannel: stable
  apphttpsPort: ${APP_HTTPS_PORT}
env:
  APP_DOMAIN: ${APP_DOMAIN}
  CNE_API_TOKEN: ${APP_TOKEN}
ingress:
  host: ${APP_DOMAIN}
EOF

cat > /tmp/operator.yaml <<EOF
minio:
  auth:
    password: Z6Ho2LdLZb8AAXuv
    username: FaBux6M2
  ingress:
    enabled: true
    host: s3.${APP_DOMAIN}
EOF

helm upgrade -i ingress install/nginx-ingress-controller -n quickon-system
helm upgrade -i cne-operator install/cne-operator -f /tmp/operator.yaml -n quickon-system
helm upgrade -i qucheng install/qucheng -f /tmp/qucheng.yaml -n quickon-system

[ -d "/qcli/root/.kube" ] || mkdir -pv /qcli/root/.kube
[ -d "/qcli/root/.qc/config" ] || mkdir -pv /qcli/root/.qc/config

cp -a /qcadmin_linux_amd64 /qcli/qbin/q
cp -a /qcadmin_linux_amd64 /qcli/qbin/qcadmin
cp -a /usr/local/bin/kubectl /qcli/qbin/kubectl
cp -a /usr/local/bin/helm /qcli/qbin/helm

cat > /qcli/root/.qc/config/cluster.yaml <<EOF
api_token: ${APP_TOKEN}
cluster:
  cid: b31b9178-ca9f-4acd-b9a1-c5277d631fe2
  cni: flannel
  init_node: ${APP_NODE_IP}
  master:
  - host: ${APP_NODE_IP}
    init: true
  pod-cidr: 10.42.0.0/16
  registry: hub.zentao.net
  svc-cidr: 10.43.0.0/16
  token: YywCEEPKVhaDEgF4
  worker: null
console-password: pass4Quickon
datadir: /opt/quickon
db: ""
domain: ${APP_DOMAIN}
global:
  ssh: {}
install:
  pkg: ""
  type: online
quickon:
  type: oss
storage:
  type: local
s3:
  password: Z6Ho2LdLZb8AAXuv
  username: FaBux6M2
EOF

[ -f "/qcli/k3syaml/k3s.yaml" ] && cp -a /qcli/k3syaml/k3s.yaml /qcli/root/.kube/config

# wait tls ok
wait_for_tls() {
    local retries=${MAXWAIT:-300}
    echo "check the tls is available."
    for ((i = 1; i <= $retries; i += 1)); do
        code=$(curl -s -o /dev/null -w "%{http_code}" https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml)
        if [ $code != 404 ] ;
        then
            echo "tls is ready."
            break
        fi

        echo "Waiting tls gen $i seconds"
        sleep 1
        if [ "$i" == "$retries" ]; then
            echo "unable to load tls"
            return 1
        fi
    done
    return 0
}

wait_for_tls && (
  kubectl apply -f https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml -n default
  kubectl apply -f https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml -n quickon-app
  kubectl apply -f https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml -n quickon-system
  kubectl apply -f https://pkg.qucheng.com/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml -n kube-system
) || echo "load tls failed"
