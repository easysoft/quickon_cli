// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
)

type bugReportCmd struct {
	log log.Logger
}

func newCmdBugReport(f factory.Factory) *cobra.Command {
	br := bugReportCmd{
		log: f.GetLog(),
	}
	cmd := &cobra.Command{
		Use:     "bugreport",
		Aliases: []string{"bug-report"},
		Short:   "Display system information for bug report",
		Long:    "this command shares no personally-identifiable information, and is unused unless you share the bug identifier with our team.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return br.BugReport()
		},
	}
	return cmd
}

func (br bugReportCmd) BugReport() error {
	debugShell := fmt.Sprintf("%s/hack/manifests/scripts/diagnose.sh", common.GetDefaultDataDir())
	br.log.Debugf("gen debug message script: %v", debugShell)
	// 移除qcadmin初始化文件
	if err := qcexec.CommandRun("/bin/bash", debugShell, os.Args[0]); err != nil {
		return err
	}
	bugMsg := "found bug: submit the error message to Github or Gitee\n\t Github: https://github.com/easysoft/quickon_cli/issues/new?assignees=&labels=&template=bug-report.md\n\t Gitee: https://gitee.com/wwccss/qucheng_cli/issues\n"
	br.log.Info(bugMsg)
	return nil
}
