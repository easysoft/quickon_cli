#!/bin/bash
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.

set -xe

if [ ! -f "hack/bin/k3s-linux-amd64" ]; then
	wget -O hack/bin/k3s-linux-amd64 https://github.com/k3s-io/k3s/releases/download/v1.28.11%2Bk3s1/k3s
fi
if [ ! -f "hack/bin/k3s-linux-arm64" ]; then
  wget -O hack/bin/k3s-linux-arm64 https://github.com/k3s-io/k3s/releases/download/v1.28.11%2Bk3s1/k3s-arm64
fi

chmod +x hack/bin/k3s-linux-amd64 hack/bin/k3s-linux-arm64
