// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/ergoapi/util/zos"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/cmd/experimental"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

func newCmdExperimental(f factory.Factory) *cobra.Command {
	experimentalCmd := &cobra.Command{
		Use:     "experimental",
		Aliases: []string{"x", "exp"},
		Short:   "Experimental commands that may be modified or deprecated",
	}
	experimentalCmd.AddCommand(experimental.KubectlCommand(f))
	experimentalCmd.AddCommand(experimental.HelmCommand(f))
	experimentalCmd.AddCommand(experimental.ToolsCommand(f))
	experimentalCmd.AddCommand(experimental.SSHCommand(f))
	experimentalCmd.AddCommand(experimental.SCPCommand(f))
	experimentalCmd.AddCommand(experimental.K3sTPLCommand(f))
	experimentalCmd.AddCommand(experimental.DebugCommand(f))
	experimentalCmd.AddCommand(experimental.CheckCommand(f))
	if zos.IsLinux() {
		experimentalCmd.AddCommand(experimental.InstallCommand(f))
	}
	return experimentalCmd
}
