#!/bin/bash

if [ ! -f "hack/bin/k3s-linux-amd64" ]; then
	wget -O hack/bin/k3s-linux-amd64 https://pkg.qucheng.com/qucheng/cli/stable/k3s/v1.23.7/k3s-linux-amd64
fi
if [ ! -f "hack/bin/k3s-linux-arm64" ]; then
  wget -O hack/bin/k3s-linux-arm64 https://pkg.qucheng.com/qucheng/cli/stable/k3s/v1.23.7/k3s-linux-arm64
fi

chmod +x hack/bin/k3s-linux-amd64 hack/bin/k3s-linux-arm64

# if [ ! -f "hack/bin/helm-linux-amd64" ]; then
# wget -O hack/bin/helm-linux-amd64  https://pkg.qucheng.com/qucheng/cli/stable/helm/helm-linux-amd64
# fi

# if [ ! -f "hack/bin/helm-linux-arm64" ]; then
# wget -O hack/bin/helm-linux-arm64  https://pkg.qucheng.com/qucheng/cli/stable/helm/helm-linux-arm64
# fi

# chmod +x hack/bin/helm-linux-amd64 hack/bin/helm-linux-arm64
