#!/usr/bin/env bash

[ $(id -u) -eq 0 ] || exec sudo $0 $@

qcmd=${1:-/usr/local/bin/qcadmin}

current_dir=$(pwd)
tmpdir=/tmp
timestamp=$(date +%s)
diagnose_dir=/tmp/diagnose_${timestamp}
mkdir -p $diagnose_dir
is_ps_hang=false

run() {
    echo
    echo "-----------------run $@------------------"
    timeout 10s $@
    if [ "$?" != "0" ]; then
        echo "failed to collect info: $@"
    fi
    echo "------------End of ${1}----------------"
}

get_distribution() {
    cat /etc/issue*
    if [ -r /etc/os-release ]; then
      cat /etc/os-release
    fi
}

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

# Service status
service_status() {
    run service firewalld status | tee $diagnose_dir/service_status
    run service ntpd status | tee $diagnose_dir/service_status
    run service chronyd status | tee $diagnose_dir/service_status
}


#system info

system_info() {
    run uname -a | tee -a ${diagnose_dir}/os_info
    run uname -r | tee -a ${diagnose_dir}/os_info
    run get_distribution | tee -a ${diagnose_dir}/os_info
    if command_exists lsb_release; then
        run lsb_release | tee -a ${diagnose_dir}/os_info
    fi
    run ulimit -a | tee -a ${diagnose_dir}/os_info
    run sysctl -a | tee -a ${diagnose_dir}/os_info
    run cat /proc/vmstat | tee -a ${diagnose_dir}/os_info
}

#network
network_info() {
    # mkdir -p ${diagnose_dir}/network_info
    #run ifconfig
    run ip --details ad show | tee -a ${diagnose_dir}/network_info
    run ip --details link show | tee -a ${diagnose_dir}/network_info
    run ip route show | tee -a ${diagnose_dir}/network_info
    run iptables-save | tee -a ${diagnose_dir}/network_info
    run cat /proc/net/nf_conntrack | tee -a ${diagnose_dir}/network_info
    run netstat -nt | tee -a ${diagnose_dir}/network_info
    run netstat -nu | tee -a ${diagnose_dir}/network_info
    run netstat -ln | tee -a ${diagnose_dir}/network_info
}

memory_info() {
    run cat /proc/meminfo | tee -a ${diagnose_dir}/memory_info
    run cat /proc/buddyinfo | tee -a ${diagnose_dir}/memory_info
    run cat /proc/vmallocinfo | tee -a ${diagnose_dir}/memory_info
    run cat /proc/slabinfo | tee -a ${diagnose_dir}/memory_info
    run cat /proc/zoneinfo | tee -a ${diagnose_dir}/memory_info
}


# check ps -ef command is hung
check_ps_hang() {
  echo "check if ps -ef command hang" | tee -a ${diagnose_dir}/ps_command_status
  checkD=$(timeout -s 9 2 ps -ef)
  if [ "$?" != "0" ]; then
      echo "ps -ef command is hung" | tee -a ${diagnose_dir}/ps_command_status
      is_ps_hang=true
      echo "start to check which process lead to ps -ef command hang" | tee -a ${diagnose_dir}/ps_command_status
      for f in `find /proc/*/task -name status`
      do
          checkD=$(cat $f|grep "State.*D")
          if [ "$?" == "0" ]; then
              cmdline=$(echo ${f%%status}"cmdline")
              pid=$(echo ${f%%status}"")
              stack=$(echo ${f%%status}"stack")
              timeout -s 9 2 cat $cmdline
              if [ "$?" != "0" ]; then
                  echo "process $pid is in State D and lead to ps -ef process hang,stack info:" | tee -a ${diagnose_dir}/ps_command_status
                  cat $stack | tee -a ${diagnose_dir}/ps_command_status
              fi
          fi
      done
      echo "finish to check which process lead to ps -ef command hang" | tee -a ${diagnose_dir}/ps_command_status
  else
      echo "ps -ef command works fine" | tee -a ${diagnose_dir}/ps_command_status
  fi
}


