#!/bin/bash
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.


if [ ! -f "hack/packaging/q-helm/bin/helm-linux-amd64" ]; then
  temp=$(mktemp -d)
	wget wget -q -O - https://get.helm.sh/helm-v3.11.3-linux-amd64.tar.gz | tar -xzf - -C "$temp"
	mv "$temp/linux-amd64/helm" hack/packaging/q-helm/bin/helm-linux-amd64
	rm -rf "$temp"
fi
if [ ! -f "hack/packaging/q-helm/bin/helm-linux-arm64" ]; then
  temp=$(mktemp -d)
	wget wget -q -O - https://get.helm.sh/helm-v3.11.3-linux-arm64.tar.gz | tar -xzf - -C "$temp"
	mv "$temp/linux-arm64/helm" hack/packaging/q-helm/bin/helm-linux-arm64
	rm -rf "$temp"
fi

chmod +x hack/packaging/q-helm/bin/helm-linux-amd64 hack/packaging/q-helm/bin/helm-linux-arm64
