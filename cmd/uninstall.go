// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
)

var (
	cleanCluster bool
)

func newCmdUninstall(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	defaultArgs := os.Args
	globalToolPath := defaultArgs[0]
	uninstallCmd := &cobra.Command{
		Use:     "uninstall",
		Short:   "Uninstall cluster",
		Aliases: []string{"un", "clean"},
	}
	uninstallCmd.Run = func(cmd *cobra.Command, args []string) {
		log.Debugf("start uninstall platform")
		if err := qcexec.CommandRun(globalToolPath, "platform", "uninstall", fmt.Sprintf("--debug=%v", globalFlags.Debug)); err != nil {
			log.Errorf("uninstall platform failed, reason: %v", err)
			return
		}
		if cleanCluster {
			// TODO 检查集群是否是quickon安装的
			log.Debugf("start uninstall cluster")
			if err := qcexec.CommandRun(globalToolPath, "cluster", "clean", fmt.Sprintf("--debug=%v", globalFlags.Debug)); err != nil {
				log.Errorf("uninstall cluster failed, reason: %v", err)
				return
			}
		}
		log.Donef("uninstall success")
	}
	uninstallCmd.PersistentFlags().BoolVar(&cleanCluster, "all", true, "clean all")
	return uninstallCmd
}
