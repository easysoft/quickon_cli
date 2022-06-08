// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/cmd/plugin"
	"github.com/spf13/cobra"
)

func newCmdPlugin() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "plugins",
		Short:   "plugin",
		Aliases: []string{"plugin"},
	}
	cmd.AddCommand(plugin.ListPluginCmd())
	cmd.AddCommand(plugin.InstallPluginCmd())
	cmd.AddCommand(plugin.UnInstallPluginCmd())
	return cmd
}
