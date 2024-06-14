#!/usr/bin/env bash

ip=${1:-127.0.0.1}

kubectl get sc | grep q-nfs >/dev/null 2>&1 && exit 0

helm repo add install https://hub.zentao.net/chartrepo/stable

helm repo update

helm upgrade -i q-nfs install/nfs-subdir-external-provisioner \
  -n quickon-storage \
  --set nfs.server=${ip} \
  --set nfs.path=/opt/quickon/storage/nfs \
  --set storageClass.defaultClass=true \
  --set storageClass.name=q-nfs
