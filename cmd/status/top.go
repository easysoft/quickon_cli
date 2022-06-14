// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package status

import (
	"github.com/easysoft/qcadmin/internal/pkg/status/top"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	KRNodeExample = templates.Examples(`
	q status node
	`)
)

func TopNodeCmd() *cobra.Command {
	o := top.NodeOption{}
	nodeCmd := &cobra.Command{
		Use:                   "node",
		DisableFlagsInUseLine: true,
		Short:                 "node provides an overview of the node",
		Aliases:               []string{"nodes", "no"},
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
