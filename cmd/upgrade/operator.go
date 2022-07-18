// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

import (
	"os"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func NewUpgradeOperator(f factory.Factory) *cobra.Command {
	up := option{
		log: f.GetLog(),
	}
	upq := &cobra.Command{
		Use:   "operator",
		Short: "upgrade operator to the newest version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			up.DoOperator()
		},
	}
	return upq
}

func (up option) DoOperator() {
	qcexec.CommandRun(os.Args[0], "manage", "plugins", "sync")
	if err := qcexec.CommandRun(os.Args[0], "manage", "plugins", "enable", "cne-operator"); err != nil {
		up.log.Errorf("upgrade plugin cne-operator err: %v", err)
		return
	}
}
