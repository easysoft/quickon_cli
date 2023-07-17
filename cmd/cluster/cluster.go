// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/cmd/precheck"
	statussubcmd "github.com/easysoft/qcadmin/cmd/status"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/pkg/cluster"
	"github.com/ergoapi/util/confirm"
	"github.com/ergoapi/util/exnet"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	initExample = templates.Examples(`
		# init default cluster
		q cluster init

		# init cluster with custom cidr
		q cluster init --pod-cidr 10.100.0.0/16 --service-cidr 10.200.0.0/16

		# init cluster use mysql as datastore
		q cluster init --datastore mysql://root:123456@localhost:3306/k3s

		# more args
		q cluster init --help
	`)
)

// k3s server --tls-san --data-dir --cluster-cidr --service-cidr \
// --token --server --cluster-init --datastore-endpoint --disable  servicelb, traefik, local-storage
// --disable-network-policy --disable-helm-controller \
// --docker  \
// --pause-image \

func InitCommand(f factory.Factory) *cobra.Command {
	var preCheck precheck.PreCheck
	myCluster := cluster.NewCluster(f)
	init := &cobra.Command{
		Use:     "init",
		Short:   "init cluster",
		Example: initExample,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(myCluster.MasterIPs) == 0 {
				myCluster.MasterIPs = append(myCluster.MasterIPs, exnet.LocalIPs()[0])
			}
			// 禁止重复初始化
			if err := preCheck.Run(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return myCluster.InitNode()
		},
	}
	init.Flags().AddFlagSet(flags.ConvertFlags(init, myCluster.GetInitFlags()))
	return init
}

func JoinCommand(f factory.Factory) *cobra.Command {
	myCluster := cluster.NewCluster(f)
	authStatus := myCluster.CheckAuthExist()
	join := &cobra.Command{
		Use:   "join",
		Short: "join cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !authStatus && (len(myCluster.SSH.Passwd) == 0 && len(myCluster.SSH.Pk) == 0) {
				return errors.New("missing ssh user or passwd or pk")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return myCluster.JoinNode()
		},
	}
	fs := myCluster.GetIPFlags()
	if !authStatus {
		fs = append(fs, myCluster.GetSSHFlags()...)
	}
	join.Flags().AddFlagSet(flags.ConvertFlags(join, fs))
	return join
}

func DeleteCommand(f factory.Factory) *cobra.Command {
	myCluster := cluster.NewCluster(f)
	deleteCmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete node(s)",
		Aliases: []string{"del"},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(myCluster.IPs) == 0 {
				return errors.New("missing node ips")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return myCluster.DeleteNode()
		},
	}
	deleteCmd.Flags().StringSliceVar(&myCluster.IPs, "ips", nil, "ips, like 192.168.0.1:22")
	return deleteCmd
}

func CleanCommand(f factory.Factory) *cobra.Command {
	myCluster := cluster.NewCluster(f)
	log := f.GetLog()
	clean := &cobra.Command{
		Use:     "clean",
		Short:   "clean cluster",
		Version: "2.0.4",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, _ := confirm.Confirm("Are you sure to clean cluster")
			if status {
				return myCluster.Clean()
			}
			log.Donef("cancel clean cluster")
			return nil
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
