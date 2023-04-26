#!/bin/sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

# Determine OS platform
# shellcheck source=/dev/null

ensure_sudo() {
	if [ "$(id -u)" = "0" ]; then
		echo "Sudo permissions detected"
	else
		echo "No sudo permission detected, please run as sudo"
		exit 1
	fi
}

{
  ensure_sudo
}
