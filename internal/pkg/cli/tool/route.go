// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package tool

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/route"
)

var (
	routeHost      string
	routeGatewayIP string
)

func EmbedRouteCommand() *cobra.Command {
	// routeCmd represents the route command
	var routeCmd = &cobra.Command{
		Use:   "route",
		Short: "set default route gateway",
	}
	// check route for host
	routeCmd.PersistentFlags().StringVar(&routeHost, "host", "", "route host ip address for iFace")
	routeCmd.AddCommand(newCheckRouteCmd())
	routeCmd.AddCommand(newDelRouteCmd())
	routeCmd.AddCommand(newAddRouteCmd())
	return routeCmd
}

func newCheckRouteCmd() *cobra.Command {
	var checkRouteCmd = &cobra.Command{
		Use:   "check",
		Short: "check route host via gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			return route.CheckIsDefaultRoute(routeHost)
		},
	}
	return checkRouteCmd
}

func newAddRouteCmd() *cobra.Command {
	var addRouteCmd = &cobra.Command{
		Use:   "add",
		Short: "set route host via gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := route.NewRoute(routeHost, routeGatewayIP)
			return r.SetRoute()
		},
	}
	// manually to set host via gateway
	addRouteCmd.Flags().StringVar(&routeGatewayIP, "gateway", "", "route gateway ,ex ip route add host via gateway")
	return addRouteCmd
}

func newDelRouteCmd() *cobra.Command {
	var delRouteCmd = &cobra.Command{
		Use:   "del",
		Short: "del route host via gateway, like ip route del host via gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := route.NewRoute(routeHost, routeGatewayIP)
			return r.DelRoute()
		},
	}
	// manually to set host via gateway
	delRouteCmd.Flags().StringVar(&routeGatewayIP, "gateway", "", "route gateway ,ex ip route del host via gateway")
	return delRouteCmd
}
