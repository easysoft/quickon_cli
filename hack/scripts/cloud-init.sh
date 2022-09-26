#!/usr/bin/env sh

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

if ! command_exists k3s ; then
  echo "download k3s"
  wget https://ghproxy.com/https://github.com/k3s-io/k3s/releases/download/v1.24.6-rc1%2Bk3s1/k3s
  chmod +x k3s
  mv k3s /usr/local/bin/k3s
fi

if ! command_exists kubectl; then
  cp -a /usr/local/bin/k3s /usr/local/bin/kubectl
fi

cat > /etc/systemd/system/k3s.service <<EOF
[Unit]
Description=k3s server
ConditionFileIsExecutable=/usr/local/bin/k3s
After=network-online.target

[Service]
Type=notify
Type=process
Delegate=yes
EnvironmentFile=-/etc/sysconfig/%N
EnvironmentFile=-/etc/default/%N
EnvironmentFile=-/etc/systemd/system/k3s.service.env
StartLimitInterval=5
StartLimitBurst=10
ExecStartPre=-/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/k3s "server" "--docker" "--kubelet-arg=max-pods=220" "--kube-proxy-arg=proxy-mode=ipvs" "--kube-proxy-arg=masquerade-all=true" "--kube-proxy-arg=metrics-bind-address=0.0.0.0" "--data-dir=/opt/quickon/platform" "--pause-image=hub.qucheng.com/library/k3s-pause:3.6" "--disable-network-policy" "--disable-helm-controller" "--disable=servicelb,traefik" " --tls-san=kapi.qucheng.local" "--service-node-port-range=22767-32767" "--system-default-registry=hub.qucheng.com/library" "--cluster-cidr=10.42.0.0/16" "--service-cidr=10.43.0.0/16"
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
Restart=always
RestartSec=30
[Install]
WantedBy=multi-user.target
EOF

mkdir -pv /opt/quickon/platform/server/manifests /opt/quickon/backup
chmod 777 /opt/quickon/backup

[ -f "/opt/quickon/init/env" ] && source /opt/quickon/init/env

[ -z "$APP_DOMAIN" ] && export APP_DOMAIN=demo.haogs.cn

[ -z "$APP_HTTPS_PORT" ] && export APP_HTTPS_PORT=443

[ -z "$APP_HTTP_PORT" ] && export APP_HTTP_PORT=80

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
        image: hub.qucheng.com/platform/podinstall:2022092612
        imagePullPolicy: Always
        env:
        - name: APP_DOMAIN
          value: "$APP_DOMAIN"
        - name: APP_TOKEN
          value: "$APP_TOKEN"
        - name: APP_HTTP_PORT
          value: "$APP_HTTP_PORT"
        - name: APP_HTTPS_PORT
          value: "$APP_HTTPS_PORT"
        volumeMounts:
        - name: k3syaml
          mountPath: /qcli/k3syaml
        - name: qcfg
          mountPath: /qcli/root
      volumes:
      - name: k3syaml
        hostPath:
          path: /etc/rancher/k3s
      - name: qcfg
        hostPath:
          path: /root
      restartPolicy: Never
EOF

systemctl enable k3s --now
