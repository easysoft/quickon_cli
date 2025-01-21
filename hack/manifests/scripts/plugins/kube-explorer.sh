#!/bin/sh
set -e
set -o noglob

# --- helper functions for logs ---
info()
{
    echo '[INFO] ' "$@"
}
warn()
{
    echo '[WARN] ' "$@" >&2
}
fatal()
{
    echo '[ERROR] ' "$@" >&2
    exit 1
}

verify_system() {
    # --- use sudo if we are not already root ---
    SUDO=sudo
    if [ $(id -u) -eq 0 ]; then
        SUDO=
    fi
    if [ ! -x /usr/local/bin/devops-kube-explorer ]; then
        $SUDO cp -a /root/.qc/bin/qc-explorer /usr/local/bin/devops-kube-explorer
        $SUDO chmod +x /usr/local/bin/devops-kube-explorer
    fi
}

# --- disable current service if loaded --
systemd_disable() {
    $SUDO systemctl disable devops-kube-explorer >/dev/null 2>&1 || true
    $SUDO rm -f /etc/systemd/system/devops-kube-explorer.service || true
    $SUDO rm -f /etc/systemd/system/devops-kube-explorer.service.env || true
}

# --- write systemd service file ---
create_service_file() {
  info "systemd: Creating service file /etc/systemd/system/devops-kube-explorer.service"
    $SUDO tee /etc/systemd/system/devops-kube-explorer.service >/dev/null << EOF
[Unit]
Description=devops plugin kube-explorer
Documentation=https://www.zentao.net
Wants=network-online.target
After=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
Type=simple
EnvironmentFile=-/etc/default/%N
EnvironmentFile=-/etc/sysconfig/%N
EnvironmentFile=-/etc/systemd/system/devops-kube-explorer.service.env
User=root
Restart=on-failure
RestartSec=5s
DynamicUser=true
ExecStart=/usr/local/bin/devops-kube-explorer \
            --http-listen-port=2525 \
            --https-listen-port=0 \
            --pod-image=hub.zentao.net/app/shell:v0.2.1-rc.7

EOF

    $SUDO restorecon -R -i /etc/systemd/system/devops-kube-explorer.service 2>/dev/null || true
}

# --- startup systemd service ---
service_enable_and_start() {
    info "systemd: Enabling devops-kube-explorer unit"
    $SUDO systemctl enable devops-kube-explorer >/dev/null
    $SUDO systemctl daemon-reload >/dev/null
    info "systemd: Starting devops plugin kube-explorer"
    $SUDO systemctl restart devops-kube-explorer
}

# --- run the install process --
{
    verify_system
    systemd_disable
    create_service_file
    service_enable_and_start
}
