// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import (
	"fmt"
	"os"
	"runtime"
	"strings"

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

// GetChannel 获取chartrepo channel地址
func GetVersion(p string) string {
	if strings.Contains(p, "-") {
		v := strings.Split(p, "-")
		if len(v) != 2 {
			return DefaultQuchengVersion
		}
		return v[1]
	}
	return ""
}

func GetDefaultConfig() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCfgDir + "/cluster.yaml"
}

func GetDefaultKubeConfig() string {
	d := fmt.Sprintf("%v/.kube", zos.GetHomeDir())
	os.MkdirAll(d, FileMode0644)
	return fmt.Sprintf("%v/config", d)
}

func GetCustomConfig(name string) string {
	home := zos.GetHomeDir()
	return fmt.Sprintf("%s/%s/%s", home, DefaultCfgDir, name)
}

func GetAPI(path string) string {
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("https://api.qucheng.com/%s", path)
}

func GetDefaultQuickonDir() string {
	if zos.IsMacOS() {
		return fmt.Sprintf("%v/%v", zos.GetHomeDir(), DefaultQuickonDir)
	}
	return DefaultQuickonDir
}
