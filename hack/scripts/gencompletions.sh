#!/bin/sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

set -ex

rm -rf completions
mkdir completions
go build -o /tmp/qcadmin
/tmp/qcadmin version | grep -q linux && touch /tmp/.qcadmin-linux
for sh in bash zsh fish; do
	/tmp/qcadmin completion "$sh" >"completions/qcadmin.$sh"
  cp -a "completions/qcadmin.$sh" "completions/q.$sh"
  [ -f "/tmp/.qcadmin-linux" ] && (
    sed -i "s#qcadmin#q#g" "completions/q.$sh"
  ) || (
    sed -i "" "s#qcadmin#q#g" "completions/q.$sh"
  )
done
rm -rf /tmp/qcadmin /tmp/.qcadmin-linux
