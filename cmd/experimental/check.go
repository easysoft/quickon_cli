// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package experimental

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/cli/check"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

func CheckCommand(f factory.Factory) *cobra.Command {
	tCmd := &cobra.Command{
		Use:     "check",
		Short:   "check tools",
		Long:    "检查数据库等中间件是否可用",
		Version: "4.0.0",
	}
	tCmd.AddCommand(check.CheckMySQLCommand(f))
	return tCmd
}
