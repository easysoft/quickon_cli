// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/internal/pkg/cluster"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func newCmdUninstall(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("start uninstall cluster")
			c := cluster.NewCluster()
			c.Log = log
			err := c.Uninstall()
			if err != nil {
				log.Fatalf("uninstall cluster failed, reason: %v", err)
			}
			log.Info("uninstall cluster success")
		},
	}
}
