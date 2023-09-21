// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/cmd/upgrade"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func newCmdUpgrade(f factory.Factory) *cobra.Command {
	up := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrades the cli or plugin to the newest version",
		Aliases: []string{"ug", "ugc"},
	}
	up.AddCommand(upgrade.NewUpgradeQ(f))
	up.AddCommand(upgrade.NewUpgradeOperator(f))
	up.AddCommand(upgrade.NewUpgradePlatform(f))
	return up
}
