# qcadmin(q)

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/easysoft/quickon_cli?filename=go.mod&style=flat-square)
[![Release](https://github.com/easysoft/quickon_cli/actions/workflows/release.yml/badge.svg)](https://github.com/easysoft/quickon_cli/actions/workflows/release.yml)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/easysoft/quickon_cli?style=flat-square)
![GitHub](https://img.shields.io/badge/license-ZPL%20%2B%20AGPL-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/easysoft/quickon_cli)](https://goreportcard.com/report/github.com/easysoft/quickon_cli)
[![Releases](https://img.shields.io/github/release-pre/easysoft/quickon_cli.svg)](https://github.com/easysoft/quickon_cli/releases)
[![TODOs](https://img.shields.io/endpoint?url=https://api.tickgit.com/badge?repo=github.com/easysoft/quickon_cli)](https://www.tickgit.com/browse?repo=github.com/easysoft/quickon_cli)
[![docs](https://img.shields.io/badge/docs-done-green)](https://www.qucheng.com/)
[![Chat on QQ](https://img.shields.io/badge/chat-768721743-blueviolet?logo=TencentQQ)](https://img.qucheng.com/group/qq.jpg)

> English | [中文](README.md)

qcadmin is an open-source lightweight cli tool for managing quickon.

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
      <td><b>Ubuntu</b> <i>20.04</i>, <i>18.04</i></td>
      <td>2 CPU cores, 4 GB memory, 40 GB disk space</td>
    </tr>
		<tr>
    <td><b>CentOS</b> <i>7.x</i></td>
      <td>2 CPU cores, 4 GB memory, 40 GB disk space</td>
    </tr>
		<tr>
    <td><b>Rocky</b> <i>9.x</i>, <i>8.x</i></td>
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
echo "deb [trusted=yes] https://repo.qucheng.com/quickon/apt/ /" | tee /etc/apt/sources.list.d/quickon.list
apt update
apt search qcadmin
apt install qcadmin
# centos7
cat > /etc/yum.repos.d/quickon.repo << EOF
[quickon]
name=Quickon Repo
baseurl=https://repo.qucheng.com/quickon/yum/
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
```
