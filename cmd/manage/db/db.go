// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package db

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

func NewCmdDB(f factory.Factory) *cobra.Command {
	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Manage Platform Databases",
	}
	dbListCmd := &cobra.Command{
		Use:   "list",
		Short: "List Platform DB or DbService",
	}
	dbListCmd.AddCommand(cmdDbsList(f))
	dbListCmd.AddCommand(cmdDbSvcList(f))
	dbCmd.AddCommand(dbListCmd)
	dbCmd.AddCommand(cmdExternalDb(f))
	return dbCmd
}
