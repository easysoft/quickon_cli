// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/internal/pkg/status/top"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/pkg/cluster"
)

var (
	KRNodeExample = templates.Examples(`
	z cluster nodes list
	`)
)

func NewCmdClusterNodes(f factory.Factory) *cobra.Command {
	clusterNodesCmd := &cobra.Command{
		Use:     "nodes",
		Short:   "cluster nodes manage commands",
		Version: "20250211",
	}
	clusterNodesCmd.AddCommand(joinCommand(f))
	clusterNodesCmd.AddCommand(deleteCommand(f))
	clusterNodesCmd.AddCommand(topNodeCommand())
	return clusterNodesCmd
}

func joinCommand(f factory.Factory) *cobra.Command {
	myCluster := cluster.NewCluster(f)
	authStatus := myCluster.CheckAuthExist()
	join := &cobra.Command{
		Use:     "join",
		Short:   "join cluster",
		Aliases: []string{"add"},
		Example: templates.Examples(i18n.T(`
	# join cluster by pass4Quickon
	z cluster join --worker 192.168.99.52 --password pass4Quickon

	# join cluster by pkfile
	z cluster join --worker 192.168.99.52 --pkfile /root/.ssh/id_rsa
	`)),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !authStatus && (len(myCluster.SSH.Passwd) == 0 && len(myCluster.SSH.Pk) == 0) {
				return errors.New("missing ssh user or passwd or pk")
			}
			if myCluster.SSH.User != "root" {
				return errors.New("only support root user")
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

func deleteCommand(f factory.Factory) *cobra.Command {
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

func topNodeCommand() *cobra.Command {
	o := top.NodeOption{}
	nodeCmd := &cobra.Command{
		Use:                   "toplist",
		DisableFlagsInUseLine: true,
		Short:                 "show node top usage info",
		Example:               KRNodeExample,
		Run: func(cmd *cobra.Command, args []string) {
			o.Validate()
			o.RunResourceNode()
		},
	}
	nodeCmd.PersistentFlags().StringVarP(&o.KubeCtx, "context", "", "", "context to use for Kubernetes config")
	nodeCmd.PersistentFlags().StringVarP(&o.KubeConfig, "kubeconfig", "", "", "kubeconfig file to use for Kubernetes config")
	nodeCmd.PersistentFlags().StringVarP(&o.Output, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	nodeCmd.PersistentFlags().StringVarP(&o.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	nodeCmd.PersistentFlags().StringVarP(&o.SortBy, "sortBy", "s", "cpu", "sort by cpu or memory")
	return nodeCmd
}
