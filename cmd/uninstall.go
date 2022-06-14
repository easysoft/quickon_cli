// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/internal/pkg/cluster"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/spf13/cobra"
)

func newCmdUninstall() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall",
		Run: func(cmd *cobra.Command, args []string) {
			log.Flog.Info("start uninstall cluster")
			c := cluster.NewCluster()
			err := c.Uninstall()
			if err != nil {
				log.Flog.Fatalf("uninstall cluster failed, reason: %v", err)
			}
			log.Flog.Info("uninstall cluster success")
		},
	}
}
