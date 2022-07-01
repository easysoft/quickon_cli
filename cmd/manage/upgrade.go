// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package manage

import (
	"fmt"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/qucheng/upgrade"
	"github.com/spf13/cobra"
)

// UpgradeCmd is a struct that defines a command call for "upgrade"
type UpgradeCmd struct {
	Version string
	log     log.Logger
}

func NewUpgradeQucheg(f factory.Factory) *cobra.Command {
	upcmd := &UpgradeCmd{
		log: f.GetLog(),
	}
	up := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades the QuCheng to the newest version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return upcmd.Run()
		},
	}
	return up
}

// Run executes the command logic
func (cmd *UpgradeCmd) Run() error {
	// Run the upgrade command
	cmd.log.Infof("check update...")
	err := upgrade.Upgrade(cmd.Version, cmd.log)
	if err != nil {
		return fmt.Errorf("couldn't upgrade: %v", err)
	}
	return nil
}
