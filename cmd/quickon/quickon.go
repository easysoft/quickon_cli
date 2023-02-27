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
			quickonCliet := quickon.New(f)
			if err := quickonCliet.GetKubeClient(); err != nil {
				return err
			}
			return quickonCliet.Check()
		},
	}
	return check
}

func InitCommand(f factory.Factory) *cobra.Command {
	quickonCliet := quickon.New(f)
	init := &cobra.Command{
		Use:   "init",
		Short: "init quickon",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := quickonCliet.GetKubeClient(); err != nil {
				return err
			}
			if len(quickonCliet.IP) == 0 {
				cfg, _ := config.LoadConfig()
				quickonCliet.IP = cfg.InitNode
			}
			return quickonCliet.Init()
		},
	}
	init.Flags().StringVar(&quickonCliet.Domain, "domain", "", "global domain")
	init.Flags().StringVar(&quickonCliet.IP, "ip", "", "ip")
	init.Flags().StringVar(&quickonCliet.ConsolePassword, "quickon-password", expass.PwGenAlphaNum(32), "quickon console password")
	init.Flags().StringVar(&quickonCliet.Version, "version", common.DefaultQuchengVersion, "quick version")
	return init
}

func UninstallCommand(f factory.Factory) *cobra.Command {
	quickonCliet := quickon.New(f)
	uninstall := &cobra.Command{
		Use:     "uninstall",
		Short:   "uninstall quickon",
		Aliases: []string{"clean"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := quickonCliet.GetKubeClient(); err != nil {
				return err
			}
			return quickonCliet.UnInstall()
		},
	}
	return uninstall
}
