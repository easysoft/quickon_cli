// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

import (
	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/qucheng/upgrade"
	"github.com/spf13/cobra"
)

// Option is a struct that defines a command call for "upgrade"
type Option struct {
	Version string
	Test    bool
	log     log.Logger
}

func NewUpgradePlatform(f factory.Factory) *cobra.Command {
	upcmd := &Option{
		log: f.GetLog(),
	}
	up := &cobra.Command{
		Use:     "platform",
		Aliases: []string{"pt"},
		Short:   "upgrade platform to the newest version",
		Args:    cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// 检查升级前的环境
			upcmd.log.Donef("check upgrade environment...")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return upcmd.Run()
		},
	}
	up.Flags().BoolVar(&upcmd.Test, "t", false, "upgrade test mode, only for test")
	return up
}

// Run executes the command logic
func (cmd *Option) Run() error {
	// Run the upgrade command
	cmd.log.Info("check update...")
	cmd.log.Debugf("gen new version manifest")
	err := upgrade.Upgrade(cmd.Version, cmd.Test, cmd.log)
	if err != nil {
		return errors.Errorf("couldn't upgrade: %v", err)
	}
	return nil
}
