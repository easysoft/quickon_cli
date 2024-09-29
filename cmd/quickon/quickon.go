// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package quickon

import (
	"fmt"

	"github.com/ergoapi/util/confirm"
	"github.com/ergoapi/util/exnet"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/internal/api/statistics"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/pkg/providers"
	"github.com/easysoft/qcadmin/pkg/quickon"
)

var (
	initCmd = &cobra.Command{
		Use: "init",
	}
	cProvider = "devops"
	cp        providers.Provider
)

func init() {
	initCmd.Flags().StringVarP(&cProvider, "provider", "p", cProvider, "install provider, support devops, quickon")
}

func InitCommand(f factory.Factory) *cobra.Command {
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
	meta := cp.GetMeta()
	// quickonClient := quickon.New(f)
	// fs = append(fs, quickonClient.GetFlags()...)
	initCmd.Short = fmt.Sprintf("init %s platform", pStr)
	initCmd.Example = cp.GetUsageExample()
	initCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if err := cp.GetKubeClient(); err != nil {
			return err
		}
		return cp.Check()
	}

	initCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := cp.GetKubeClient(); err != nil {
			return err
		}
		if len(meta.IP) == 0 {
			cfg, _ := config.LoadConfig()
			ip := cfg.Cluster.InitNode
			if len(ip) == 0 {
				ip = exnet.LocalIPs()[0]
			}
			meta.IP = ip
		}
		if err := cp.Install(); err != nil {
			return err
		}
		cp.Show()
		return nil
	}
	initCmd.Flags().AddFlagSet(flags.ConvertFlags(initCmd, fs))
	return initCmd
}

func UninstallCommand(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	quickonClient := quickon.New(f)
	uninstall := &cobra.Command{
		Use:     "uninstall",
		Short:   "uninstall platform",
		Aliases: []string{"clean"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := quickonClient.GetKubeClient(); err != nil {
				return err
			}
			status, _ := confirm.Confirm("Are you sure to uninstall platform")
			if status {
				if err := quickonClient.UnInstall(); err != nil {
					return err
				}
				log.Done("uninstall platform success")
				statistics.SendStatistics("uninstall-platform")
				return nil
			}
			log.Donef("cancel uninstall platform")
			return nil
		},
	}
	return uninstall
}

func CheckCommand(f factory.Factory) *cobra.Command {
	check := &cobra.Command{
		Use:   "check",
		Short: "check cluster ready",
		RunE: func(cmd *cobra.Command, args []string) error {
			quickonClient := quickon.New(f)
			if err := quickonClient.GetKubeClient(); err != nil {
				return err
			}
			return quickonClient.Check()
		},
	}
	return check
}
