#!/bin/bash

apt update

apt install -y --no-install-recommends bridge-utils qemu-system libvirt-clients libvirt-daemon-system

# apt install -y virt-manager
