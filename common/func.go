// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import (
	"fmt"
	"strings"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
)

// GetUG 获取user-agent
func GetUG() string {
	return fmt.Sprintf("%v QAdm/%v", DownloadAgent, Version)
}

func GetDefaultBinDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultBinDir
}

func GetDefaultCacheDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCacheDir
}

func GetDefaultDataDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultDataDir
}

func GetDefaultLogDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultLogDir
}

// GetChartRepo 获取chartrepo地址
func GetChartRepo(p string) string {
	if strings.HasPrefix(p, "test") || strings.HasPrefix(p, "edge") {
		p = "test"
	} else {
		p = "stable"
	}
	return fmt.Sprintf("https://hub.zentao.net/chartrepo/%s", p)
}

// GetChannel 获取chartrepo channel地址
func GetChannel(p string) string {
	if strings.HasPrefix(p, "test") || strings.HasPrefix(p, "edge") {
		p = "test"
	} else {
		p = "stable"
	}
	return p
}

// GetVersion 获取版本地址
func GetVersion(devops bool, typ, version string) string {
	if !devops {
		// 渠成不支持版本
		return DefaultQuickonOSSVersion
	}
	v := strings.Split(version, "-")
	if len(v) == 2 {
		return v[1]
	}
	defaultVersion := DefaultZentaoDevOPSOSSVersion
	switch typ {
	case string(ZenTaoIPDType):
		defaultVersion = DefaultZentaoDevOPSIPDVersion
	case string(ZenTaoBizType):
		defaultVersion = DefaultZentaoDevOPSBizVersion
	case string(ZenTaoMaxType):
		defaultVersion = DefaultZentaoDevOPSMaxVersion
	}
	if strings.HasSuffix(defaultVersion, ".0") {
		defaultVersion = strings.TrimSuffix(defaultVersion, ".0")
	}
	return defaultVersion
}

func GetDefaultConfig() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCfgDir + "/cluster.yaml"
}

func DefaultKubeConfig() string {
	d := fmt.Sprintf("%v/.kube", zos.GetHomeDir())
	// os.MkdirAll(d, FileMode0644)
	return fmt.Sprintf("%v/config", d)
}

func DefaultQuickONKubeConfig() string {
	home := zos.GetHomeDir()
	d := home + "/" + DefaultCfgDir + "/.kube"
	// os.MkdirAll(d, FileMode0644)
	return fmt.Sprintf("%v/config", d)
}

// GetKubeConfig get kubeconfig
func GetKubeConfig() string {
	kubeCfg := DefaultQuickONKubeConfig()
	if file.CheckFileExists(kubeCfg) {
		return kubeCfg
	}
	return DefaultKubeConfig()
}

func GetCustomConfig(name string) string {
	home := zos.GetHomeDir()
	return fmt.Sprintf("%s/%s/%s", home, DefaultCfgDir, name)
}

func GetAPI(path string) string {
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("https://api.qucheng.com/%s", path)
}

func GetCustomQuickonDir(path string) string {
	if zos.IsMacOS() {
		return fmt.Sprintf("%v/%v", zos.GetHomeDir(), path)
	}
	return path
}

func GetDefaultQuickonBackupDir(path string) string {
	if len(path) == 0 {
		path = DefaultQuickonDataDir
	} else {
		path = strings.TrimSuffix(path, "/")
	}
	return fmt.Sprintf("%s/backup", path)
}

func GetDefaultQuickonPlatformDir(path string) string {
	if len(path) == 0 {
		path = DefaultQuickonDataDir
	} else {
		path = strings.TrimSuffix(path, "/")
	}
	return fmt.Sprintf("%s/platform", path)
}

// GetDefaultSystemNamespace get quickon default system ns
func GetDefaultSystemNamespace(newVersion bool) string {
	if newVersion {
		return DefaultSystemNamespace
	}
	return DefaultSystem
}

func GetDefaultQuickONNamespace() []string {
	var ns []string
	ns = append(ns, DefaultAppNamespace, DefaultCINamespace, DefaultSystemNamespace, DefaultStorageNamespace)
	return ns
}

func GetCustomScripts(path string) string {
	return fmt.Sprintf("%s/%s", GetDefaultDataDir(), path)
}

// GetReleaseName get chart release name
func GetReleaseName(devops bool) string {
	if devops {
		return DefaultZentaoPaasName
	}
	return DefaultQuchengName
}

// GetDefaultQuickonDir get quickon default nfs dir
func GetDefaultNFSStoragePath(dir string) string {
	if len(dir) == 0 {
		dir = DefaultQuickonDataDir
	}
	return fmt.Sprintf("%s/storage/nfs", dir)
}

// GetInstallType get install type
func GetInstallType(devops bool) string {
	if devops {
		return "devops"
	}
	return "quickon"
}

// GetCustomLogFile 获取logfile完整路径
func GetCustomLogFile(name string) string {
	if name == "" {
		name = Version
	}
	return GetDefaultLogDir() + "/" + name
}
