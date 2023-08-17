#!/bin/sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

# Determine OS platform
# shellcheck source=/dev/null

create_link() {
  # shellcheck disable=SC2015
	[ -f "/usr/bin/q" ] && rm -rf /usr/bin/q || true
  # shellcheck disable=SC2015
  [ -f "/usr/bin/z" ] && rm -rf /usr/bin/z || true
	# shellcheck disable=SC2015
	[ -f "/usr/local/bin/q" ] && rm -rf /usr/local/bin/q || true
	# shellcheck disable=SC2015
	[ -f "/usr/local/bin/z" ] && rm -rf /usr/local/bin/z || true
	# shellcheck disable=SC2015
	[ -f "/usr/local/bin/qcadmin" ] && rm -rf /usr/local/bin/qcadmin || true
  # shellcheck disable=SC2015
  [ -f "/usr/bin/qcadmin" ] && ln -s /usr/bin/qcadmin /usr/bin/q || true
  # shellcheck disable=SC2015
  [ -f "/usr/bin/qcadmin" ] && ln -s /usr/bin/qcadmin /usr/bin/z || true

}

summary() {
	echo "----------------------------------------------------------------------"
	echo "Zentao(z)/Quickon(q) package has been successfully installed."
	echo ""
	echo " Please follow the next steps to start the software:"
	echo ""
	echo "    z --help"
	echo ""
	echo "----------------------------------------------------------------------"
}

{
  create_link
  summary
}
