#!/usr/bin/env bash
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.


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
