#!/bin/sh

[ $(id -u) -eq 0 ] || exec sudo $0 $@

echo "incluster mode"
if [ -f "/usr/local/bin/helm" ]; then
	helm delete cne-api -n cne-system
	helm delete qucheng -n cne-system
	helm repo list | grep -q install && helm repo remove install || true
fi
