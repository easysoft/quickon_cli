#!/usr/bin/env bash

if type apt >/dev/null 2>&1; then
  apt update
  apt install -y nfs-common nfs-kernel-server
  systemctl enable rpcbind --now
  systemctl enable nfs-server --now
elif type yum >/dev/null 2>&1; then
  yum install -y nfs-utils
  systemctl enable rpcbind --now
  systemctl enable nfs --now
else
  echo "Neither apt-get nor yum found" >&2
  exit 1
fi

[ -f "/etc/exports" ] && cp -a /etc/exports /etc/exports.bak

SPATH=${1:-/opt/quickon/storage/nfs}

[ -d "$SPATH" ] || mkdir -p $SPATH

chmod 777 $SPATH

echo "$SPATH *(insecure,rw,sync,no_root_squash,no_subtree_check)" > /etc/exports

exportfs -r

showmount -e 127.0.0.1

touch $SPATH/.quickon
