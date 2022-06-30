// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/spf13/cobra"
)

func newCmdBugReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bug-report",
		Short: "Display system information for bug report",
		Long:  "this command shares no personally-identifiable information, and is unused unless you share the bug identifier with our team.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bugReport()
		},
	}
	return cmd
}

func bugReport() error {
	log := log.GetInstance()
	log.Info("Issue: 🐛Bug Report: https://github.com/easysoft/qucheng_cli/issues/new?assignees=&labels=&template=bug-report.md")
	// TODO 详细信息
	sprintf := func(key, val string) string {
		return fmt.Sprintf("%-24s%s\n", key, val)
	}
	report := sprintf("q version:", common.Version)
	report += sprintf("GOOS:", runtime.GOOS)
	report += sprintf("GOARCH:", runtime.GOARCH)
	report += sprintf("NumCPU:", fmt.Sprint(runtime.NumCPU()))
	vcs, ok := debug.ReadBuildInfo()
	if ok && vcs != nil {
		report += fmt.Sprintln("Build info:")
		vcs := *vcs
		report += sprintf("\tGo version:", vcs.GoVersion)
		report += sprintf("\tModule path:", vcs.Path)
		report += sprintf("\tMain version:", vcs.Main.Version)
		report += sprintf("\tMain path:", vcs.Main.Path)
		report += sprintf("\tMain checksum:", vcs.Main.Sum)

		report += fmt.Sprintln("\tBuild settings:")
		for _, set := range vcs.Settings {
			report += sprintf(fmt.Sprintf("\t\t%s:", set.Key), set.Value)
		}
	}
	fmt.Println(report)
	return nil
}
