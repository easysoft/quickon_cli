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

使用 `qcadmin`(`q`, `z`)，您可以轻松、高效、灵活地单独或整体安装渠成平台或者禅道DevOPS解决方案。

## 支持的环境

### Linux 发行版

* **Debian**  *12(推荐), 11, 10*
* **Ubuntu**  *24.04,22.04*
* **CentOS**  *7(不推荐)*
* **Rocky**  *9(不推荐), 8*

> 建议内核版本`5.14及之后版本`

### 容器运行时

- 内置`containerd`
- docker(本地已安装)

> 多节点时不推荐混用

### k8s/k3s版本

* 对接已有k8s集群, 推荐1.21+版本
* 默认k3s版本为`1.28`

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

> 需要go环境，推荐使用`1.22`

```bash
# Clone the repo
# Build and run the executable
#make generate
#make build
task local
task
```

#### 2. 二进制安装

> 使用我们提供的编译二进制文件。可以从下面或者github获取

```bash
# 稳定版本 / stable / tag (Recommended)
curl -sfL https://pkg.zentao.net/get.sh | sh -
# 安装渠成平台
z init --provider quickon
# 安装禅道DevOPS
z init --provider devops
# 一键安装
curl -sfL https://pkg.zentao.net/devops.sh | INSTALL_DOMAIN=example.com sh -
```

#### 3. 包安装

> 目前仅提供deb或者rpm包方式安装。

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

### 快速开始

> 快速入门使用 `all-in-one` 安装，这是熟悉 渠成平台 的良好开始。

```bash
# create quickon cluster
z init --provider quickon
# create zentao devops cluster(default)
z init --provider zentao
z init
```

### 高级使用

```bash
# step 1. 创建k3s集群
z cluster init
# step 2. 初始化渠成平台或devops
z quickon init --provider quickon
z quickon init
```

## 相关文档

[文档](./docs/index.md)

## 问题反馈

* GitHub Issues
* QQGroup: 768721743

## Contributors

<!-- readme: collaborators,contributors -start -->
<!-- readme: collaborators,contributors -end -->
<a href="https://github.com/easysoft/quickon_cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=easysoft/quickon_cli" />
</a>
