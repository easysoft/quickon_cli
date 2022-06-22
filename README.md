# qcadmin(q)

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/easysoft/qucheng_cli?filename=go.mod&style=flat-square)
![GitHub Workflow Status (event)](https://img.shields.io/github/workflow/status/easysoft/qucheng_cli/Release?style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/easysoft/qucheng_cli?style=flat-square)
![GitHub](https://img.shields.io/badge/license-ZPL%20%2B%20AGPL-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/easysoft/qucheng_cli)](https://goreportcard.com/report/github.com/easysoft/qucheng_cli)
[![Releases](https://img.shields.io/github/release-pre/easysoft/qucheng_cli.svg)](https://github.com/easysoft/qucheng_cli/releases)
[![docs](https://img.shields.io/badge/docs-done-green)](https://www.qucheng.com/)


> 中文 | [English](README-EN.md)

使用 `qcadmin`(`q`)，您可以轻松、高效、灵活地单独或整体安装渠成平台。

## 支持的环境

### Linux 发行版

* **Debian**  *Buster(推荐), Stretch*
* **Ubuntu**  *16.04, 18.04*
* **CentOS**  *7*

> 建议内核版本`5.14及之后版本`

### 容器运行时

默认使用k3s内置的`Containerd`, 如果本地已经安装docker，则优先使用docker

### k8s/k3s版本

* 对接已有k8s集群, 推荐1.20+版本
* 默认k3s版本为`1.23`

## 要求和建议

* 最低资源要求：
  * 2 核虚拟 CPU
  * 4 GB 内存
  * 40 GB 储存空间

* 操作系统要求：

  * 节点时间同步。
  * `sudo`/`curl` 节点需已安装。
  * 网络正常。

> * 建议您的操作系统环境足够干净 (不安装任何其他软件)，否则可能会发生冲突。

## 使用

### 安装二进制

#### 1. 从源码安装

> 需要go环境，且版本大于1.16

```bash
# Clone the repo
# Build and run the executable
make generate
make build
```

#### 2. 二进制安装

> 使用我们提供的编译二进制文件。可以从下面或者github获取

```bash
# 稳定版本 / stable / tag (Recommended)
curl https://pkg.qucheng.com/qucheng/cli/stable/get.sh | sh -
q init
# 开发版 / edge / master
curl https://pkg.qucheng.com/qucheng/cli/edge/get.sh | sh -
q init -q edge
```

#### 3. 包安装

> 目前仅提供deb或者rpm包方式安装。

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

### 快速开始

> 快速入门使用 `all-in-one` 安装，这是熟悉 渠成平台 的良好开始。

```bash
# create qucheng cluster
q init
# create a k3s cluster with other cidr
q init --podsubnet 10.42.0.0/16 --svcsubnet 10.43.0.0/16
```

## 相关文档

[版本升级](https://github.com/easysoft/qucheng_cli/wiki/%E7%89%88%E6%9C%AC%E5%8D%87%E7%BA%A7)
