// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package manage

import (
	"os"
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	pluginapi "github.com/easysoft/qcadmin/internal/pkg/plugin"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/output"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

var show string

func NewCmdPlugin(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "plugins",
		Short:   "manage plugin",
		Aliases: []string{"plugin"},
	}
	cmd.AddCommand(listPluginCmd(f))
	cmd.AddCommand(installPluginCmd(f))
	cmd.AddCommand(unInstallPluginCmd(f))
	cmd.AddCommand(upgradePluginCmd(f))
	return cmd
}

// func syncPluginFileCmd(f factory.Factory) *cobra.Command {
// 	sync := &cobra.Command{
// 		Use:   "sync",
// 		Short: "sync plugin file",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			dataDir := common.GetDefaultDataDir()
// 			templateVars := map[string]string{
// 				"%{NAMESPACE}%": common.DefaultSystem,
// 			}
// 			if err := deploy.StageFunc(dataDir, templateVars); err != nil {
// 				return err
// 			}
// 			return nil
// 		},
// 	}
// 	return sync
// }

func listPluginCmd(f factory.Factory) *cobra.Command {
	listcmd := &cobra.Command{
		Use:     "list",
		Short:   "list plugin",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := pluginapi.GetAll()
			if err != nil {
				return err
			}
			switch strings.ToLower(show) {
			case "json":
				return output.EncodeJSON(os.Stdout, ps)
			case "yaml":
				return output.EncodeYAML(os.Stdout, ps)
			default:
				table := uitable.New()
				table.MaxColWidth = 80
				table.Wrap = true
				for _, d := range ps {
					table.AddRow("type: ", d.Type)
					table.AddRow("default: ", d.Default)
					if len(d.Item) > 1 {
						str := []string{}
						for _, v := range d.Item {
							if v.Name != d.Default {
								str = append(str, v.Name)
							}
						}
						table.AddRow("optional: ", strings.Join(str, ", "))
					}
					table.AddRow("------------")
				}
				return output.EncodeTable(os.Stdout, table)
			}
		},
	}
	listcmd.Flags().StringVarP(&show, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	return listcmd
}

func installPluginCmd(f factory.Factory) *cobra.Command {
	var version, kubecfg string
	cmd := &cobra.Command{
		Use:     "enable",
		Short:   "install plugin",
		Aliases: []string{"i", "install"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := pluginapi.GetMeta(args...)
			if err != nil {
				return err
			}
			c, err := k8s.NewClient("", kubecfg)
			if err != nil {
				return err
			}
			ps.Client = c
			ps.InstallVersion = version
			return ps.Install()
		},
	}
	cmd.Flags().StringVarP(&version, "version", "v", "", "plugin")
	cmd.Flags().StringVarP(&kubecfg, "kubeconfig", "k", common.GetDefaultNewKubeConfig(), "kubeconfig file")
	return cmd
}

func unInstallPluginCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "disable",
		Short:   "uninstall plugin",
		Aliases: []string{"uninstall", "x", "un", "delete", "del"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := pluginapi.GetMeta(args...)
			if err != nil {
				return err
			}
			c, err := k8s.NewClient("", "")
			if err != nil {
				return err
			}
			ps.Client = c
			return ps.UnInstall()
		},
	}
	return cmd
}

func upgradePluginCmd(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade plugin",
		Aliases: []string{"up", "ug"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := pluginapi.GetMeta(args...)
			if err != nil {
				return err
			}
			c, err := k8s.NewClient("", "")
			if err != nil {
				return err
			}
			ps.Client = c
			return ps.Upgrade()
		},
	}
	return cmd
}
