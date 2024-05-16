// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/cmd/app"
	"github.com/easysoft/qcadmin/cmd/manage"
	"github.com/easysoft/qcadmin/cmd/manage/db"
	"github.com/easysoft/qcadmin/cmd/quickon"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
)

func newCmdApp(f factory.Factory) *cobra.Command {
	appCmd := &cobra.Command{
		Use:     "app",
		Short:   "Manage applications",
		Version: "20230906",
	}
	cfg, _ := config.LoadConfig()
	if cfg == nil || !cfg.Quickon.DevOps {
		appCmd.AddCommand(app.NewCmdAppExec(f))
		appCmd.AddCommand(app.NewCmdAppGet(f))
		appCmd.AddCommand(app.NewCmdAppLogs(f))
		appCmd.AddCommand(app.NewCmdAppInstall(f))
		appCmd.AddCommand(app.NewCmdAppMarket(f))
	}
	appCmd.AddCommand(app.NewCmdAppList(f))
	return appCmd
}

func newCmdPlatform(f factory.Factory) *cobra.Command {
	platformCmd := &cobra.Command{
		Use:     "platform",
		Short:   "Platform commands",
		Aliases: []string{"qc", "quickon", "pt"},
		Version: "20230330",
	}
	platformCmd.AddCommand(newCmdApp(f))
	platformCmd.AddCommand(quickon.CheckCommand(f))
	platformCmd.AddCommand(quickon.InitCommand(f))
	platformCmd.AddCommand(quickon.UninstallCommand(f))
	platformCmd.AddCommand(manage.NewCmdPlugin(f))
	platformCmd.AddCommand(db.NewCmdDB(f))
	cfg, _ := config.LoadConfig()
	if cfg == nil || !cfg.Quickon.DevOps {
		platformCmd.AddCommand(manage.NewResetPassword(f))
	}
	if cfg == nil || kutil.IsLegalDomain(cfg.Domain) {
		platformCmd.AddCommand(manage.NewRenewTLS(f))
	}
	return platformCmd
}
