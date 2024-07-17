#!/bin/bash

datadir=${1:-"/opt/quickon"}
platformdir="$datadir/platform"
logdir=${2:-"/root/.qc/log"}

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

if ! command_exists docker; then
    echo "k3s is not installed"
    exit 1
fi

k3sbin=$(command -v k3s)

$k3sbin etcd-snapshot save --name qcli --data-dir $platformdir --snapshot-compress --log $logdir/etcd-snapshot-cli.log
$k3sbin etcd-snapshot ls --data-dir $platformdir | grep -v "auto"
