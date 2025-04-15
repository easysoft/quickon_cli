---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: q-nfs-nfs-subdir-external-provisioner
  namespace: quickon-storage
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: q-nfs
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
provisioner: cluster.local/q-nfs-nfs-subdir-external-provisioner
allowVolumeExpansion: true
reclaimPolicy: Delete
volumeBindingMode: Immediate
parameters:
  archiveOnDelete: "true"
---
# Source: nfs-subdir-external-provisioner/templates/clusterrole.yaml
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: q-nfs-nfs-subdir-external-provisioner-runner
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "update", "patch"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: run-q-nfs-nfs-subdir-external-provisioner
subjects:
  - kind: ServiceAccount
    name: q-nfs-nfs-subdir-external-provisioner
    namespace: quickon-storage
roleRef:
  kind: ClusterRole
  name: q-nfs-nfs-subdir-external-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-q-nfs-nfs-subdir-external-provisioner
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-q-nfs-nfs-subdir-external-provisioner
subjects:
  - kind: ServiceAccount
    name: q-nfs-nfs-subdir-external-provisioner
    namespace: quickon-storage
roleRef:
  kind: Role
  name: leader-locking-q-nfs-nfs-subdir-external-provisioner
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: q-nfs-nfs-subdir-external-provisioner
  namespace: quickon-storage
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: nfs-subdir-external-provisioner
      release: q-nfs
  template:
    metadata:
      annotations:
      labels:
        app: nfs-subdir-external-provisioner
        release: q-nfs
    spec:
      serviceAccountName: q-nfs-nfs-subdir-external-provisioner
      securityContext:
        {}
      containers:
        - name: nfs-subdir-external-provisioner
          image: hub.zentao.net/app/nfs-subdir-external-provisioner:v4.0.2
          imagePullPolicy: IfNotPresent
          securityContext:
            {}
          volumeMounts:
            - name: nfs-subdir-external-provisioner-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: cluster.local/q-nfs-nfs-subdir-external-provisioner
            - name: NFS_SERVER
              value: __NFS_IP__
            - name: NFS_PATH
              value: __NFS_PATH__
      volumes:
        - name: nfs-subdir-external-provisioner-root
          nfs:
            server: __NFS_IP__
            path: __NFS_PATH__
