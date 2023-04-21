#!/usr/bin/env bash

# shellcheck disable=SC2068
[ $(id -u) -eq 0 ] || exec sudo "$0" $@

current_dir=$(pwd)
tmpdir=/tmp
timestamp=$(date +%s)
diagnose_dir=/tmp/diagnose_${timestamp}
mkdir -p $diagnose_dir/{k3s,docker,quickon,system,k8s}
is_ps_hang=false

run() {
    echo
    # shellcheck disable=SC2145
    echo "-----------------run $@------------------"
    timeout 10s $@
    # shellcheck disable=SC2181
    if [ "$?" != "0" ]; then
        # shellcheck disable=SC2145
        echo "failed to collect info: $@"
    fi
    echo "------------End of ${1}----------------"
}

os_env() {
  run env | tee $diagnose_dir/system/os_env
}

get_distribution() {
    cat /etc/issue*
    if [ -r /etc/os-release ]; then
      cat /etc/os-release
    fi
    if [ -r /etc/redhat-release ]; then
      cat  /etc/redhat-release
    fi
}

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

# Get Service Status
service_status() {
    run service firewalld status | tee $diagnose_dir/system/service_status
    run service ntpd status | tee $diagnose_dir/system/service_status
    run service chronyd status | tee $diagnose_dir/system/service_status
}

# get System Info
system_info() {
    run uname -a | tee -a ${diagnose_dir}/system/os_info
    run uname -r | tee -a ${diagnose_dir}/system/os_info
    run get_distribution | tee -a ${diagnose_dir}/system/os_info
    if command_exists lsb_release; then
        run lsb_release | tee -a ${diagnose_dir}/system/os_info
    fi
    run ulimit -a | tee -a ${diagnose_dir}/system/ulimit
    run sysctl -a | tee -a ${diagnose_dir}/system/sysctl
    run cat /proc/vmstat | tee -a ${diagnose_dir}/system/vmstat
}

# Get Network Info
network_info() {
    # mkdir -p ${diagnose_dir}/network_info
    #run ifconfig
    run ip --details ad show | tee -a ${diagnose_dir}/system/network_ipad_info
    run ip --details link show | tee -a ${diagnose_dir}/system/network_iplink_info
    run ip route show | tee -a ${diagnose_dir}/system/network_route_info
    run cat /proc/net/nf_conntrack | tee -a ${diagnose_dir}/system/network_nf_conntrack_info
    run netstat -nt | tee -a ${diagnose_dir}/system/netstat_info
    run netstat -nu | tee -a ${diagnose_dir}/system/netstat_info
    run netstat -ln | tee -a ${diagnose_dir}/system/netstat_info
    if command_exists iptables-save; then
        run iptables-save | tee -a ${diagnose_dir}/system/iptables_info
    fi
    if command_exists ipvsadm; then
        run ipvsadm -L | tee -a ${diagnose_dir}/system/ipvsadm_list
    fi
}

memory_info() {
    run cat /proc/meminfo | tee -a ${diagnose_dir}/system/memory_info
    run cat /proc/buddyinfo | tee -a ${diagnose_dir}/system/memory_info
    run cat /proc/vmallocinfo | tee -a ${diagnose_dir}/system/memory_info
    run cat /proc/slabinfo | tee -a ${diagnose_dir}/system/memory_info
    run cat /proc/zoneinfo | tee -a ${diagnose_dir}/system/memory_info
}


