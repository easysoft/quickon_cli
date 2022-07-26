// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

import (
	"fmt"
	"os"

	"github.com/easysoft/qcadmin/common"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/qucheng/upgrade"
	"github.com/spf13/cobra"
)

// Option is a struct that defines a command call for "upgrade"
type Option struct {
	Version string
	log     log.Logger
}

func NewUpgradeQucheng(f factory.Factory) *cobra.Command {
	upcmd := &Option{
		log: f.GetLog(),
	}
	up := &cobra.Command{
		Use:     "quickon",
		Aliases: []string{"qc", "qucheng"},
		Short:   "Upgrades the QuCheng to the newest version",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return upcmd.Run()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// 升级成功
			return upcmd.CleanOrInstall()
		},
	}
	return up
}

// Run executes the command logic
func (cmd *Option) Run() error {
	// Run the upgrade command
	cmd.log.Info("check update...")
	cmd.log.Debugf("gen new version manifest")
	err := upgrade.Upgrade(cmd.Version, cmd.log)
	if err != nil {
		return fmt.Errorf("couldn't upgrade: %v", err)
	}
	return nil
}

// Clean executes the command logic
func (cmd *Option) CleanOrInstall() error {
	cmd.log.Info("cleanup deprecated resources")
	if err := qcexec.CommandRun(os.Args[0], "exp", "helm", "uninstall", "--name", "cne-api", "--namespace", common.DefaultSystem); err != nil {
		cmd.log.Errorf("clean cne-api err: %v", err)
	}
	cmd.log.Info("patch new resources")
	cmd.log.Info("start deploy operator plugins: cne-operator")
	if err := qcexec.CommandRun(os.Args[0], "manage", "plugins", "enable", "cne-operator"); err != nil {
		cmd.log.Errorf("deploy plugin cne-operator err: %v", err)
	}
	return nil
}
