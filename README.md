# qcadmin(q)

![GitHub Workflow Status (event)](https://img.shields.io/github/workflow/status/easysoft/qucheng_cli/tag?style=flat-square)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/easysoft/qucheng_cli?filename=go.mod&style=flat-square)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/easysoft/qucheng_cli?style=flat-square)
![GitHub all releases](https://img.shields.io/github/downloads/easysoft/qucheng_cli/total?style=flat-square)
![GitHub](https://img.shields.io/github/license/easysoft/qucheng_cli?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/easysoft/qucheng_cli)](https://goreportcard.com/report/github.com/easysoft/qucheng_cli)
[![Releases](https://img.shields.io/github/release-pre/easysoft/qucheng_cli.svg)](https://github.com/easysoft/qucheng_cli/releases)
[![docs](https://img.shields.io/badge/docs-done-green)](https://www.qucheng.cn/)

compatibility:

- [x] 100% support `Debian 11+`

## Quick start

```bash
q init
```

## Quick build

```bash
make generate
make build
```

## Install

### Building From Source

`qcadmin(q)` is currently using go v1.16 or above. In order to build ergo from source you must:

```bash
# Clone the repo
# Build and run the executable
make generate
make build
```

### Linux Binary

Downloaded from pre-compiled binaries

```bash
# 稳定版本 / stable / tag (Recommended)
curl https://pkg.qucheng.com/qucheng/cli/stable/get.sh | sh -
q init
# 开发版 / edge / master
curl https://pkg.qucheng.com/qucheng/cli/edge/get.sh | sh -
q init -q edge
```

## Upgrade

```bash
# upgrade self
q upgrade q
```
