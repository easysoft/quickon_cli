// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/easysoft/qcadmin/cmd/cluster"
	"github.com/easysoft/qcadmin/cmd/precheck"
	"github.com/easysoft/qcadmin/cmd/storage"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func newCmdPreCheck(f factory.Factory) *cobra.Command {
	var pc precheck.PreCheck
	cmd := &cobra.Command{
		Use:   "precheck",
		Short: "Precheck system",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pc.Run()
		},
		Args: cobra.NoArgs,
	}
	cmd.PersistentFlags().BoolVar(&pc.IgnorePreflightErrors, "ignore", false, "ignore precheck error")
	return cmd
}

func newCmdCluster(f factory.Factory) *cobra.Command {
	clusterCmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Cluster commands",
		Version: "20230330",
	}
	clusterCmd.AddCommand(newCmdPreCheck(f))
	clusterCmd.AddCommand(cluster.InitCommand(f))
	clusterCmd.AddCommand(cluster.JoinCommand(f))
	clusterCmd.AddCommand(cluster.DeleteCommand(f))
	clusterCmd.AddCommand(cluster.CleanCommand(f))
	clusterCmd.AddCommand(cluster.StatusCommand(f))
	clusterCmd.AddCommand(storage.NewCmdStorage(f))
	return clusterCmd
}
