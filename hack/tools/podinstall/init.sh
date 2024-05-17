#!/bin/bash
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.


[ -f "/.q.init" ] && exit 0

[ -d "/opt/quickon/init" ] || mkdir -pv /opt/quickon/init

mkdir -pv /opt/quickon/platform/server/manifests /opt/quickon/backup

chmod 777 /opt/quickon/backup

wait_for_mount() {
    local retries=${MAXWAIT:-300}
    echo "Check whether the Mount is available."

    for ((i = 1; i <= $retries; i += 1)); do
        if [ -f "/opt/quickon/init/env" ] ;
        then
            echo "Mount is ready."
            break
        fi

        if [ -f "/mnt/env" ] ;
        then
            cp -a /mnt/env /opt/quickon/init/env
        fi

        echo "Waiting Mount $i seconds"
        sleep 1

        if [ "$i" == "$retries" ]; then
            echo "Unable to connect to mount"
            return 1
        fi
    done
    return 0
}

wait_for_mount && source /opt/quickon/init/env || echo "mount env failed"

[ -z "$QUICKON_DOMAIN" ] && export QUICKON_DOMAIN=demo.corp.cc

[ -z "$QUICKON_HTTP_PORT" ] && export QUICKON_HTTP_PORT=443

[ -z "$QUICKON_HTTPS_PORT" ] && export QUICKON_HTTPS_PORT=80

[ -z "$APP_TOKEN" ] && export APP_TOKEN=$(pwgen 30 1)

cat > /opt/quickon/platform/server/manifests/initcluster.yaml <<EOF
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: qcli-cm
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: qcli-cm-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: qcli-cm
  namespace: kube-system
---
apiVersion: batch/v1
kind: Job
metadata:
  name: qcli-cm
  namespace: kube-system
  labels:
    app: qcli-cm
spec:
  backoffLimit: 1
  template:
    metadata:
      name: qcli-cm
      labels:
        app: qcli-cm
    spec:
      serviceAccountName: qcli-cm
      containers:
      - name: qcli-cm
        image: hub.zentao.net/platform/podinstall:2022102713
        imagePullPolicy: Always
        env:
        - name: APP_NODE_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        - name: APP_DOMAIN
          value: "$QUICKON_DOMAIN"
        - name: APP_TOKEN
          value: "$APP_TOKEN"
        - name: APP_HTTP_PORT
          value: "$QUICKON_HTTP_PORT"
        - name: APP_HTTPS_PORT
          value: "$QUICKON_HTTPS_PORT"
        volumeMounts:
        - name: k3syaml
          mountPath: /qcli/k3syaml
        - name: qcfg
          mountPath: /qcli/root
        - name: qbin
          mountPath: /qcli/qbin
      volumes:
      - name: k3syaml
        hostPath:
          path: /etc/rancher/k3s
      - name: qcfg
        hostPath:
          path: /root
      - name: qbin
        hostPath:
          path: /usr/local/bin
      restartPolicy: Never
EOF

systemctl enable k3s --now

touch /.q.init

exit 0