#system status
system_status() {
    #mkdir -p ${diagnose_dir}/system_status
    run uptime | tee -a ${diagnose_dir}/system_uptime_status
    run top -b -n 1 | tee -a ${diagnose_dir}/system_top_status
    if [ "$is_ps_hang" == "false" ]; then
        run ps -ef | tee -a ${diagnose_dir}/system_status
    else
        echo "ps -ef command hang, skip [ps -ef] check" | tee -a ${diagnose_dir}/system_status
    fi

    run sar -A | tee -a ${diagnose_dir}/system_status

    run df -h | tee -a ${diagnose_dir}/system_df_status

    run cat /proc/mounts | tee -a ${diagnose_dir}/system_mounts_status

    if [ "$is_ps_hang" == "false" ]; then
        run pstree -al | tee -a ${diagnose_dir}/system_status
    else
        echo "ps -ef command hang, skip [pstree -al] check" | tee -a ${diagnose_dir}/system_status
    fi

    run lsof | tee -a ${diagnose_dir}/system_status

    (
        cd /proc
        find -maxdepth 1 -type d -name '[0-9]*' \
         -exec bash -c "ls {}/fd/ | wc -l | tr '\n' ' '" \; \
         -printf "fds (PID = %P), command: " \
         -exec bash -c "tr '\0' ' ' < {}/cmdline" \; \
         -exec echo \; | sort -rn | head | tee -a ${diagnose_dir}/system_status
    )

    echo "----------------start pid leak detect---------------------" | tee -a ${diagnose_dir}/system_status
    ps -elT | awk '{print $4}' | sort | uniq -c | sort -k 1 -g | tail -5 | tee -a ${diagnose_dir}/system_status
    echo "----------------done pid leak detect---------------------" | tee -a ${diagnose_dir}/system_status
}


daemon_status() {
    if command_exists docker; then
        run systemctl status docker -l | tee -a ${diagnose_dir}/docker_status
        run systemctl cat docker | tee -a ${diagnose_dir}/k3s_status
        docker_check
    fi
    if command_exists kubelet; then
        run systemctl status kubelet -l | tee -a ${diagnose_dir}/kubelet_status
    fi
    if command_exists containerd; then
        run systemctl status containerd -l | tee -a ${diagnose_dir}/containerd_status
    fi
    if command_exists k3s; then
        run systemctl status k3s -l | tee -a ${diagnose_dir}/k3s_status
        run systemctl cat k3s | tee -a ${diagnose_dir}/k3s_status
        k3s_check
    fi
    # run systemctl status container-storaged -l | tee -a ${diagnose_dir}/container-storaged_status
}

k3s_check() {
  log_tail_lines=10000
  pidof systemd && journalctl -n ${log_tail_lines} -u k3s.service &> ${diagnose_dir}/logs/k3s.log
}

docker_check() {
    log_tail_lines=10000
    pidof systemd && journalctl -n ${log_tail_lines} -u docker.service &> ${diagnose_dir}/logs/docker.log || tail -n ${log_tail_lines} /var/log/upstart/docker.log &> ${diagnose_dir}/logs/docker.log
    echo "check dockerd process"
    if [ "$is_ps_hang" == "false" ]; then
        run ps -ef|grep -E 'dockerd|docker daemon'|grep -v grep| tee -a ${diagnose_dir}/docker_status
    else
        echo "ps -ef command hang, skip [ps -ef|grep -E 'dockerd|docker daemon'] check" | tee -a ${diagnose_dir}/docker_status
    fi

    #docker info
    run docker info | tee -a ${diagnose_dir}/docker_status
    run docker version | tee -a ${diagnose_dir}/docker_status
    sudo kill -SIGUSR1 $(cat /var/run/docker.pid)
    cp /var/run/docker/libcontainerd/containerd/events.log ${diagnose_dir}/containerd_events.log
    sleep 10
    cp /var/run/docker/*.log ${diagnose_dir}
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
    run dmesg -T | tail -n ${log_tail_lines}  | tee ${diagnose_dir}/logs/dmesg.log
    tail -c 500M /var/log/messages &> ${diagnose_dir}/logs/messages
}

archive() {
    tar -zcvf $tmpdir/diagnose_${timestamp}.tar.gz ${diagnose_dir}
    echo "please get $tmpdir/diagnose_${timestamp}.tar.gz for diagnostics"
}

varlogmessage(){
    grep cloud-init /var/log/messages > $diagnose_dir/logs/cloud-init.log
}

cluster_dump(){
    run kubectl cluster-info dump | tee $diagnose_dir/cluster_dump.log
}

cluster_events(){
    run kubectl get events | tee $diagnose_dir/cluster.log
}

component() {
    local ns="$1"
    mkdir -p $diagnose_dir/cs/$ns/
    local pods=`kubectl get -n $ns po  | awk '{print $1}'|grep -v NAME`
    for po in ${pods}
    do
        kubectl logs -n ${ns} ${po} &> $diagnose_dir/cs/${ns}/${po}.log
    done
}

quchenglog() {
    mkdir -p ${diagnose_dir}/logs/qucheng
    cp ~/.qc/log/* ${diagnose_dir}/logs/qucheng/
     if command_exists q; then
        run q version | tee $diagnose_dir/logs/qucheng/q.log
        run q status | tee $diagnose_dir/logs/qucheng/q.log
        run q status node | tee $diagnose_dir/logs/qucheng/q.log
     fi
}

kube_status() {
  if command_exists kubectl; then
     cluster_dump
     cluster_events
     component kube-system
     component cne-system
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
    quchenglog

    # memory
    memory_info
    varlogmessage
    kube_status
    archive
}

pd_collect

echo "请上传 $tmpdir/diagnose_${timestamp}.tar.gz"
