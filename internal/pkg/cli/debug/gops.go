// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package debug

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"

	egops "github.com/easysoft/qcadmin/internal/pkg/util/gops"
)

func GOPSCommand(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gops",
		Short: "gops is a tool to list and diagnose Go processes.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(os.Args) > 3 {
				_, err := strconv.Atoi(os.Args[3])
				if err == nil {
					f.GetLog().Infof("fetch pid %s process info", os.Args[3])
					egops.ProcessInfo(os.Args[3:]) // shift off the command name
					return
				}
			}
			f.GetLog().Info("fetch all process info")
			egops.Processes()
		},
	}
	cmd.AddCommand(treeCommand())
	cmd.AddCommand(processCommand())
	return cmd
}

// treeCommand displays a process tree.
func treeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tree",
		Short: "Display parent-child tree for Go processes.",
		Run: func(cmd *cobra.Command, args []string) {
			egops.DisplayProcessTree()
		},
	}
}

// processCommand displays information about a Go process.
func processCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "process",
		Aliases: []string{"pid", "proc"},
		Short:   "Prints information about a Go process.",
		RunE: func(cmd *cobra.Command, args []string) error {
			egops.ProcessInfo(args)
			return nil
		},
	}
}
