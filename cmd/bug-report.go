package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/easysoft/qcadmin/common"
	"github.com/spf13/cobra"
)

func newCmdBugReport() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bug-report",
		Short: "Display system information for bug report",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bugReport()
		},
	}
	return cmd
}

func bugReport() error {
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
	return nil
}
