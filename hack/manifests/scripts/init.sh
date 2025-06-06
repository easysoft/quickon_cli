#!/bin/sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.


[ -f "/.q.initdone" ] && exit 0

args=$1

echo "init system package..."

# --- install package
install_package() {

  if [ "$args" = "offline" ]; then
    echo "offline skip install package"
    return
  fi

  if type apt >/dev/null 2>&1; then
    export DEBIAN_FRONTEND=noninteractive
    apt update -qq
    apt remove -y -qq ufw lxd lxd-client lxcfs lxc-common
    apt install --no-install-recommends --no-install-suggests -y -qq nfs-common iptables conntrack jq socat bash-completion open-iscsi rsync ipset ipvsadm htop net-tools wget psmisc git curl nload ebtables ethtool procps
    systemctl enable --now iscsid
    # ufw disable
    # ufw allow 6443/tcp #apiserver
    # ufw allow from 10.42.0.0/16 to any #pods
    # ufw allow from 10.43.0.0/16 to any #services
  fi

  if type yum >/dev/null 2>&1; then
    yum install -y -q nfs-utils iptables conntrack jq socat bash-completion rsync ipset ipvsadm htop net-tools wget psmisc git curl nload ebtables ethtool
    yum --setopt=tsflags=noscripts install -y -q iscsi-initiator-utils
    systemctl enable --now iscsid
    systemctl disable firewalld || true
    systemctl stop firewalld || true
    systemctl mask firewalld || true
    #  systemctl disable NetworkManager
    #  systemctl stop NetworkManager
    status=$(getenforce)
    if [ "$status" != "Disabled" ]; then
      setenforce 0
      cp /etc/selinux/config /etc/selinux/config.bak
      sed -i 's/^SELINUX=.*/SELINUX=disabled/g' /etc/selinux/config
    fi
  fi
}


config_system() {

if command -v systemctl; then

mkdir -pv /etc/systemd/system.conf.d
cat > /etc/systemd/system.conf.d/30-k8s-ulimits.conf <<EOF
[Manager]
DefaultLimitCORE=infinity
DefaultLimitNOFILE=100000
DefaultLimitNPROC=100000
EOF

mkdir -pv /etc/systemd/journald.conf.d /var/log/journal

cat > /etc/systemd/journald.conf.d/95-k8s-journald.conf <<EOF
[Journal]
# 持久化保存到磁盘
Storage=persistent
# 最大占用空间 2G
SystemMaxUse=2G
# 单日志文件最大 100M
SystemMaxFileSize=100M
# 日志保存时间 1 周
MaxRetentionSec=1week
# 禁止转发
ForwardToSyslog=no
ForwardToWall=no
EOF

systemctl daemon-reload
systemctl restart systemd-journald

swapoff -a && sysctl -w vm.swappiness=0

cat > /etc/modules-load.d/10-k8s-modules.conf <<EOF
br_netfilter
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
nf_conntrack
EOF

systemctl daemon-reload
systemctl restart systemd-modules-load

fi

sed -i  's/^.*net.ip.*/# &/g' /etc/sysctl.conf

cat > /etc/sysctl.d/95-k8s-sysctl.conf <<EOF
# 转发
net.ipv4.ip_forward = 1
# 对直接连接的网络进行反向路径过滤
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
#不允许接受含有源路由信息的ip包
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0
#打开TCP SYN cookies保护, 一定程度预防SYN攻击
net.ipv4.tcp_syncookies = 1
#SYN队列的长度,适当增大该值,有助于抵挡SYN攻击
net.ipv4.tcp_max_syn_backlog = 3072
#SYN的重试次数,适当降低该值,有助于防范SYN攻击
net.ipv4.tcp_synack_retries = 3
net.ipv4.tcp_syn_retries = 3
#关闭Linux kernel的路由重定向功能
# net.ipv4.conf.all.send_redirects = 0
# net.ipv4.conf.default.send_redirects = 0
#不允许ip重定向信息
# net.ipv4.conf.all.accept_redirects = 0
#取消安全重定向
# net.ipv4.conf.all.secure_redirects = 0
# icmp ping
# net.ipv4.icmp_echo_ignore_broadcasts = 1
# net.ipv4.icmp_ignore_bogus_error_responses = 1
#进程快速回收,避免系统中存在大量TIME_WAIT进程
# net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_fin_timeout = 30 # 缩短TIME_WAIT时间,加速端口回收
#临时端口范围
net.ipv4.ip_local_port_range = 20000 65535
#预留给kubernetes service的nodeport端口范围,不设置可能会造成
#kubernetes在做服务探针时使用下列范围端口,造成连接被占用而失败,引起探针失效
net.ipv4.ip_local_reserved_ports = 30000-32768
# 以下三个参数是 arp 缓存的 gc 阀值,相比默认值提高了,避免在某些场景下arp缓存溢出导致网络超时,参考: https://imroc.cc/k8s/troubleshooting/arp-cache-overflow-causes-healthcheck-failed/
# net.ipv4.neigh.default.gc_thresh1 = 2048
# net.ipv4.neigh.default.gc_thresh2 = 4096
# net.ipv4.neigh.default.gc_thresh3 = 8192
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 30
net.ipv4.tcp_keepalive_probes = 10
net.core.somaxconn = 65535
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-arptables = 1

# conntrack优化
net.netfilter.nf_conntrack_tcp_be_liberal = 1 # 容器环境下, 开启这个参数可以避免 NAT 过的 TCP 连接 带宽上不去。
net.netfilter.nf_conntrack_tcp_loose = 1
net.netfilter.nf_conntrack_max = 3200000
net.netfilter.nf_conntrack_buckets = 1600512
net.netfilter.nf_conntrack_tcp_timeout_time_wait = 30

# fd优化
fs.file-max = 6553600 # 提升文件句柄上限，像 nginx 这种代理，每个连接实际分别会对 downstream 和 upstream 占用一个句柄，连接量大的情况下句柄消耗就大。
fs.inotify.max_user_instances = 524288 # 表示同一用户同时最大可以拥有的 inotify 实例 (每个实例可以有很多 watch)
fs.inotify.max_user_watches = 524288 # 表示同一用户同时可以添加的watch数目（watch一般是针对目录，决定了同时同一用户可以监控的目录数量) 默认值 8192 在容器场景下偏小，在某些情况下可能会导致 inotify watch 数量耗尽，使得创建 Pod 不成功或者 kubelet 无法启动成功，将其优化到 524288
#如果wattch数过多可以打开 inotify_add_watch 跟踪，进一步 debug inotify watch 耗尽的原因:
#echo 1 >> /sys/kernel/debug/tracing/events/syscalls/sys_exit_inotify_add_watch/enable

vm.swappiness = 0
vm.max_map_count = 655360
EOF

sysctl -p /etc/sysctl.d/95-k8s-sysctl.conf

cat /etc/security/limits.conf | grep -vE "(^#|^$)" | wc | grep 0 && (
	cat > /etc/security/limits.conf <<EOF
* soft nofile 1000000
* hard nofile 1000000
* soft stack 10240
* soft nproc 65536
* hard nproc 65536
EOF
)

}

{
  install_package
  config_system
  touch /.q.initdone
  exit 0
}
