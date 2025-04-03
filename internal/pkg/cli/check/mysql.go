// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package check

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/db"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

// CheckMySQLCommand check mysql command.
func CheckMySQLCommand(f factory.Factory) *cobra.Command {
	var (
		host     string
		port     int
		username string
		password string
	)
	cmd := &cobra.Command{
		Use:   "mysql",
		Short: "check mysql availability and database creation",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if password == "" {
				return fmt.Errorf("password is required")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &db.Config{
				Host:     host,
				Port:     port,
				Username: username,
				Password: password,
			}
			db.CheckMySQL(cfg)
		},
	}
	cmd.Flags().StringVarP(&host, "host", "", "localhost", "MySQL host")
	cmd.Flags().IntVarP(&port, "port", "", 3306, "MySQL port")
	cmd.Flags().StringVarP(&username, "username", "u", "root", "MySQL username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "MySQL password")
	return cmd
}
