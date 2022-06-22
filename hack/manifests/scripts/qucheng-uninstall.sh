#!/bin/sh

[ $(id -u) -eq 0 ] || exec sudo $0 $@

echo "incluster mode"

qcmd=${1:-"/usr/local/bin/qcadmin"}

if [ -f "${qcmd}" ]; then
  echo "${qcmd} clean helm"
  ${qcmd} experimental dns clean
  ${qcmd} experimental helm uninstall --name cne-api -n cne-system
  ${qcmd} experimental helm uninstall --name qucheng -n cne-system
  ${qcmd} experimental helm repo-del
fi
