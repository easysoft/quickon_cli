#!/bin/bash

set -xe

if [ ! -f "hack/bin/k3s-linux-amd64" ]; then
	wget -O hack/bin/k3s-linux-amd64 https://github.com/k3s-io/k3s/releases/download/v1.24.9%2Bk3s2/k3s
fi
if [ ! -f "hack/bin/k3s-linux-arm64" ]; then
  wget -O hack/bin/k3s-linux-arm64 https://github.com/k3s-io/k3s/releases/download/v1.24.9%2Bk3s2/k3s-arm64
fi

chmod +x hack/bin/k3s-linux-amd64 hack/bin/k3s-linux-arm64
