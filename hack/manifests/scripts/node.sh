#!/bin/sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.


set -e
set -o noglob

# --- helper functions for logs ---
info()
{
    echo '[INFO]  ' "$@"
}
warn()
{
    echo '[WARN]  ' "$@" >&2
}
fatal()
{
    echo '[ERROR]  ' "$@" >&2
    exit 1
}

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

# --- add additional utility links ---
create_symlinks() {
    for cmd in kubectl crictl ctr; do
        if [ ! -e /usr/local/bin/${cmd} ]; then
            which_cmd=$(command -v ${cmd} 2>/dev/null || true)
            if [ -z "${which_cmd}" ]; then
                $SUDO ln -sf k3s /usr/local/bin/${cmd}
            fi
        fi
    done
}

check_docker() {
  which_cmd=$(command -v docker 2>/dev/null || true)
  if [ -z "${which_cmd}" ]; then
    sed -i "s#--docker \\\##g" /root/.k3s.service
    sed -i "s#--docker##g" /root/.k3s.service
  else
    docker_status=$(systemctl is-active docker)
    if [ "${docker_status}" != "active" ]; then
      sed -i "s#--docker \\\##g" /root/.k3s.service
      sed -i "s#--docker##g" /root/.k3s.service
    fi
  fi
}

# --- disable current service if loaded --
systemd_disable() {
    $SUDO systemctl disable k3s >/dev/null 2>&1 || true
    $SUDO rm -f /etc/systemd/system/k3s || true
    $SUDO rm -f /etc/systemd/system/k3s.env || true
}

# --- capture current env and create file containing k3s_ variables ---
# create_env_file() {
#     info "env: Creating environment file ${FILE_K3S_ENV}"
#     $SUDO touch /etc/systemd/system/k3s.service.env
#     $SUDO chmod 0600 /etc/systemd/system/k3s.service.env
#     sh -c export | while read x v; do echo $v; done | grep -E '^(K3S|CONTAINERD)_' | $SUDO tee ${FILE_K3S_ENV} >/dev/null
#     sh -c export | while read x v; do echo $v; done | grep -Ei '^(NO|HTTP|HTTPS)_PROXY' | $SUDO tee -a ${FILE_K3S_ENV} >/dev/null
# }

systemd_env() {
  info "init k3s tls env"
  cat > /tmp/.k3s.service.env <<EOF
CATTLE_NEW_SIGNED_CERT_EXPIRATION_DAYS=3650
EOF

  $SUDO mv /tmp/.k3s.service.env /etc/systemd/system/k3s.service.env
  $SUDO chmod 0777 /etc/systemd/system/k3s.service.env
}

# --- enable and start systemd service ---
systemd_enable() {
    info "enabling k3s unit"
    $SUDO cp /root/.k3s.service /etc/systemd/system/k3s.service
    $SUDO systemctl enable k3s >/dev/null
    $SUDO systemctl daemon-reload >/dev/null
}

systemd_start() {
    info "starting k3s"
    $SUDO systemctl restart k3s
}

# --- startup systemd or openrc service ---
service_enable_and_start() {
    if [ -f "/proc/cgroups" ] && [ "$(grep memory /proc/cgroups | while read -r n n n enabled; do echo $enabled; done)" -eq 0 ];
    then
        info 'Failed to find memory cgroup, you may need to add "cgroup_memory=1 cgroup_enable=memory" to your linux cmdline (/boot/cmdline.txt on a Raspberry Pi)'
    fi
    systemd_env
    systemd_enable
    systemd_start
    return 0
}

# --- run the install process --
{
    create_symlinks
    systemd_disable
    # check_docker
    # create_env_file
    service_enable_and_start
}
