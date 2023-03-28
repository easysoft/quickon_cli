// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package quickon

import (
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/pkg/quickon"
	"github.com/ergoapi/util/expass"
	"github.com/spf13/cobra"
)

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

func InitCommand(f factory.Factory) *cobra.Command {
	quickonClient := quickon.New(f)
	init := &cobra.Command{
		Use:   "init",
		Short: "init quickon",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := quickonClient.GetKubeClient(); err != nil {
				return err
			}
			return quickonClient.Check()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := quickonClient.GetKubeClient(); err != nil {
				return err
			}
			if len(quickonClient.IP) == 0 {
				cfg, _ := config.LoadConfig()
				quickonClient.IP = cfg.InitNode
			}
			return quickonClient.Init()
		},
	}
	init.Flags().StringVar(&quickonClient.Domain, "domain", "", "global domain")
	init.Flags().StringVar(&quickonClient.IP, "ip", "", "ip")
	init.Flags().StringVar(&quickonClient.ConsolePassword, "quickon-password", expass.PwGenAlphaNum(32), "quickon console password")
	init.Flags().StringVar(&quickonClient.Version, "version", common.DefaultQuchengVersion, "quickon version")
	init.Flags().BoolVar(&quickonClient.OssMode, "oss", true, "quickon type, default oss, also support enterprise")
	return init
}

func UninstallCommand(f factory.Factory) *cobra.Command {
	quickonClient := quickon.New(f)
	uninstall := &cobra.Command{
		Use:     "uninstall",
		Short:   "uninstall quickon",
		Aliases: []string{"clean"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := quickonClient.GetKubeClient(); err != nil {
				return err
			}
			return quickonClient.UnInstall()
		},
	}
	return uninstall
}
