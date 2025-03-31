// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/expass"
	"github.com/ergoapi/util/exstr"
	"github.com/ergoapi/util/file"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/cli/k3stpl"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/downloader"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/pkg/util/ssh"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var defaultBackoff = wait.Backoff{
	Duration: 6 * time.Second,
	Factor:   1,
	Steps:    10,
}

type Cluster struct {
	log                   log.Logger
	MasterIPs             []string
	WorkerIPs             []string
	IPs                   []string
	SSH                   types.SSH
	CNI                   string
	DataDir               string
	PodCIDR               string
	ServiceCIDR           string
	OffLine               bool
	Registry              string
	Storage               string
	DataStore             string
	IgnorePreflightErrors bool
}

func NewCluster(f factory.Factory) *Cluster {
	return &Cluster{
		log:         f.GetLog(),
		CNI:         "host-gw",
		PodCIDR:     common.DefaultClusterPodCidr,
		ServiceCIDR: common.DefaultClusterServiceCidr,
		DataDir:     common.DefaultQuickonDataDir,
		SSH: types.SSH{
			User: "root",
		},
		DataStore:             "",
		Storage:               common.DefaultStorageType,
		Registry:              common.DefaultHub,
		OffLine:               false,
		IgnorePreflightErrors: false,
	}
}

func (c *Cluster) getInitFlags() []types.Flag {
	return []types.Flag{
		{
			Name:  "offline",
			P:     &c.OffLine,
			V:     c.OffLine,
			Usage: `offline install, only whitelist users are supported`,
		},
		{
			Name:   "hub",
			P:      &c.Registry,
			V:      c.Registry,
			EnvVar: common.DefaultHub,
			Usage:  `custom image hub`,
		},
		{
			Name:   "storage",
			P:      &c.Storage,
			V:      c.Storage,
			EnvVar: common.DefaultStorageType,
			Usage:  `storage, e.g: nfs,local,none`,
		},
		{
			Name:  "datastore",
			P:     &c.DataStore,
			V:     c.DataStore,
			Usage: `datastore, e.g: mysql://root:123456@tcp(localhost:3306)/k3s?charset=utf8&parseTime=True&loc=Local`,
		},
		{
			Name:  "ignore-preflight-errors",
			P:     &c.IgnorePreflightErrors,
			V:     c.IgnorePreflightErrors,
			Usage: `ignore precheck error`,
		},
	}
}

func (c *Cluster) GetInitFlags() []types.Flag {
	fs := c.GetSSHFlags()
	fs = append(fs, c.GetIPFlags()...)
	fs = append(fs, c.getMasterFlags()...)
	fs = append(fs, c.getInitFlags()...)
	return fs
}

func (c *Cluster) GetSSHFlags() []types.Flag {
	return []types.Flag{
		{
			Name:      "username",
			ShortHand: "u",
			P:         &c.SSH.User,
			V:         c.SSH.User,
			Usage:     "ssh user, only support root",
		},
		{
			Name:  "password",
			P:     &c.SSH.Passwd,
			V:     c.SSH.Passwd,
			Usage: "ssh password",
		},
		{
			Name:  "pkfile",
			P:     &c.SSH.Pk,
			V:     c.SSH.Pk,
			Usage: "ssh private key, if not set, will use password",
		},
		{
			Name:  "pkpass",
			P:     &c.SSH.PkPasswd,
			V:     c.SSH.PkPasswd,
			Usage: "ssh private key password",
		},
	}
}

func (c *Cluster) GetIPFlags() []types.Flag {
	return []types.Flag{
		{
			Name:  "master",
			P:     &c.MasterIPs,
			V:     c.MasterIPs,
			Usage: `master ip list, e.g: 192.168.0.1:22`,
		},
		{
			Name:  "worker",
			P:     &c.WorkerIPs,
			V:     c.WorkerIPs,
			Usage: `worker ip list, e.g: 192.168.0.1:22`,
		},
	}
}

func (c *Cluster) getMasterFlags() []types.Flag {
	return []types.Flag{
		{
			Name:  "cni",
			P:     &c.CNI,
			V:     c.CNI,
			Usage: "k8s networking plugin, support flannel, wireguard, host-gw, custom",
		},
		{
			Name:  "pod-cidr",
			P:     &c.PodCIDR,
			V:     c.PodCIDR,
			Usage: "k8s cluster pod cidr",
		},
		{
			Name:  "service-cidr",
			P:     &c.ServiceCIDR,
			V:     c.ServiceCIDR,
			Usage: "k8s cluster service cidr",
		},
		{
			Name:      "data-dir",
			P:         &c.DataDir,
			V:         c.DataDir,
			ShortHand: "d",
			Usage:     "cluster & quickon data dir",
		},
	}
}

