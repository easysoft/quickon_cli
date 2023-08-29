#!/bin/sh
# Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
# Use of this source code is covered by the following dual licenses:
# (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
# (2) Affero General Public License 3.0 (AGPL 3.0)
# license that can be found in the LICENSE file.


[ $(id -u) -eq 0 ] || exec sudo $0 $@

for bin in /var/lib/rancher/k3s/data/**/bin/; do
    [ -d $bin ] && export PATH=$PATH:$bin:$bin/aux
done

for service in /etc/systemd/system/k3s*.service; do
    [ -s $service ] && systemctl stop $(basename $service)
    [ -s $service ] && systemctl disable $(basename $service)
done

for service in /etc/init.d/k3s*; do
    [ -x $service ] && $service stop
done

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

pschildren() {
    ps -e -o ppid= -o pid= | \
    sed -e 's/^\s*//g; s/\s\s*/\t/g;' | \
    grep -w "^$1" | \
    cut -f2
}

pstree() {
    for pid in $@; do
        echo $pid
        for child in $(pschildren $pid); do
            pstree $child
        done
    done
}

killtree() {
    kill -9 $(
        { set +x; } 2>/dev/null;
        pstree $@;
        set -x;
    ) 2>/dev/null
}

remove_interfaces() {
    # Delete network interface(s) that match 'master cni0'
    ip link show 2>/dev/null | grep 'master cni0' | while read ignore iface ignore; do
        iface=${iface%%@*}
        [ -z "$iface" ] || ip link delete $iface
    done

    # Delete cni related interfaces
    ip link delete cni0
    ip link delete flannel.1
    ip link delete flannel-v6.1
    ip link delete kube-ipvs0
    ip link delete flannel-wg
    ip link delete flannel-wg-v6
}

getshims() {
    ps -e -o pid= -o args= | sed -e 's/^ *//; s/\s\s*/\t/;' | grep -w 'data/[^/]*/bin/containerd-shim' | cut -f1
}

killtree $({ set +x; } 2>/dev/null; getshims; set -x)

do_unmount_and_remove() {
    # set +x
    while read -r _ path _; do
        case "$path" in $1*) echo "$path" ;; esac
    done < /proc/self/mounts | sort -r | xargs -r -t -n 1 sh -c 'umount "$0" && rm -rf "$0"'
    # set -x
}

do_unmount_and_remove '/run/k3s'
do_unmount_and_remove '/opt/quickon/platform'
do_unmount_and_remove '/var/lib/rancher/k3s'
do_unmount_and_remove '/var/lib/kubelet/pods'
do_unmount_and_remove '/var/lib/kubelet/plugins'
do_unmount_and_remove '/run/netns/cni-'

# Remove CNI namespaces
ip netns show 2>/dev/null | grep cni- | xargs -r -t -n 1 ip netns delete

remove_interfaces

rm -rf /var/lib/cni/

iptables-save | grep -v KUBE- | grep -v CNI- | grep -v flannel | iptables-restore
ip6tables-save | grep -v KUBE- | grep -v CNI- | grep -v flannel | ip6tables-restore

if command_exists docker && [ -e /var/run/docker.sock ]; then
		(
			rm_ctns=$(docker ps -a -q --filter 'name=k8s')
			if [ -z "$rm_ctns" ];then
    		echo "no containers need to delete"
			else
        docker rm -f $rm_ctns
			fi
		) || true
fi

exit 0
