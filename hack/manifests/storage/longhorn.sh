#!/usr/bin/env bash
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

z exp helm repo-list | grep install || (
  z exp helm repo-add --name install --url https://hub.zentao.net/chartrepo/stable
)
z exp helm repo-update

# z helm upgrade -i longhorn install/longhorn -n quickon-storage --create-namespace --set ingress.host=lh.local

z exp helm upgrade --repo install --name longhorn -n quickon-storage --chart longhorn --set ingress.host=lh.local
