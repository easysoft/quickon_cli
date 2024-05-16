// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package debug

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

func PortForwardCommand(f factory.Factory) *cobra.Command {
	var ns, svc string
	var port int
	pf := &cobra.Command{
		Use:     "port-forward",
		Aliases: []string{"pf"},
		Short:   "forward local port to kube pod or svc",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k8s.PortForwardCommand(context.Background(), ns, svc, port)
		},
	}
	pf.Flags().StringVarP(&ns, "ns", "n", "default", "namespace")
	pf.Flags().StringVar(&svc, "svc", "kubernetes", "svc name")
	pf.Flags().IntVarP(&port, "port", "p", 443, "port")
	return pf
}
