# qcadmin(q)

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/easysoft/qucheng_cli?filename=go.mod&style=flat-square)
![GitHub Workflow Status (event)](https://img.shields.io/github/workflow/status/easysoft/qucheng_cli/Release?style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/easysoft/qucheng_cli?style=flat-square)
![GitHub](https://img.shields.io/badge/license-ZPL%20%2B%20AGPL-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/easysoft/qucheng_cli)](https://goreportcard.com/report/github.com/easysoft/qucheng_cli)
[![Releases](https://img.shields.io/github/release-pre/easysoft/qucheng_cli.svg)](https://github.com/easysoft/qucheng_cli/releases)
[![docs](https://img.shields.io/badge/docs-done-green)](https://www.qucheng.com/)


> qcadmin is an open-source lightweight cli tool for managing qucheng.

## Requirements

<table>
  <tbody>
    <tr>
    	<th width='320'>OS</th>
    	<th>Minimum Requirements</th>
    </tr>
    <tr>
      <td><b>Debian(Recommended)</b> <i>Bullseye</i>, <i>Buster</i></td>
      <td>2 CPU cores, 4 GB memory, 40 GB disk space</td>
    </tr>
    <tr>
      <td><b>Ubuntu</b> <i>16.04</i>, <i>18.04</i></td>
      <td>2 CPU cores, 4 GB memory, 40 GB disk space</td>
    </tr>
		<tr>
    <td><b>CentOS</b> <i>7.x</i></td>
      <td>2 CPU cores, 4 GB memory, 40 GB disk space</td>
    </tr>
  </tbody>
</table>

> Recommended Linux Kernel Version: 5.14 or later

## Container runtimes

> If you use q to set up a cluster, use containerd by default. Alternatively, you can manually install Docker  runtimes before you create a cluster.

## Usage

### Install

#### 1. Building From Source

`qcadmin(q)` is currently using go v1.16 or above. In order to build ergo from source you must:

```bash
# Clone the repo
# Build and run the executable
make generate
make build
```

#### 2. Linux Binary

Downloaded from pre-compiled binaries

```bash
# 稳定版本 / stable / tag (Recommended)
curl https://pkg.qucheng.com/qucheng/cli/stable/get.sh | sh -
q init
# 开发版 / edge / master
curl https://pkg.qucheng.com/qucheng/cli/edge/get.sh | sh -
q init -q edge
```

#### 3. Debian/CentOS 7

```bash
# debian
echo "deb [trusted=yes] https://apt.fury.io/qucheng/ /" | tee /etc/apt/sources.list.d/qcadmin.list
apt update
apt search qcadmin
apt install qcadmin
# centos7
cat > /etc/yum.repos.d/qcadmin.repo << EOF
[fury]
name=Qucheng Yum Repo
baseurl=https://yum.fury.io/qucheng/
enabled=1
gpgcheck=0
EOF
yum makecache
yum install qcadmin
```

### Quick Start

```bash
# create qucheng cluster
q init
# create a k3s cluster with other cidr
q init --podsubnet 10.42.0.0/16 --svcsubnet 10.43.0.0/16
```
