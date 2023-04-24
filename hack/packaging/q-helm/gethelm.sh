#!/bin/bash

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
