// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/cmd/manage"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func newCmdManage(f factory.Factory) *cobra.Command {
	m := &cobra.Command{
		Use:     "manage",
		Short:   "Manage qucheng tools",
		Aliases: []string{"m", "op"},
	}
	m.AddCommand(manage.NewCmdPlugin(f))
	m.AddCommand(manage.NewResetPassword(f))
	m.AddCommand(manage.NewUpgradeQucheg(f))
	return m
}

func newCmdManageGet(f factory.Factory) *cobra.Command {
	m := &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources.",
	}
	m.AddCommand(manage.NewCmdGetNode(f))
	m.AddCommand(manage.NewCmdGetApp(f))
	return m
}
