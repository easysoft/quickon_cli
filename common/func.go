// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
)

// GetUG 获取user-agent
func GetUG() string {
	return fmt.Sprintf("%v QAdm/%v", DownloadAgent, Version)
}

// GetK3SURL 获取k3s地址
func GetK3SURL() string {
	// url := fmt.Sprintf("%s/%s/k3s", K3sBinURL, K3sBinVersion)
	url := "https://pkg.qucheng.com/qucheng/k3s/1.23/k3s"
	return url
}

// GetQURL 获取qcadmin地址
func GetQURL() string {
	// url := fmt.Sprintf("%s/%s/k3s", K3sBinURL, K3sBinVersion)
	url := "https://pkg.qucheng.com/qucheng/cli/edge/qcadmin_%s_%s"
	return fmt.Sprintf(url, runtime.GOOS, runtime.GOARCH)
}

// GetBinURL 获取bin地址
func GetBinURL(binName string) string {
	// url := fmt.Sprintf("%s/%s/k3s", K3sBinURL, K3sBinVersion)
	url := "https://pkg.qucheng.com/qucheng/cli/stable/%s/%s-linux-%s"
	return fmt.Sprintf(url, binName, binName, runtime.GOARCH)
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
	return fmt.Sprintf("https://hub.qucheng.com/chartrepo/%s", p)
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
func GetVersion(devops bool, p, version string) string {
	v := strings.Split(version, "-")
	if len(v) != 2 {
		switch p {
		case string(ZenTaoIPDType):
			return DefaultZentaoDevOPSIPDVersion
		case string(ZenTaoBizType):
			return DefaultZentaoDevOPSBizVersion
		case string(ZenTaoMaxType):
			return DefaultZentaoDevOPSMaxVersion
		default:
			if devops {
				return DefaultZentaoDevOPSOSSVersion
			}
			return DefaultQuickonOSSVersion
		}
	}
	return v[1]
}

// GetZenTaoVersion 获取chartRepo channel地址
// func GetZenTaoVersion(p string, qt QuickonType) string {
// 	v := strings.Split(p, "-")
// 	if len(v) != 2 {
// 		if qt == QuickonOSSType {
// 			return GetVersion(DefaultZentaoDevOPSOSSVersion, qt)
// 		}
// 		return GetVersion(DefaultZentaoDevOPSOSSVersion, qt)
// 	}
// 	return v[1]
// }

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
	if devops == true {
		return DefaultZentaoPassName
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