func (c *Cluster) GetMasterFlags() []types.Flag {
	fs := c.GetSSHFlags()
	fs = append(fs, c.GetIPFlags()...)
	fs = append(fs, c.getMasterFlags()...)
	return fs
}

func (c *Cluster) GetWorkerFlags() []types.Flag {
	return nil
}

func (c *Cluster) preinit(mip, ip string, sshClient ssh.Interface) error {
	k3sbin := fmt.Sprintf("%s/hack/bin/k3s-%s-%s", common.GetDefaultDataDir(), runtime.GOOS, runtime.GOARCH)
	if err := sshClient.Copy(ip, k3sbin, common.K3sBinPath); err != nil {
		return errors.Errorf("copy k3s bin (%s:%s -> %s:%s) failed, reason: %v", ip, mip, k3sbin, common.K3sBinPath, ip, err)
	}
	qbin, _ := os.Executable()
	if qbin != common.QcAdminBinPath || mip != ip {
		if err := sshClient.Copy(ip, qbin, common.QcAdminBinPath); err != nil {
			return errors.Errorf("copy cli bin (%s:%s -> %s:%s) failed, reason: %v", mip, qbin, ip, common.QcAdminBinPath, err)
		}
	}

	if err := sshClient.CmdAsync(ip, fmt.Sprintf("%s version", common.QcAdminBinPath)); err != nil {
		return errors.Errorf("load cli version failed, reason: %v", err)
	}
	c.log.StartWait(ip + " start run init script")
	initArgs := []string{common.GetCustomFile("hack/manifests/scripts/init.sh")}
	if c.OffLine {
		initArgs = append(initArgs, "offline")
	}
	if err := sshClient.CmdAsync(ip, strings.Join(initArgs, " ")); err != nil {
		return errors.Errorf("%s run init script failed, reason: %v", ip, err)
	}
	c.log.StopWait()
	c.log.Donef("%s run init script success", ip)
	// add master0 ip
	hostsArgs := fmt.Sprintf("%s exp tools hosts add --domain %s --ip %s", common.QcAdminBinPath, common.DefaultKubeAPIDomain, mip)
	if err := sshClient.CmdAsync(ip, hostsArgs); err != nil {
		c.log.Debugf("cmd: %s", hostsArgs)
		return errors.Errorf("%s add master0 (%s --> %s) failed, reason: %v", ip, common.DefaultKubeAPIDomain, mip, err)
	}
	if err := sshClient.CmdAsync(ip, common.GetCustomFile("hack/manifests/scripts/node.sh")); err != nil {
		return errors.Errorf("%s run init script failed, reason: %v", ip, err)
	}
	return nil
}

