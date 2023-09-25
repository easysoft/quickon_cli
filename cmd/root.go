// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/excmd"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	globalFlags *flags.GlobalFlags
)

// func init() {
// 	cobra.OnInitialize(initConfig)
// }

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
		if !strings.Contains(err.Error(), "unknown command") {
			f.GetLog().Info("----------------------------")
			bugmsg := "found bug: submit the error message to Github or Gitee\n\t Github: https://github.com/easysoft/quickon_cli/issues/new?assignees=&labels=&template=bug-report.md\n\t Gitee: https://gitee.com/wwccss/qucheng_cli/issues"
			f.GetLog().Info(bugmsg)
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
	rootCmd.AddCommand(newCmdVersion(f))
	rootCmd.AddCommand(newCmdInit(f))
	rootCmd.AddCommand(newCmdUninstall(f))
	rootCmd.AddCommand(newCmdStatus(f))
	rootCmd.AddCommand(newCmdUpgrade(f))
	rootCmd.AddCommand(newCmdCluster(f))
	rootCmd.AddCommand(newCmdPlatform(f))
	rootCmd.AddCommand(newCmdBackup(f))
	// Add plugin commands
	rootCmd.AddCommand(newCmdExperimental(f))
	rootCmd.AddCommand(newManCmd())
	rootCmd.AddCommand(newCmdBugReport(f))
	rootCmd.AddCommand(newCmdDebug(f))

	// Deprecated commands, will be removed in the future
	deprecatedAppCommand := newCmdApp(f)
	deprecatedAppCommand.Deprecated = fmt.Sprintf("use %s instead", color.SGreen("%s platform app", os.Args[0]))
	rootCmd.AddCommand(deprecatedAppCommand)

	args := os.Args
	if len(args) > 1 {
		pluginHandler := excmd.NewDefaultPluginHandler(common.GetDefaultBinDir(), common.ValidPrefixes)
		cmdPathPieces := args[1:]
		if _, _, err := rootCmd.Find(cmdPathPieces); err != nil {
			var cmdName string // first "non-flag" arguments
			for _, arg := range cmdPathPieces {
				if !strings.HasPrefix(arg, "-") {
					cmdName = arg
					break
				}
			}
			switch cmdName {
			case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
				// Don't search for a plugin
			default:
				if err := excmd.HandlePluginCommand(pluginHandler, cmdPathPieces); err != nil {
					fmt.Fprintf(os.Stdout, "Error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}
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
			qlog := f.GetLog()
			if globalFlags.Silent {
				qlog.SetLevel(logrus.FatalLevel)
			} else if globalFlags.Debug {
				qlog.SetLevel(logrus.DebugLevel)
			}

			log.StartFileLogging()
			return nil
		},
	}
}

func newManCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates q's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcobra.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return cmd
}
