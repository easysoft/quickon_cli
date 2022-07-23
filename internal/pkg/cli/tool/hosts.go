// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package tool

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/hosts"
	"github.com/spf13/cobra"
)

var hostsPath string

func EmbedHostsCommand(f factory.Factory) *cobra.Command {
	var hostsCmd = &cobra.Command{
		Use:   "hosts",
		Short: "hosts manager",
	}
	// check route for host
	hostsCmd.AddCommand(newHostsListCmd())
	hostsCmd.AddCommand(newHostsAddCmd(f))
	hostsCmd.AddCommand(newHostsDeleteCmd(f))
	hostsCmd.PersistentFlags().StringVar(&hostsPath, "path", "/etc/hosts", "default hosts path")
	return hostsCmd
}

func newHostsListCmd() *cobra.Command {
	var hostsListCmd = &cobra.Command{
		Use:   "list",
		Short: "hosts manager list",
		Run: func(cmd *cobra.Command, args []string) {
			hf := &hosts.HostFile{Path: hostsPath}
			hf.ListCurrentHosts()
		},
	}
	return hostsListCmd
}

func newHostsAddCmd(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var ip, domain string
	var hostsAddCmd = &cobra.Command{
		Use:   "add",
		Short: "hosts manager add",
		PreRun: func(cmd *cobra.Command, args []string) {
			if ip == "" {
				log.Fatal("ip not empty")
			}
			if domain == "" {
				log.Fatal("domain not empty")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			hf := &hosts.HostFile{Path: hostsPath}
			if hf.HasDomain(domain) {
				hf.DeleteDomain(domain)
				log.Donef("domain %s delete success", domain)
			}
			hf.AppendHost(domain, ip)
			log.Donef("domain %s:%s append success", domain, ip)
		},
	}
	hostsAddCmd.Flags().StringVar(&ip, "ip", "", "ip address")
	hostsAddCmd.Flags().StringVar(&domain, "domain", "", "domain address")

	return hostsAddCmd
}

func newHostsDeleteCmd(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var domain string
	var hostsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "hosts manager delete",
		PreRun: func(cmd *cobra.Command, args []string) {
			if domain == "" {
				log.Fatal("domain not empty")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			hf := &hosts.HostFile{Path: hostsPath}
			hf.DeleteDomain(domain)
			log.Donef("domain %s delete success", domain)
		},
	}
	hostsDeleteCmd.Flags().StringVar(&domain, "domain", "", "domain address")

	return hostsDeleteCmd
}
