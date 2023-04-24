// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/cmd/flags"
	statussubcmd "github.com/easysoft/qcadmin/cmd/status"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/pkg/cluster"
	"github.com/ergoapi/util/exnet"
	"github.com/spf13/cobra"
)

const initExample = `q cluster init --podsubnet "10.42.0.0/16" \
 			--svcsubnet "10.43.0.0/16" \
			--eip 1.1.1.1  \
			--san kubeapi.k8s.io`

// k3s server --tls-san --data-dir --cluster-cidr --service-cidr \
// --token --server --cluster-init --datastore-endpoint --disable  servicelb, traefik, local-storage
// --disable-network-policy --disable-helm-controller \
// --docker  \
// --pause-image \

func InitCommand(f factory.Factory) *cobra.Command {
	cluster := cluster.NewCluster(f)
	init := &cobra.Command{
		Use:     "init",
		Short:   "init cluster",
		Example: initExample,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(cluster.MasterIPs) == 0 {
				cluster.MasterIPs = append(cluster.MasterIPs, exnet.LocalIPs()[0])
			}
			// 禁止重复初始化
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.InitNode()
		},
	}
	init.Flags().AddFlagSet(flags.ConvertFlags(init, cluster.GetInitFlags()))
	return init
}

func JoinCommand(f factory.Factory) *cobra.Command {
	cluster := cluster.NewCluster(f)
	authStatus := cluster.CheckAuthExist()
	join := &cobra.Command{
		Use:   "join",
		Short: "join cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !authStatus && (len(cluster.SSH.Passwd) == 0 && len(cluster.SSH.Pk) == 0) {
				return errors.New("missing ssh user or passwd or pk")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.JoinNode()
		},
	}
	fs := cluster.GetIPFlags()
	if !authStatus {
		fs = append(fs, cluster.GetSSHFlags()...)
	}
	join.Flags().AddFlagSet(flags.ConvertFlags(join, fs))
	return join
}

func DeleteCommand(f factory.Factory) *cobra.Command {
	cluster := cluster.NewCluster(f)
	delete := &cobra.Command{
		Use:     "delete",
		Short:   "delete node(s)",
		Aliases: []string{"del"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(cluster.IPs) == 0 {
				return errors.New("missing node ips")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.DeleteNode()
		},
	}
	delete.Flags().StringSliceVar(&cluster.IPs, "ips", nil, "ips, like 192.168.0.1:22")
	return delete
}

func CleanCommand(f factory.Factory) *cobra.Command {
	cluster := cluster.NewCluster(f)
	clean := &cobra.Command{
		Use:   "clean",
		Short: "clean cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.Clean()
		},
	}
	return clean
}

func StatusCommand(f factory.Factory) *cobra.Command {
	status := &cobra.Command{
		Use:   "status",
		Short: "status cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	status.AddCommand(statussubcmd.TopNodeCmd())
	return status
}

func StorageCommand(f factory.Factory) *cobra.Command {
	storage := &cobra.Command{
		Use:     "csi",
		Aliases: []string{"storage"},
		Short:   "cluster csi",
		Hidden:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return storage
}
