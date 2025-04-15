#!/usr/bin/env bash

nfsIP=${1:-127.0.0.1}
nfsPath=${2:-/opt/quickon/storage/nfs}

sourcePath=$(dirname "$0")/nfs.yaml.tpl
targetPath=$(dirname "$0")/nfs.deploy.yaml

cp -a "${sourcePath}" "${targetPath}"

sed -i "s|__NFS_IP__|${nfsIP}|g" "${targetPath}"
sed -i "s|__NFS_PATH__|${nfsPath}|g" "${targetPath}"
