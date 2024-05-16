// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package sshutil

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/ssh"
)

func EmbedSSHCommand(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var cfg types.SSH
	var hosts []string
	ssh := &cobra.Command{
		Use:   "ssh",
		Short: "ssh util tool",
		Run: func(cmd *cobra.Command, args []string) {
			osargs := args[:]
			if len(hosts) == 0 {
				// local
				log.Warnf("hosts is null, local exec: %v", osargs)
			} else {
				sshClient := ssh.NewSSHClient(&cfg, true)
				for _, host := range hosts {
					if err := sshClient.Ping(host); err != nil {
						continue
					}
					sshClient.CmdsAsync(host, osargs...)
				}
			}
		},
	}
	ssh.Flags().StringVar(&cfg.User, "user", "root", "ssh user")
	ssh.Flags().StringVar(&cfg.Passwd, "passwd", "", "ssh password")
	ssh.Flags().StringVar(&cfg.Pk, "pkfile", "", "ssh pk file")
	ssh.Flags().StringVar(&cfg.PkPasswd, "pkpass", "", "ssh key passwd")
	ssh.Flags().StringSliceVar(&hosts, "hosts", nil, "ips, like 192.168.0.1:22")
	return ssh
}

func EmbedScpCommand(f factory.Factory) *cobra.Command {
	var cfg types.SSH
	var source, dest string
	var hosts []string
	log := f.GetLog()
	ssh := &cobra.Command{
		Use:   "scp",
		Short: "scp util tool",
		Run: func(cmd *cobra.Command, args []string) {
			if len(hosts) == 0 {
				// local
				log.Warnf("hosts is null")
			} else {
				sshClient := ssh.NewSSHClient(&cfg, true)
				for _, host := range hosts {
					if err := sshClient.Ping(host); err != nil {
						continue
					}
					if err := sshClient.Copy(host, source, dest); err != nil {
						log.Errorf("%s scp copy file %s -> %s err: %v", host, source, dest, err)
					}
				}
			}
		},
	}
	ssh.Flags().StringVar(&cfg.User, "user", "root", "ssh user")
	ssh.Flags().StringVar(&cfg.Passwd, "passwd", "", "ssh password")
	ssh.Flags().StringVar(&cfg.Pk, "pkfile", "", "ssh pk file")
	ssh.Flags().StringVar(&cfg.PkPasswd, "pkpass", "", "ssh key passwd")
	ssh.Flags().StringSliceVar(&hosts, "hosts", nil, "ips, like 192.168.0.1:22")
	ssh.Flags().StringVar(&source, "source", "", "source")
	ssh.Flags().StringVar(&dest, "dest", "", "dest")
	return ssh
}
