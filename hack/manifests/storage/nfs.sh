#!/usr/bin/env bash

kubectl get sc | grep q-nfs >/dev/null 2>&1 && exit 0

helm repo add install https://hub.qucheng.com/chartrepo/stable

helm repo update

helm upgrade -i q-nfs install/nfs-subdir-external-provisioner \
  -n quickon-storage \
  --set nfs.server=127.0.0.1 \
  --set nfs.path=/opt/quickon/storage/nfs \
  --set storageClass.defaultClass=true \
  --set storageClass.name=q-nfs
