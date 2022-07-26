// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/cmd/app"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func newCmdApp(f factory.Factory) *cobra.Command {
	appCmd := &cobra.Command{
		Use:    "app",
		Short:  " Manage applications",
		Hidden: false,
	}
	appCmd.AddCommand(app.NewCmdAppExec(f))
	appCmd.AddCommand(app.NewCmdAppGet(f))
	appCmd.AddCommand(app.NewCmdAppLogs(f))
	appCmd.AddCommand(app.NewCmdAppList(f))
	return appCmd
}