func (c *Cluster) initMaster0(cfg *config.Config, sshClient ssh.Interface) error {
	c.log.Infof("master0 ip: %s", cfg.Cluster.InitNode)
	k3sargs := k3stpl.K3sArgs{
		Master0:      true,
		TypeMaster:   true,
		KubeAPI:      common.DefaultKubeAPIDomain,
		KubeToken:    expass.PwGenAlphaNum(16),
		DataDir:      c.DataDir,
		PodCIDR:      c.PodCIDR,
		ServiceCIDR:  c.ServiceCIDR,
		CNI:          c.CNI,
		DataStore:    c.DataStore,
		LocalStorage: false,
		// LocalStorage: strings.ToLower(c.Storage) == "local",
		Registry:  c.Registry,
		OffLine:   c.OffLine,
		Master0IP: cfg.Cluster.InitNode,
	}
	master0tplSrc := fmt.Sprintf("%s/master0.%s", common.GetDefaultCacheDir(), cfg.Cluster.InitNode)
	master0tplDst := fmt.Sprintf("/%s/.k3s.service", c.SSH.User)
	file.WriteFile(master0tplSrc, k3sargs.Manifests(""), true)
	if err := sshClient.Copy(cfg.Cluster.InitNode, master0tplSrc, master0tplDst); err != nil {
		return errors.Errorf("copy master0 %s tpl failed, reason: %v", cfg.Cluster.InitNode, err)
	}
	if err := c.preinit(cfg.Cluster.InitNode, cfg.Cluster.InitNode, sshClient); err != nil {
		return err
	}
	// waiting k3s ready
	if err := c.waitk3sReady(cfg.Cluster.InitNode, sshClient); err != nil {
		return err
	}
	c.log.Infof("install %s as default storage", c.Storage)
	if c.Storage == "nfs" {
		if err := qcexec.CommandRun("bash", "-c", common.GetCustomFile("hack/manifests/storage/nfs-server.sh")); err != nil {
			return errors.Errorf("%s run install nfs script failed, reason: %v", cfg.Cluster.InitNode, err)
		}
	} else if c.Storage == "local" {
		kubeargs := []string{"experimental", "kubectl", "apply", "-f", common.GetCustomFile("hack/manifests/storage/local.yaml")}
		if err := qcexec.CommandRun(os.Args[0], kubeargs...); err != nil {
			return errors.Errorf("%s run install local storage failed, reason: %v", cfg.Cluster.InitNode, err)
		}
	} else if c.Storage == "none" {
		c.log.Infof("skip install storage")
	}
	kclient, _ := k8s.NewSimpleClient()
	if ns, _ := kclient.GetNamespace(context.TODO(), common.DefaultKubeSystem, metav1.GetOptions{}); ns != nil {
		cfg.Cluster.ID = string(ns.GetUID())
	}
	if err := c.installNerdctl(); err != nil {
		c.log.Warnf("install nerdctl failed, after install cluster, you can use `%s` retry install nerdctl", color.SGreen("%s exp install nerdctl", os.Args[0]))
	}
	cfg.Cluster.PodCIDR = c.PodCIDR
	cfg.Cluster.ServiceCIDR = c.ServiceCIDR
	cfg.Cluster.CNI = c.CNI
	cfg.Cluster.Registry = c.Registry
	cfg.Storage.Type = c.Storage
	cfg.DB = k3sargs.DataStore
	cfg.DataDir = k3sargs.DataDir
	cfg.Cluster.Master = append(cfg.Cluster.Master, config.Node{
		Host: cfg.Cluster.InitNode,
		Init: true,
	})
	cfg.Cluster.Token = k3sargs.KubeToken
	if c.OffLine {
		cfg.Install.Type = "offline"
		cfg.Install.Pkg = common.GetDefaultDataDir()
	} else {
		cfg.Install.Type = "online"
	}
	cfg.Install.Version = common.Version
	return cfg.SaveConfig()
}

