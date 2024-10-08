#!/usr/bin/env sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

set -e

[ ! -z $DEBUG ] && set -x

export APP_DOMAIN=${DOMAIN:-k3s.local}
export TOP_DOMAIN=${APP_DOMAIN#*.}
export NS=${NS:-quickon-system}

kubectl apply -f  https://pkg.zentao.net/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml -n kube-system
kubectl apply -f  https://pkg.zentao.net/ssl/${TOP_DOMAIN}/${APP_DOMAIN}/tls.yaml -n $NS