# check ps -ef command is hung
check_ps_hang() {
  echo "check if ps -ef command hang" | tee -a ${diagnose_dir}/system/ps_command_status
  checkD=$(timeout -s 9 2 ps -ef)
  if [ "$?" != "0" ]; then
      echo "ps -ef command is hung" | tee -a ${diagnose_dir}/system/ps_command_status
      is_ps_hang=true
      echo "start to check which process lead to ps -ef command hang" | tee -a ${diagnose_dir}/system/ps_command_status
      for f in `find /proc/*/task -name status`
      do
          checkD=$(cat $f|grep "State.*D")
          if [ "$?" == "0" ]; then
              cmdline=$(echo ${f%%status}"cmdline")
              pid=$(echo ${f%%status}"")
              stack=$(echo ${f%%status}"stack")
              timeout -s 9 2 cat $cmdline
              if [ "$?" != "0" ]; then
                  echo "process $pid is in State D and lead to ps -ef process hang,stack info:" | tee -a ${diagnose_dir}/system/ps_command_status
                  cat $stack | tee -a ${diagnose_dir}/system/ps_command_status
              fi
          fi
      done
      echo "finish to check which process lead to ps -ef command hang" | tee -a ${diagnose_dir}/system/ps_command_status
  else
      echo "ps -ef command works fine" | tee -a ${diagnose_dir}/system/ps_command_status
  fi
}

#system status
system_status() {
    #mkdir -p ${diagnose_dir}/system_status
    run uptime | tee -a ${diagnose_dir}/system/uptime_status
    run top -b -n 1 | tee -a ${diagnose_dir}/system/top_status
    if [ "$is_ps_hang" == "false" ]; then
        run ps -ef | tee -a ${diagnose_dir}/system/ps-ef_status
    else
        echo "ps -ef command hang, skip [ps -ef] check" | tee -a ${diagnose_dir}/system/ps-hang_status
    fi

    run sar -A | tee -a ${diagnose_dir}/system/sar_status

    run df -h | tee -a ${diagnose_dir}/system/df_status

    run cat cat /etc/resolv.conf | tee -a ${diagnose_dir}/system/dns

    run cat cat /etc/hosts | tee -a ${diagnose_dir}/system/hosts

    run cat /proc/mounts | tee -a ${diagnose_dir}/system/mounts_status

    if [ "$is_ps_hang" == "false" ]; then
        run pstree -al | tee -a ${diagnose_dir}/system/pstree_status
    else
        echo "ps -ef command hang, skip [pstree -al] check" | tee -a ${diagnose_dir}/ps-hang_status
    fi

    run lsof | tee -a ${diagnose_dir}/system/lsof_status

    (
        cd /proc
        find -maxdepth 1 -type d -name '[0-9]*' \
         -exec bash -c "ls {}/fd/ | wc -l | tr '\n' ' '" \; \
         -printf "fds (PID = %P), command: " \
         -exec bash -c "tr '\0' ' ' < {}/cmdline" \; \
         -exec echo \; | sort -rn | head | tee -a ${diagnose_dir}/system/proc_status
    )

    echo "----------------start pid leak detect---------------------" | tee -a ${diagnose_dir}/system/ps_status
    ps -elT | awk '{print $4}' | sort | uniq -c | sort -k 1 -g | tail -5 | tee -a ${diagnose_dir}/system/ps_status
    echo "----------------done pid leak detect---------------------" | tee -a ${diagnose_dir}/system/ps_status
}


daemon_status() {
    if command_exists docker; then
        run systemctl status docker -l | tee -a ${diagnose_dir}/docker/service_status
        run systemctl cat docker | tee -a ${diagnose_dir}/docker/service_status
        docker_check
    fi
    if command_exists kubelet; then
        run systemctl status kubelet -l | tee -a ${diagnose_dir}/k8s/kubelet_status
    fi
    if command_exists containerd; then
        run systemctl status containerd -l | tee -a ${diagnose_dir}/k8s/containerd_status
    fi
    if command_exists k3s; then
        run systemctl status k3s -l | tee -a ${diagnose_dir}/k3s/service_status
        run systemctl cat k3s | tee -a ${diagnose_dir}/k3s/service_status
        run k3s check-config | tee $diagnose_dir/k3s/k3s.check.log
        k3s_check
    fi
}

k3s_check() {
  log_tail_lines=10000
  pidof systemd && journalctl -n ${log_tail_lines} -u k3s.service &> ${diagnose_dir}/k3s/k3s_systemd.log
}

