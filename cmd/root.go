// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	globalFlags *flags.GlobalFlags
)

func Execute() {
	// create a new factory
	f := factory.DefaultFactory()
	// build the root command
	rootCmd := BuildRoot(f)
	// before hook
	// execute command
	err := rootCmd.Execute()
	// after hook
	if err != nil {
		if globalFlags.Debug {
			f.GetLog().Fatalf("%v", err)
		} else {
			f.GetLog().Fatal(err)
		}
	}
}

// BuildRoot creates a new root command from the
func BuildRoot(f factory.Factory) *cobra.Command {
	// build the root cmd
	rootCmd := NewRootCmd(f)
	persistentFlags := rootCmd.PersistentFlags()
	globalFlags = flags.SetGlobalFlags(persistentFlags)
	// Add main commands
	rootCmd.AddCommand(newCmdVersion())
	rootCmd.AddCommand(newCmdPreCheck())
	rootCmd.AddCommand(newCmdInit(f))
	rootCmd.AddCommand(newCmdJoin(f))
	rootCmd.AddCommand(newCmdUninstall())
	rootCmd.AddCommand(newCmdStatus())
	rootCmd.AddCommand(newCmdServe())
	rootCmd.AddCommand(newCmdUpgrade())
	rootCmd.AddCommand(newCmdManage())
	// Add plugin commands

	rootCmd.AddCommand(KubectlCommand())
	return rootCmd
}

// NewRootCmd returns a new root command
func NewRootCmd(f factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:           "qcadmin",
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "Easily bootstrap a secure control plane for QuCheng",
		Aliases:       []string{"q"},
		Example:       common.RootTPl,
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			if cobraCmd.Annotations != nil {
				return nil
			}
			log := f.GetLog()
			if globalFlags.Silent {
				log.SetLevel(logrus.FatalLevel)
			} else if globalFlags.Debug {
				log.SetLevel(logrus.DebugLevel)
			}

			// TODO apply extra flags
			return nil
		},
	}
}
