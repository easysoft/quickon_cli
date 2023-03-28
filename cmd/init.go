// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Run this command in order to set up the QuCheng control plane",
	}
	skip    bool
	appName string
)

func init() {
	initCmd.PersistentFlags().BoolVar(&skip, "skip-precheck", false, "skip precheck")
	initCmd.PersistentFlags().StringVar(&appName, "app", "zentao", "app name")
}

func newCmdInit(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	defaultArgs := os.Args
	globalToolPath := defaultArgs[0]
	name := "native"
	if file.CheckFileExists(common.GetKubeConfig()) {
		name = "incluster"
	}

	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
			log.Donef("quickon is already initialized, just run %s get cluster status", color.SGreen("%s status", globalToolPath))
			os.Exit(0)
		}
		if name == "incluster" {
			// TODO Check k8s ready
			if _, err := k8s.NewSimpleClient(); err != nil {
				log.Errorf("k8s is not ready, please check your k8s cluster, just run %s ", color.SGreen("%s exp kubectl get nodes", globalToolPath))
				os.Exit(0)
			}
		}
	}
	initCmd.Run = func(cmd *cobra.Command, args []string) {
		if name == "native" {
			log.Infof("start init native provider")
			if err := qcexec.CommandRun(globalToolPath, "cluster", "init", fmt.Sprintf("--debug=%v", globalFlags.Debug)); err != nil {
				log.Errorf("init k8s cluster failed, reason: %v", err)
				return
			}
		}
		if err := qcexec.CommandRun(globalToolPath, "quickon", "init", fmt.Sprintf("--debug=%v", globalFlags.Debug)); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
		if err := qcexec.CommandRun(globalToolPath, "app", "install", "--name", appName, "--api-useip", fmt.Sprintf("--debug=%v", globalFlags.Debug)); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
	}
	return initCmd
}
