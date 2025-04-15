apiVersion: v1
kind: ServiceAccount
metadata:
  name: q-local-provisioner-service-account
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: q-local-provisioner-role
rules:
- apiGroups: [""]
  resources: ["nodes", "persistentvolumeclaims", "configmaps", "pods/log"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["endpoints", "persistentvolumes", "pods"]
  verbs: ["*"]
- apiGroups: [""]
  resources: ["events"]
  verbs: ["create", "patch"]
- apiGroups: ["storage.k8s.io"]
  resources: ["storageclasses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: q-local-provisioner-bind
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: q-local-provisioner-role
subjects:
- kind: ServiceAccount
  name: q-local-provisioner-service-account
  namespace: kube-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: q-local-provisioner
  namespace: kube-system
spec:
  revisionHistoryLimit: 0
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
  selector:
    matchLabels:
      app: q-local-provisioner
  template:
    metadata:
      labels:
        app: q-local-provisioner
    spec:
      priorityClassName: "system-node-critical"
      serviceAccountName: q-local-provisioner-service-account
      tolerations:
          - key: "CriticalAddonsOnly"
            operator: "Exists"
          - key: "node-role.kubernetes.io/control-plane"
            operator: "Exists"
            effect: "NoSchedule"
          - key: "node-role.kubernetes.io/master"
            operator: "Exists"
            effect: "NoSchedule"
      containers:
      - name: q-local-provisioner
        image: "hub.zentao.net/rancher/local-path-provisioner:v0.0.28"
        imagePullPolicy: IfNotPresent
        command:
        - local-path-provisioner
        - start
        - --config
        - /etc/config/config.json
        - --configmap-name
        - q-local-config
        - --service-account-name
        - q-local-provisioner-service-account
        volumeMounts:
        - name: config-volume
          mountPath: /etc/config/
        env:
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
      volumes:
        - name: config-volume
          configMap:
            name: q-local-config
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: q-local
  annotations:
    defaultVolumeType: "local"
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: rancher.io/local-path
volumeBindingMode: WaitForFirstConsumer
reclaimPolicy: Delete
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: q-local-config
  namespace: kube-system
data:
  config.json: |-
    {
      "nodePathMap":[
      {
        "node":"DEFAULT_PATH_FOR_NON_LISTED_NODES",
        "paths":["__LOCAL_PATH__"]
      }
      ]
    }
  setup: |-
    #!/bin/sh
    set -eu
    mkdir -m 0777 -p "${VOL_DIR}"
    chmod 700 "${VOL_DIR}/.."
  teardown: |-
    #!/bin/sh
    set -eu
    rm -rf "${VOL_DIR}"
  helperPod.yaml: |-
    apiVersion: v1
    kind: Pod
    metadata:
      name: helper-pod
    spec:
      containers:
      - name: helper-pod
        image: "hub.zentao.net/rancher/mirrored-library-busybox:1.36.1"
        imagePullPolicy: IfNotPresent
