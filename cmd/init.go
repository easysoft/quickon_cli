// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"os"

	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/cmd/precheck"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/db"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/pkg/providers"

	_ "github.com/easysoft/qcadmin/pkg/providers/devops"
	_ "github.com/easysoft/qcadmin/pkg/providers/quickon"

	nativeCluster "github.com/easysoft/qcadmin/pkg/cluster"
)

var (
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize Platform",
	}
	skip      bool
	cProvider = "devops"
	cp        providers.Provider
)

func init() {
	initCmd.Flags().StringVarP(&cProvider, "provider", "p", cProvider, "install provider, support devops, quickon")
	initCmd.PersistentFlags().BoolVar(&skip, "skip-precheck", false, "skip precheck")
}

func newCmdInit(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	pStr := flags.FlagHackLookup("--provider")
	var fs []types.Flag
	if pStr == "" {
		pStr = cProvider
	}
	if reg, err := providers.GetProvider(pStr); err != nil {
		log.Warn(err)
	} else {
		cp = reg
	}
	fs = append(fs, cp.GetFlags()...)

	var preCheck precheck.PreCheck

	defaultArgs := os.Args
	globalToolPath := defaultArgs[0]
	name := "native"
	nCluster := nativeCluster.NewCluster(f)
	if file.CheckFileExists(common.GetKubeConfig()) {
		name = "incluster"
		initCmd.Long = `Found kubeconfig, run this command in order to set up Quickon Control Plane`
	} else {
		fs = append(fs, nCluster.GetInitFlags()...)
		initCmd.Long = `Run this command in order to set up the Kubernetes & Quickon Control Plane`
	}
	meta := cp.GetMeta()
	initCmd.Flags().AddFlagSet(flags.ConvertFlags(initCmd, fs))
	initCmd.Example = cp.GetUsageExample()
	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
			log.Donef("quickon is already initialized, just run %s get cluster status", color.SGreen("%s status", globalToolPath))
			os.Exit(0)
		}
		if name == "incluster" {
			// TODO Check k8s ready
			if _, err := k8s.NewSimpleClient(); err != nil {
				log.Errorf("kube cluster is not ready, please check your cluster, just run %s ", color.SGreen("%s exp kubectl get nodes", globalToolPath))
				os.Exit(0)
			}
		} else {
			preCheck.OffLine = nCluster.OffLine
			preCheck.IgnorePreflightErrors = nCluster.IgnorePreflightErrors
			preCheck.Devops = cProvider == "devops"
			if err := preCheck.Run(); err != nil {
				log.Errorf("precheck failed, reason: %v", err)
				os.Exit(-1)
			}
			if len(nCluster.MasterIPs) == 0 {
				nCluster.MasterIPs = append(nCluster.MasterIPs, exnet.LocalIPs()[0])
			}
		}
		if len(meta.ExtDBHost) > 0 && len(meta.ExtDBPassword) > 0 {
			cfg := &db.Config{
				Host:     meta.ExtDBHost,
				Port:     meta.ExtDBPort,
				Username: meta.ExtDBUser,
				Password: meta.ExtDBPassword,
			}
			if !db.CheckMySQL(cfg) {
				log.Errorf("external db check failed, please check your external db")
				os.Exit(-1)
			}
		}
	}
	initCmd.Run = func(cmd *cobra.Command, args []string) {
		if name == "native" {
			log.Infof("start init native provider")
			if err := nCluster.InitNode(); err != nil {
				log.Errorf("init kube cluster failed, reason: %v", err)
				return
			}
		}
		if err := cp.GetKubeClient(); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
		if err := cp.Check(); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
		if len(meta.IP) == 0 {
			// TODO fix ip, not allow PublicIP
			cfg, _ := config.LoadConfig()
			ip := cfg.Cluster.InitNode
			if len(ip) == 0 {
				ip = exnet.LocalIPs()[0]
			}
			meta.IP = ip
		}
		if err := cp.Install(); err != nil {
			log.Errorf("init quickon failed, reason: %v", err)
			return
		}
		cp.Show()
	}
	return initCmd
}
