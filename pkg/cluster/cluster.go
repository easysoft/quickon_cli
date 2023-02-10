// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"fmt"
	"runtime"

	"github.com/cockroachdb/errors"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/cli/k3stpl"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/pkg/util/ssh"
	"github.com/ergoapi/util/exstr"
	"github.com/ergoapi/util/file"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

type Cluster struct {
	log       log.Logger
	MasterIPs []string
	WorkerIPs []string
	SSH       types.SSH
}

func NewCluster(f factory.Factory) *Cluster {
	return &Cluster{
		log: f.GetLog(),
	}
}

func (c *Cluster) copyk3s(mip, ip string, sshClient ssh.Interface) error {
	k3sbin := fmt.Sprintf("%s/hack/bin/k3s-%s-%s", common.GetDefaultDataDir(), runtime.GOOS, runtime.GOARCH)
	if err := sshClient.Copy(ip, k3sbin, "/usr/local/bin/k3s"); err != nil {
		return errors.Errorf("copy k3s bin (%s:%s -> %s:/usr/local/bin/k3s) failed, reason: %v", ip, mip, k3sbin, ip, err)
	}
	return nil
}

func (c *Cluster) initMaster0(ip string, sshClient ssh.Interface) error {
	c.log.Infof("master0 ip: %s", ip)
	k3sargs := k3stpl.K3sArgs{}
	master0tplSrc := fmt.Sprintf("%s/master0.%s", common.GetDefaultCacheDir(), ip)
	master0tplDst := fmt.Sprintf("%s/k3s.service", c.SSH.User)
	file.Writefile(master0tplSrc, k3sargs.Manifests(""))
	if err := sshClient.Copy(ip, master0tplSrc, master0tplDst); err != nil {
		return errors.Errorf("copy master0 %s tpl failed, reason: %v", ip, err)
	}
	return c.copyk3s(ip, ip, sshClient)
}

func (c *Cluster) joinNode(mip, ip string, master bool, sshClient ssh.Interface) error {
	t := "worker"
	if master {
		t = "master"
	}
	c.log.Infof("node role: %s, ip: %s", t, ip)
	k3sargs := k3stpl.K3sArgs{}
	tplSrc := fmt.Sprintf("%s/%s.%s", common.GetDefaultCacheDir(), t, ip)
	tplDst := fmt.Sprintf("/%s/k3s.service", c.SSH.User)
	file.Writefile(tplSrc, k3sargs.Manifests(""))
	if err := sshClient.Copy(ip, tplSrc, tplDst); err != nil {
		return errors.Errorf("%s copy tpl (%s:%s->%s:%s) failed, reason: %v", t, mip, tplSrc, ip, tplDst, err)
	}
	return c.copyk3s(mip, ip, sshClient)
}

func (c *Cluster) InitNode() error {
	c.log.Info("init node")
	c.MasterIPs = exstr.DuplicateStrElement(c.MasterIPs)
	c.WorkerIPs = exstr.DuplicateStrElement(c.WorkerIPs)
	master0 := c.MasterIPs[0]
	otherMaster := c.MasterIPs[1:]
	sshClient := ssh.NewSSHClient(&c.SSH, true)
	if err := c.initMaster0(master0, sshClient); err != nil {
		return err
	}
	for _, host := range otherMaster {
		c.log.Debugf("ping master %s", host)
		if err := sshClient.Ping(host); err != nil {
			c.log.Warnf("skip join master: %s, reason: %v", host, err)
			continue
		}
		if err := c.joinNode(master0, host, true, sshClient); err != nil {
			c.log.Warnf("skip join master: %s, reason: %v", host, err)
		}
	}
	for _, host := range c.WorkerIPs {
		c.log.Debugf("ping worker %s", host)
		if err := sshClient.Ping(host); err != nil {
			c.log.Warnf("skip join worker: %s, reason: %v", host, err)
			continue
		}
		if err := c.joinNode(master0, host, false, sshClient); err != nil {
			c.log.Warnf("skip join worker: %s, reason: %v", host, err)
		}
	}
	return nil
}
