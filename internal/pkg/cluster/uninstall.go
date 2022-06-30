// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"fmt"
	"os"

	"github.com/easysoft/qcadmin/common"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/ergoapi/util/file"
)

func (p *Cluster) Uninstall() error {
	if !file.CheckFileExists(common.GetDefaultConfig()) {
		p.Log.Done("uninstall cluster success")
		return nil
	}
	var uninstallFile string
	checkfile := common.GetCustomConfig(common.InitModeCluster)
	mode := "native"
	if file.CheckFileExists(checkfile) {
		uninstallFile = "qucheng-uninstall.sh"
		mode = "incluster"
	} else {
		uninstallFile = "k3s-uninstall.sh"
	}

	uninstallShell := fmt.Sprintf("%s/hack/manifests/scripts/%s", common.GetDefaultDataDir(), uninstallFile)
	p.Log.Debugf("gen %s uninstall script: %v", mode, uninstallShell)
	// 移除qcadmin初始化文件
	if err := qcexec.RunCmd("/bin/bash", uninstallShell, os.Args[0]); err != nil {
		return err
	}

	os.Remove(checkfile)

	// 移除qcadmin配置文件
	if file.CheckFileExists(common.GetDefaultConfig()) && mode == "native" {
		os.Remove(common.GetDefaultConfig())
	} else if mode == "incluster" {
		os.Remove(common.GetCustomConfig(common.InitModeCluster))
	}
	initfile := common.GetCustomConfig(common.InitFileName)
	if file.CheckFileExists(initfile) {
		os.Remove(initfile)
	}
	return nil
}
