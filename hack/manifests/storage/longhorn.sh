#!/usr/bin/env bash
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

qcadmin exp helm repo-list | grep install || (
  qcadmin exp helm repo-add --name install --url https://hub.zentao.net/chartrepo/stable
)
qcadmin exp helm repo-update

helm upgrade -i longhorn q-stable/longhorn -n quickon-storage --create-namespace
