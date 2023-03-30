// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
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
	rootCmd.AddCommand(newCmdQuickon(f))
	// Add plugin commands
	rootCmd.AddCommand(newCmdExperimental(f))
	rootCmd.AddCommand(newManCmd())
	rootCmd.AddCommand(newCmdBugReport(f))

	args := os.Args
	if len(args) > 1 {
		pluginHandler := NewDefaultPluginHandler(common.ValidPrefixes)
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
				if err := HandlePluginCommand(pluginHandler, cmdPathPieces); err != nil {
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
			// TODO apply extra flags
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

type PluginHandler interface {
	Lookup(filename string) (string, bool)
	Execute(executablePath string, cmdArgs, environment []string) error
}

func NewDefaultPluginHandler(validPrefixes []string) *DefaultPluginHandler {
	return &DefaultPluginHandler{
		ValidPrefixes: validPrefixes,
	}
}

type DefaultPluginHandler struct {
	ValidPrefixes []string
}

// Lookup implements PluginHandler
func (h *DefaultPluginHandler) Lookup(filename string) (string, bool) {
	p, _ := os.LookupEnv("PATH")
	qbin := common.GetDefaultBinDir()
	if !strings.Contains(p, qbin) {
		os.Setenv("PATH", fmt.Sprintf("%v:%v", p, qbin))
	}
	for _, prefix := range h.ValidPrefixes {
		path, err := exec.LookPath(fmt.Sprintf("%s-%s", prefix, filename))
		if err != nil || len(path) == 0 {
			continue
		}
		return path, true
	}

	return "", false
}

// Execute implements PluginHandler
func (h *DefaultPluginHandler) Execute(executablePath string, cmdArgs, environment []string) error {
	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(executablePath, cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = environment
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}

	// invoke cmd binary relaying the environment and args given
	// append executablePath to cmdArgs, as execve will make first argument the "binary name".
	return syscall.Exec(executablePath, append([]string{executablePath}, cmdArgs...), environment) // #nosec
}

func HandlePluginCommand(pluginHandler PluginHandler, cmdArgs []string) error {
	var remainingArgs []string // all "non-flag" arguments
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, "-") {
			break
		}
		remainingArgs = append(remainingArgs, strings.ReplaceAll(arg, "-", "_"))
	}

	if len(remainingArgs) == 0 {
		// the length of cmdArgs is at least 1
		return fmt.Errorf("flags cannot be placed before plugin name: %s", cmdArgs[0])
	}

	foundBinaryPath := ""

	// attempt to find binary, starting at longest possible name with given cmdArgs
	for len(remainingArgs) > 0 {
		path, found := pluginHandler.Lookup(strings.Join(remainingArgs, "-"))
		if !found {
			remainingArgs = remainingArgs[:len(remainingArgs)-1]
			continue
		}

		foundBinaryPath = path
		break
	}

	if len(foundBinaryPath) == 0 {
		return nil
	}

	// invoke cmd binary relaying the current environment and args given
	if err := pluginHandler.Execute(foundBinaryPath, cmdArgs[len(remainingArgs):], os.Environ()); err != nil {
		return err
	}

	return nil
}

// func initConfig() {
// 	if globalFlags.ConfigPath == "" {
// 		globalFlags.ConfigPath = common.GetDefaultConfig()
// 	}
// 	viper.SetConfigFile(globalFlags.ConfigPath)
// 	viper.AutomaticEnv()
// 	viper.ReadInConfig()
// }
