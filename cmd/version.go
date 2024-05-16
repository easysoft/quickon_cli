// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/cmd/version"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

// newCmdVersion show version
func newCmdVersion(f factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Args:  cobra.NoArgs,
		Run: func(cobraCmd *cobra.Command, args []string) {
			version.ShowVersion(f.GetLog())
		},
	}
}
