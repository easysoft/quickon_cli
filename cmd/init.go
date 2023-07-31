// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/pkg/quickon"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/cmd/precheck"
	nativeCluster "github.com/easysoft/qcadmin/pkg/cluster"
	"github.com/ergoapi/util/exnet"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	_ "github.com/easysoft/qcadmin/pkg/providers/devops"
	_ "github.com/easysoft/qcadmin/pkg/providers/quickon"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a Kubernetes & Quickon cluster",
	}
	skip    bool
	appName string
)

func init() {
	initCmd.PersistentFlags().BoolVar(&skip, "skip-precheck", false, "skip precheck")
	initCmd.PersistentFlags().StringVar(&appName, "app", "zentao", "app name")
}

func newCmdInit(f factory.Factory) *cobra.Command {
	var preCheck precheck.PreCheck
	log := f.GetLog()
	defaultArgs := os.Args
	globalToolPath := defaultArgs[0]
	name := "native"
	nCluster := nativeCluster.NewCluster(f)
	quickonClient := quickon.New(f)
	fs := quickonClient.GetFlags()
	if file.CheckFileExists(common.GetKubeConfig()) {
		name = "incluster"
		initCmd.Long = `Found k8s config, run this command in order to set up Quickon Control Plane`
	} else {
		fs = append(fs, nCluster.GetInitFlags()...)
		initCmd.Long = `Run this command in order to set up the Kubernetes & Quickon Control Plane`
	}
	initCmd.Flags().AddFlagSet(flags.ConvertFlags(initCmd, fs))
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
		} else {
			preCheck.OffLine = nCluster.OffLine
			preCheck.IgnorePreflightErrors = nCluster.IgnorePreflightErrors
			if err := preCheck.Run(); err != nil {
				log.Errorf("precheck failed, reason: %v", err)
				os.Exit(-1)
			}
			if len(nCluster.MasterIPs) == 0 {
				nCluster.MasterIPs = append(nCluster.MasterIPs, exnet.LocalIPs()[0])
			}
		}
	}
	initCmd.Run = func(cmd *cobra.Command, args []string) {
		if name == "native" {
			log.Infof("start init native provider")
			if err := nCluster.InitNode(); err != nil {
				log.Errorf("init k8s cluster failed, reason: %v", err)
				return
			}
		}
		if err := quickonClient.GetKubeClient(); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
		if err := quickonClient.Check(); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
		if !quickonClient.QuickonOSS {
			quickonClient.QuickonType = common.QuickonEEType
		}
		if len(quickonClient.IP) == 0 {
			// TODO fix ip, not allow PublicIP
			cfg, _ := config.LoadConfig()
			ip := cfg.Cluster.InitNode
			if len(ip) == 0 {
				ip = exnet.LocalIPs()[0]
			}
			quickonClient.IP = ip
		}
		if err := quickonClient.Init(); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
		if err := qcexec.CommandRun(globalToolPath, "quickon", "app", "install", "--name", appName, "--api-useip", fmt.Sprintf("--debug=%v", globalFlags.Debug)); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
	}
	return initCmd
}