func (c *Cluster) waitk3sReady(host string, sshClient ssh.Interface) error {
	c.log.StartWait("check k3s ready.")
	try := 0
	err := wait.ExponentialBackoff(defaultBackoff, func() (bool, error) {
		try++
		c.log.Debugf("the %d/%d time tring to check k3s status", try, defaultBackoff.Steps)
		err := sshClient.Copy(host, common.K3sKubeConfig, common.DefaultKubeConfig())
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	c.log.StopWait()
	if err != nil {
		return errors.Errorf("check k3s ready failed, reason: %w", err)
	}
	c.log.Done("check k3s ready.")
	return nil
}

func (c *Cluster) installNerdctl() error {
	if c.OffLine {
		return nil
	}
	if !file.CheckFileExists(common.DefaultNerdctlConfig) {
		os.MkdirAll(common.DefaultNerdctlDir, common.FileMode0777)
		file.Copy(common.GetCustomFile("hack/manifests/hub/nerdctl.toml"), common.DefaultNerdctlConfig, true)
	}
	remoteURL := fmt.Sprintf("https://pkg.zentao.net/qucheng/cli/stable/tools/nerdctl-%s-%s", runtime.GOOS, runtime.GOARCH)
	localURL := fmt.Sprintf("%s/qc-nerdctl", common.GetDefaultBinDir())
	_, err := downloader.Download(remoteURL, localURL)
	if err != nil {
		return err
	}
	_ = os.Chmod(localURL, common.FileMode0755)
	docker := fmt.Sprintf("%s/qc-docker", common.GetDefaultBinDir())
	_ = os.Remove(docker)
	_ = os.Link(localURL, docker)
	c.log.Donef("install nerdctl success")
	return nil
}

func (c *Cluster) joinNode(ip string, master bool, cfg *config.Config, sshClient ssh.Interface) error {
	t := "worker"
	if master {
		t = "master"
	}
	c.log.Infof("node role: %s, ip: %s", t, ip)
	k3sargs := k3stpl.K3sArgs{
		Master0:      false,
		TypeMaster:   master,
		CNI:          cfg.Cluster.CNI,
		KubeAPI:      common.DefaultKubeAPIDomain,
		KubeToken:    cfg.Cluster.Token,
		DataDir:      cfg.DataDir,
		PodCIDR:      cfg.Cluster.PodCIDR,
		ServiceCIDR:  cfg.Cluster.ServiceCIDR,
		DataStore:    cfg.DB,
		LocalStorage: true,
	}
	tplSrc := fmt.Sprintf("%s/%s.%s", common.GetDefaultCacheDir(), t, ip)
	tplDst := fmt.Sprintf("/%s/.k3s.service", c.SSH.User)
	file.WriteFile(tplSrc, k3sargs.Manifests(""), true)
	if err := sshClient.Copy(ip, tplSrc, tplDst); err != nil {
		return errors.Errorf("%s copy tpl (%s:%s->%s:%s) failed, reason: %v", t, cfg.Cluster.InitNode, tplSrc, ip, tplDst, err)
	}
	if err := c.preinit(cfg.Cluster.InitNode, ip, sshClient); err != nil {
		return err
	}
	if master {
		cfg.Cluster.Master = append(cfg.Cluster.Master, config.Node{
			Host: ip,
		})
	} else {
		cfg.Cluster.Worker = append(cfg.Cluster.Worker, config.Node{
			Host: ip,
		})
	}
	if cfg.Global.SSH.User == "" {
		cfg.Global.SSH = c.SSH
	}
	return cfg.SaveConfig()
}

func (c *Cluster) CheckNodeInitStatus(master, node string, sshClient ssh.Interface) error {
	// TODO 检查是否已经初始化过
	c.log.Infof("check node %s init status", node)
	return nil
}

func (c *Cluster) InitNode() error {
	c.log.Info("start init cluster node")
	c.MasterIPs = exstr.DuplicateStrElement(c.MasterIPs)
	c.WorkerIPs = exstr.DuplicateStrElement(c.WorkerIPs)
	otherMaster := c.MasterIPs[1:]
	sshClient := ssh.NewSSHClient(&c.SSH, true)
	cfg := config.LoadTruncateConfig()
	cfg.Cluster.InitNode = c.MasterIPs[0]
	if err := c.initMaster0(cfg, sshClient); err != nil {
		return err
	}
	for _, host := range otherMaster {
		c.log.Debugf("ping master %s", host)
		if err := sshClient.Ping(host); err != nil {
			c.log.Warnf("skip join master: %s, reason: %v", host, err)
			continue
		}
		if err := c.joinNode(host, true, cfg, sshClient); err != nil {
			c.log.Warnf("skip join master: %s, reason: %v", host, err)
		}
	}
	for _, host := range c.WorkerIPs {
		c.log.Debugf("ping worker %s", host)
		if err := sshClient.Ping(host); err != nil {
			c.log.Warnf("skip join worker: %s, reason: %v", host, err)
			continue
		}
		if err := c.joinNode(host, false, cfg, sshClient); err != nil {
			c.log.Warnf("skip join worker: %s, reason: %v", host, err)
		}
	}
	return nil
}

func (c *Cluster) CheckAuthExist() bool {
	cfg, _ := config.LoadConfig()
	if cfg.Global.SSH.Passwd == "" || cfg.Global.SSH.Pk == "" || cfg.Global.SSH.PkData == "" {
		return false
	}
	if cfg.Global.SSH.User != "root" {
		return false
	}
	return true
}

func (c *Cluster) JoinNode() error {
	c.log.Info("start join node")
	c.MasterIPs = exstr.DuplicateStrElement(c.MasterIPs)
	c.WorkerIPs = exstr.DuplicateStrElement(c.WorkerIPs)
	sshClient := ssh.NewSSHClient(&c.SSH, true)
	cfg, _ := config.LoadConfig()
	for _, host := range c.MasterIPs {
		c.log.Debugf("check master available %s via ssh", host)
		if err := sshClient.Ping(host); err != nil {
			c.log.Warnf("skip join master: %s, reason: %v", host, err)
			continue
		}
		c.log.Donef("check master available %s via ssh success", host)
		if err := c.joinNode(host, true, cfg, sshClient); err != nil {
			c.log.Warnf("skip join master: %s, reason: %v", host, err)
			continue
		}
		c.log.Donef("join master %s success", host)
	}
	for _, host := range c.WorkerIPs {
		c.log.Debugf("check worker available %s via ssh", host)
		if err := sshClient.Ping(host); err != nil {
			c.log.Warnf("skip join worker: %s, reason: %v", host, err)
			continue
		}
		c.log.Donef("check worker available %s via ssh success", host)
		if err := c.joinNode(host, false, cfg, sshClient); err != nil {
			c.log.Warnf("skip join worker: %s, reason: %v", host, err)
			continue
		}
		c.log.Donef("join worker %s success", host)
	}
	return nil
}

func (c *Cluster) cleanNode(ip string, sshClient ssh.Interface, wg *sync.WaitGroup) {
	defer wg.Done()
	c.log.StartWait(fmt.Sprintf("start clean node: %s", ip))
	err := sshClient.CmdAsync(ip, common.GetCustomFile("hack/manifests/scripts/cleankube.sh"))
	c.log.StopWait()
	if err != nil {
		c.log.Warnf("clean node %s failed, reason: %v", ip, err)
		return
	}
	c.log.Donef("clean node %s success", ip)
}

func (c *Cluster) deleteNode(ip string, sshClient ssh.Interface, kubeClient *k8s.Client, wg *sync.WaitGroup) error {
	c.log.Infof("start clean node %s", ip)
	// 从集群中移除节点
	c.log.Infof("delete node %s from cluster", ip)
	if err := kubeClient.DownNode(context.TODO(), ip); err != nil {
		c.log.Warnf("delete node %s from cluster failed, reason: %v", ip, err)
	}
	// 清理节点
	c.cleanNode(ip, sshClient, wg)
	return nil
}

func (c *Cluster) DeleteNode() error {
	cfg, _ := config.LoadConfig()
	var wg sync.WaitGroup
	sshClient := ssh.NewSSHClient(&cfg.Global.SSH, true)
	kubeClient, err := k8s.NewSimpleClient(common.GetKubeConfig())
	if err != nil {
		return errors.Errorf("load kube client failed, reason: %v", err)
	}
	for _, ip := range c.IPs {
		if ip == cfg.Cluster.InitNode {
			c.log.Warnf("init node %s not allow delete, can use clean subcmd", ip)
			continue
		}
		wg.Add(1)
		c.deleteNode(ip, sshClient, kubeClient, &wg)
	}
	wg.Wait()
	return cfg.SaveConfig()
}

// Clean 清理集群
func (c *Cluster) Clean() error {
	c.log.Info("start clean cluster")
	cfg, _ := config.LoadConfig()
	ips := cfg.GetIPs()
	if len(ips) == 0 {
		ips = append(ips, exnet.LocalIPs()[0])
	}
	sshClient := ssh.NewSSHClient(&cfg.Global.SSH, true)
	var wg sync.WaitGroup
	for _, ip := range ips {
		c.log.Debugf("clean node %s", ip)
		wg.Add(1)
		go c.cleanNode(ip, sshClient, &wg)
	}
	wg.Wait()
	c.log.Done("clean cluster success")
	return nil
}

// Stop 关闭集群
func (c *Cluster) Stop() error {
	c.log.Info("start stop cluster")
	cfg, _ := config.LoadConfig()
	ips := cfg.GetIPs()
	if len(ips) == 0 {
		ips = append(ips, exnet.LocalIPs()[0])
	}
	sshClient := ssh.NewSSHClient(&cfg.Global.SSH, true)
	var wg sync.WaitGroup
	for _, ip := range ips {
		c.log.Debugf("stop node %s", ip)
		wg.Add(1)
		go c.actionNode(ip, "stop", common.GetCustomFile("hack/manifests/scripts/stopnode.sh"), sshClient, &wg)
	}
	wg.Wait()
	c.log.Done("stop cluster success")
	return nil
}

// StartUP 启动集群
func (c *Cluster) StartUP() error {
	c.log.Info("startup cluster")
	cfg, _ := config.LoadConfig()
	ips := cfg.GetIPs()
	if len(ips) == 0 {
		ips = append(ips, exnet.LocalIPs()[0])
	}
	sshClient := ssh.NewSSHClient(&cfg.Global.SSH, true)
	var wg sync.WaitGroup
	for _, ip := range ips {
		c.log.Debugf("startup node %s", ip)
		wg.Add(1)
		go c.actionNode(ip, "startup", common.GetCustomFile("hack/manifests/scripts/startupnode.sh"), sshClient, &wg)
	}
	wg.Wait()
	c.log.Done("startup cluster success")
	return nil
}

func (c *Cluster) actionNode(ip, action, script string, sshClient ssh.Interface, wg *sync.WaitGroup) {
	defer wg.Done()
	c.log.StartWait(fmt.Sprintf("start %s node: %s", action, ip))
	err := sshClient.CmdAsync(ip, script)
	c.log.StopWait()
	if err != nil {
		c.log.Warnf("%s node %s failed, reason: %v", action, ip, err)
		return
	}
	c.log.Donef("%s node %s success", action, ip)
}