docker_check() {
    log_tail_lines=10000
    pidof systemd && journalctl -n ${log_tail_lines} -u docker.service &> ${diagnose_dir}/docker/docker.log || tail -n ${log_tail_lines} /var/log/upstart/docker.log &> ${diagnose_dir}/docker/docker.log
    echo "check dockerd process"
    if [ "$is_ps_hang" == "false" ]; then
        run ps -ef|grep -E 'dockerd|docker daemon'|grep -v grep| tee -a ${diagnose_dir}/docker/docker_status
    else
        echo "ps -ef command hang, skip [ps -ef|grep -E 'dockerd|docker daemon'] check" | tee -a ${diagnose_dir}/docker/docker_status
    fi

    #docker info
    run docker info | tee -a ${diagnose_dir}/docker/docker_info
    run docker version | tee -a ${diagnose_dir}/docker/docker_version
    sudo kill -SIGUSR1 $(cat /var/run/docker.pid)
    [ -f "/var/run/docker/libcontainerd/containerd/events.log" ] && (
      cp /var/run/docker/libcontainerd/containerd/events.log ${diagnose_dir}/docker/containerd_events.log
    )
    sleep 10
    # shellcheck disable=SC2045
    for i in $(ls /var/log/upstart/docker.log*); do
        cp $i ${diagnose_dir}/docker/
    done
}

showlog() {
    local file=$1
    if [ -f "$file" ]; then
        tail -n 200 $file
    fi
}

# collect log
common_logs() {
    log_tail_lines=10000
    mkdir -p ${diagnose_dir}/logs
    run dmesg -T | tail -n ${log_tail_lines}  | tee ${diagnose_dir}/system/dmesg.log
    tail -c 500M /var/log/messages &> ${diagnose_dir}/system/messages
}

archive() {
    tar -zcvf $tmpdir/diagnose_${timestamp}.tar.gz ${diagnose_dir}
    echo "please get $tmpdir/diagnose_${timestamp}.tar.gz for diagnostics"
}

varlogmessage(){
    grep cloud-init /var/log/messages > $diagnose_dir/system/cloud-init.log
}

cluster_dump(){
    run kubectl cluster-info dump | tee $diagnose_dir/k8s/cluster-info_dump.log
}

cluster_events(){
    run kubectl get events | tee $diagnose_dir/k8s/cluster_events.log
}

component() {
    local ns="$1"
    mkdir -p $diagnose_dir/k8s/$ns/
    local pods=`kubectl get -n $ns po  | awk '{print $1}'|grep -v NAME`
    for po in ${pods}
    do
        kubectl logs -n ${ns} ${po} &> $diagnose_dir/k8s/${ns}/${po}.log
    done
}

quickon_log() {
    cp ~/.qc/log/* ${diagnose_dir}/quickon/
    if command_exists q; then
      run q version | tee $diagnose_dir/quickon/q.log
      run q status | tee $diagnose_dir/quickon/q.log
      run q status node | tee $diagnose_dir/quickon/q.log
      run q debug hostinfo > $diagnose_dir/quickon/q.debug.hostinfo.log
    fi
}

kube_status() {
  if command_exists kubectl; then
     cluster_dump
     cluster_events
     component kube-system
     # 兼容2.x版本
     component default
     component cne-system
     # 3.x版本
     component quickon-app
     component quickon-ci
     component quickon-system
  fi
}

pd_collect() {
    os_env
    system_info
    service_status
    network_info
    check_ps_hang
    system_status
    daemon_status
    common_logs
    quickon_log

    # memory
    memory_info
    varlogmessage
    kube_status
    archive
}

pd_collect

echo "----------------------------------------------------------------------"
echo ""
echo "debug file: $tmpdir/diagnose_${timestamp}.tar.gz"
echo ""
echo "----------------------------------------------------------------------"
