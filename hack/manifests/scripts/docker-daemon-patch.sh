#!/bin/sh

# shellcheck disable=SC1090,SC1091,SC2046,SC2068,SC2086

[ $(id -u) -eq 0 ] || exec sudo $0 $@

cat > /etc/docker/daemon.example.json <<EOF
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://dyucrs4l.mirror.aliyuncs.com"
  ],
  "exec-opts": [
    "native.cgroupdriver=cgroupfs"
  ],
  "bip": "169.254.123.1/24",
  "max-concurrent-downloads": 10,
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-size": "30m",
    "max-file": "3"
  },
  "storage-driver": "overlay2"
}
EOF


# Change cgroup to cgroupfs because k3s does not use systemd cgroup
if [ -f "/etc/docker/daemon.json" ]; then
  cp -a /etc/docker/daemon.json /etc/docker/daemon.json.bak
fi
# echo -e '{\n  "exec-opts": ["native.cgroupdriver=cgroupfs"]\n}' | tee /etc/docker/daemon.json
cp /etc/docker/daemon.example.json /etc/docker/daemon.json
systemctl daemon-reload
systemctl restart docker
