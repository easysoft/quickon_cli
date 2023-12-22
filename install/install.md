---
sidebarTitle: qcadmin backup app
sidebarOrder: 1
---

## 支持的环境

### Linux 发行版

* **Debian**  *12(推荐), 11, 10*
* **Ubuntu**  *20.04, 18.04*
* **CentOS**  *7*
* **Rocky**  *9, 8*

> 建议内核版本`5.14及之后版本`

### 容器运行时

* 内置`containerd`
* docker(本地已安装)

> 多节点时不推荐混用

### k8s/k3s版本

* 对接已有k8s集群, 推荐1.21+版本
* 默认k3s版本为`1.24`

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

> 需要go环境，推荐使用`1.21`

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
curl -sfL https://pkg.qucheng.com/quickon/get.sh | sh -
# 安装渠成平台
z init --provider quickon
# 安装禅道DevOPS
z init --provider devops
# 一键安装
curl -sfL https://pkg.qucheng.com/quickon/install.sh | INSTALL_DOMAIN=example.com sh -
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
