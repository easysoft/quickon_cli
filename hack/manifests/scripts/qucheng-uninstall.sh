#!/usr/bin/env bash

# shellcheck disable=SC2068
[ $(id -u) -eq 0 ] || exec sudo "$0" $@

echo "clean quickon"

# shellcheck disable=SC2034
qcli=${1:-"/usr/local/bin/qcadmin"}

if [ -f "${qcmd}" ]; then
  echo "clean domain"
  ${qcli} experimental tools domain clean
  echo "clean helm app"
  ${qcli} experimental helm clean
fi
