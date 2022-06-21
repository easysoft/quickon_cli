#!/bin/sh

[ $(id -u) -eq 0 ] || exec sudo $0 $@

echo "incluster mode"

if [ -f "/usr/local/bin/qcadmin" ]; then
  qcadmin experimental dns clean
  qcadmin experimental helm uninstall --name cne-api -n cne-system
  qcadmin experimental helm uninstall --name qucheng -n cne-system
  qcadmin experimental helm repo-del
fi
