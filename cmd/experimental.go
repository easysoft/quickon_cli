// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/cmd/experimental"
	"github.com/spf13/cobra"
)

func NewCmdExperimental() *cobra.Command {
	experimentalCmd := &cobra.Command{
		Use:     "experimental",
		Aliases: []string{"x", "exp"},
		Short:   "Experimental commands that may be modified or deprecated",
		Hidden:  true,
	}
	experimentalCmd.AddCommand(experimental.KubectlCommand())
	experimentalCmd.AddCommand(experimental.HelmCommand())
	experimentalCmd.AddCommand(experimental.ToolsCommand())
	return experimentalCmd
}
