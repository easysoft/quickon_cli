// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	pluginapi "github.com/easysoft/qcadmin/internal/pkg/plugin"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/binfile"
	"github.com/easysoft/qcadmin/internal/pkg/util/initsystem"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/static/deploy"
	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/excmd"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/imroc/req/v3"
	"github.com/kardianos/service"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/syncmap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Cluster struct {
	types.Metadata `json:",inline"`
	types.Status   `json:"status"`
	M              *sync.Map
	KubeClient         *k8s.Client
}

func NewCluster() *Cluster {
	return &Cluster{
		Metadata: types.Metadata{
			ClusterCidr:    "10.42.0.0/16",
			ServiceCidr:    "10.43.0.0/16",
			Network:        "flannel",
			QuchengVersion: "stable",
			DisableIngress: false,
		},
		M: new(syncmap.Map),
	}
}

func (p *Cluster) GetCreateOptions() []types.Flag {
	return []types.Flag{
		// {
		// 	Name:      "plugins",
		// 	P:         &p.Plugins,
		// 	V:         p.Plugins,
		// 	ShortHand: "p",
		// 	Usage:     "Deploy packaged components",
		// },
		{
			Name:  "disable-ingress",
			P:     &p.DisableIngress,
			V:     p.DisableIngress,
			Usage: "disable nginx ingress plugins",
		},
		{
			Name:  "podsubnet",
			P:     &p.ClusterCidr,
			V:     p.ClusterCidr,
			Usage: "pod subnet",
		},
		{
			Name:  "svcsubnet",
			P:     &p.ServiceCidr,
			V:     p.ServiceCidr,
			Usage: "service subnet",
		},
		{
			Name:  "eip",
			P:     &p.EIP,
			V:     p.EIP,
			Usage: "external IP addresses to advertise for node",
		},
		{
			Name:  "san",
			P:     &p.TLSSans,
			V:     p.TLSSans,
			Usage: "kube api custom domain",
		},
		{
			Name:  "network",
			P:     &p.Network,
			V:     p.Network,
			Usage: "network cni",
		},
	}
}

func (p *Cluster) GetJoinOptions() []types.Flag {
	return []types.Flag{
		{
			Name:  "cne-api",
			P:     &p.CNEAPI,
			V:     p.CNEAPI,
			Usage: "Server to connect t",
		}, {
			Name:  "cne-token",
			P:     &p.CNEToken,
			V:     p.CNEToken,
			Usage: "Token to use for authentication",
		},
	}
}

func (p *Cluster) GetCreateExtOptions() []types.Flag {
	return []types.Flag{
		{
			Name:      "qucheng-version",
			P:         &p.QuchengVersion,
			V:         p.QuchengVersion,
			ShortHand: "q",
			Usage:     "qucheng version: stable, test",
		},
	}
}

func (p *Cluster) InitCluster() error {
	p.Status.Status = common.StatusCreating
	if err := p.InitK3sCluster(); err != nil {
		return err
	}
	p.Status.Status = common.StatusRunning
	dataDir := common.GetDefaultDataDir()
	templateVars := map[string]string{
		"%{NAMESPACE}%": common.DefaultSystem,
	}
	if err := deploy.StageFunc(dataDir, templateVars); err != nil {
		return err
	}
	if p.Metadata.DisableIngress {
		log.Flog.Warn("disable ingress controller")
	} else {
		log.Flog.Debug("start deploy ingress plugins: nginx-ingress-controller")
		localp, _ := pluginapi.GetMeta("ingress", "nginx-ingress-controller")
		localp.Client = p.KubeClient
		if err := localp.Install(); err != nil {
			log.Flog.Warnf("deploy ingress plugins: nginx-ingress-controller failed, reason: %v", err)
		} else {
			log.Flog.Done("deployed ingress plugins: nginx-ingress-controller success")
		}
	}
	return nil
}

func (p *Cluster) InitK3sCluster() error {
	log.Flog.Debug("executing init k3s cluster logic...")
	// Download k3s.
	getbin := binfile.Meta{}
	k3sbin, err := getbin.LoadLocalBin(common.K3sBinName)
	if err != nil {
		return err
	}
	// k3s args
	k3sargs := []string{
		"server",
	}
	// common args
	k3sargs = append(k3sargs, p.configCommonOptions()...)
	// k3s server config
	k3sargs = append(k3sargs, p.configServerOptions()...)
	// Create k3s service.
	k3sCfg := &initsystem.Config{
		Name: "k3s",
		Desc: "k3s server",
		Exec: k3sbin,
		Args: k3sargs,
	}
	options := make(service.KeyValue)
	options["Restart"] = "always"
	options["LimitNOFILE"] = 1048576
	options["LimitNPROC"] = "infinity"
	options["LimitCORE"] = "infinity"
	options["TasksMax"] = "infinity"
	options["TimeoutStartSec"] = 0
	options["RestartSec"] = "5s"
	options["Type"] = "notify"
	options["KillMode"] = "process"
	options["Delegate"] = true
	svcConfig := &service.Config{
		Name:        k3sCfg.Name,
		DisplayName: k3sCfg.Name,
		Description: k3sCfg.Desc,
		Dependencies: []string{
			"After=network-online.target",
		},
		Executable: k3sCfg.Exec,
		Arguments:  k3sCfg.Args,
		Option:     options,
		ExecStartPres: []string{
			"/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'",
			"/sbin/modprobe br_netfilter",
			"/sbin/modprobe overlay",
		},
	}
	ds := new(initsystem.DaemonService)
	s, err := service.New(ds, svcConfig)
	if err != nil {
		log.Flog.Errorf("create k3s service failed: %s", err)
		return err
	}
	if err := s.Install(); err != nil {
		log.Flog.Errorf("install k3s service failed: %s", err)
		return err
	}
	log.Flog.Done("installed k3s service success")
	// Start k3s service.
	if err := s.Start(); err != nil {
		log.Flog.Errorf("start k3s service failed: %s", err)
		return err
	}
	log.Flog.Done("started k3s service success")
	if !excmd.CheckBin("kubectl") {
		os.Symlink(k3sbin, common.KubectlBinPath)
		log.Flog.Done("create kubectl soft link")
	}
	log.Flog.StartWait("waiting for k3s cluster to be ready...")
	t1 := time.Now()
	for {
		if file.CheckFileExists(common.K3sKubeConfig) {
			break
		}
		time.Sleep(time.Second * 5)
		log.Flog.Info(".")
	}
	log.Flog.StopWait()
	t2 := time.Now()
	log.Flog.Donef("k3s cluster ready, cost: %v", t2.Sub(t1))
	d := common.GetDefaultKubeConfig()
	os.Symlink(common.K3sKubeConfig, d)
	log.Flog.Donef("create kubeconfig soft link %v ---> %v/config", common.K3sKubeConfig, d)
	kclient, _ := k8s.NewSimpleClient()
	if kclient != nil {
		_, err = kclient.CreateNamespace(context.TODO(), common.DefaultSystem, metav1.CreateOptions{})
		if err == nil {
			log.Flog.Donef("create namespace %s", common.DefaultSystem)
		}
		p.KubeClient = kclient
	}
	return nil
}

func (p *Cluster) configCommonOptions() []string {
	var args []string
	if excmd.CheckBin("docker") {
		args = append(args, "--docker")
	}
	if len(p.EIP) != 0 {
		args = append(args, fmt.Sprintf("--node-external-ip=%v", p.EIP))
	}
	args = append(args, "--kubelet-arg=max-pods=220",
		"--kube-proxy-arg=proxy-mode=ipvs",
		"--kube-proxy-arg=masquerade-all=true",
		"--kube-proxy-arg=metrics-bind-address=0.0.0.0",
		// "--token=a1b2c3d4", // TODO 随机生成
		"--pause-image=hub.qucheng.com/library/k3s-pause:3.6")

	return args
}

func (p *Cluster) configServerOptions() []string {
	/*
		--tls-san
		--cluster-cidr
		--service-cidr
		--service-node-port-range
		--flannel-backend
		--token
		--datastore-endpoint
		--disable-network-policy
		--disable-helm-controller
		--docker
		--pause-image
		--node-external-ip
		--kubelet-arg
		--flannel-backend=none
	*/
	var args []string
	args = append(args, "--disable-network-policy", "--disable-helm-controller", "--disable=servicelb,traefik")
	var tlsSans string
	for _, tlsSan := range p.TLSSans {
		tlsSans = tlsSans + fmt.Sprintf(" --tls-san=%s", tlsSan)
	}
	tlsSans = tlsSans + " --tls-san=kapi.qucheng.local"
	if len(p.EIP) != 0 {
		tlsSans = tlsSans + fmt.Sprintf(" --tls-san=%s", p.EIP)
	}
	if len(tlsSans) != 0 {
		args = append(args, tlsSans)
	}
	if p.Network != "flannel" {
		args = append(args, "--flannel-backend=none")
	}
	args = append(args, "--service-node-port-range=22767-32767", "--system-default-registry=hub.qucheng.com/library")
	args = append(args, fmt.Sprintf("--cluster-cidr=%v", p.ClusterCidr))
	args = append(args, fmt.Sprintf("--service-cidr=%v", p.ServiceCidr))
	// args = append(args, fmt.Sprintf("--cluster-dns=%v", p.DnSSvcIP))
	// if len(p.Token) != 0 {
	// 	args = append(args, "--token="+p.Token)
	// }
	// args = append(args, p.Args...)
	return args
}

func (p *Cluster) JoinCluster() error {
	log.Flog.Debug("executing init k3s cluster logic...")
	// Download k3s.
	getbin := binfile.Meta{}
	k3sbin, err := getbin.LoadLocalBin(common.K3sBinName)
	if err != nil {
		return err
	}
	// k3s args
	k3sargs := []string{
		"agent",
	}
	// common args
	k3sargs = append(k3sargs, p.configCommonOptions()...)
	// k3s agent config
	k3sargs = append(k3sargs, p.configAgentOptions()...)
	// Create k3s service.
	k3sCfg := &initsystem.Config{
		Name: "k3s",
		Desc: "k3s agent",
		Exec: k3sbin,
		Args: k3sargs,
	}
	options := make(service.KeyValue)
	options["Restart"] = "always"
	options["LimitNOFILE"] = 1048576
	options["LimitNPROC"] = "infinity"
	options["LimitCORE"] = "infinity"
	options["TasksMax"] = "infinity"
	options["TimeoutStartSec"] = 0
	options["RestartSec"] = "5s"
	options["Type"] = "exec"
	options["KillMode"] = "process"
	options["Delegate"] = true
	svcConfig := &service.Config{
		Name:        k3sCfg.Name,
		DisplayName: k3sCfg.Name,
		Description: k3sCfg.Desc,
		Dependencies: []string{
			"After=network-online.target",
		},
		Executable: k3sCfg.Exec,
		Arguments:  k3sCfg.Args,
		Option:     options,
		ExecStartPres: []string{
			"/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'",
			"/sbin/modprobe br_netfilter",
			"/sbin/modprobe overlay",
		},
	}
	ds := new(initsystem.DaemonService)
	s, err := service.New(ds, svcConfig)
	if err != nil {
		log.Flog.Errorf("create k3s agent failed: %s", err)
		return err
	}
	if err := s.Install(); err != nil {
		log.Flog.Errorf("install k3s agent failed: %s", err)
		return err
	}
	log.Flog.Done("installed k3s agent success")
	// Start k3s service.
	if err := s.Start(); err != nil {
		log.Flog.Errorf("start k3s agent failed: %s", err)
		return err
	}
	log.Flog.Done("started k3s agent success")
	return nil
}

// TODO support agent install
func (p *Cluster) configAgentOptions() []string {
	// agent
	/*
		--token
		--server
		--docker
		--pause-image
		--node-external-ip
		--kubelet-arg
	*/
	var args []string

	sever := p.getEnv(p.CNEAPI, "CNE_API", "")
	log.Flog.Debug("agent: %s, %s", p.CNEAPI, sever)
	if len(sever) > 0 {
		args = append(args, fmt.Sprintf("--server=https://%s:6443", sever))
	}
	token := p.getEnv(p.CNEToken, "CNE_TOKEN", "")
	log.Flog.Debug("agent: %s, %s", p.CNEToken, token)
	if len(token) > 0 {
		args = append(args, "--token="+token)
	}
	return args
}

func (p *Cluster) getEnv(key, envkey, defaultvalue string) string {
	if len(key) > 0 {
		return key
	}
	return environ.GetEnv(envkey, defaultvalue)
}

// Ready 渠成Ready
func (p *Cluster) Ready() {
	clusterWaitGroup, ctx := errgroup.WithContext(context.Background())
	clusterWaitGroup.Go(func() error {
		return p.ready(ctx)
	})
	if err := clusterWaitGroup.Wait(); err != nil {
		log.Flog.Error(err)
	}
}

func (p *Cluster) ready(ctx context.Context) error {
	t1 := ztime.NowUnix()
	client := req.C().SetUserAgent(common.GetUG()).SetTimeout(time.Second * 1)
	log.Flog.StartWait("waiting for qucheng ready")
	status := false
	for {
		t2 := ztime.NowUnix() - t1
		if t2 > 180 {
			log.Flog.Warnf("waiting for qucheng ready 3min timeout: check your network or storage. after install you can run: q status")
			break
		}
		_, err := client.R().Get(fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0]))
		if err == nil {
			status = true
			break
		}
		time.Sleep(time.Second * 10)
	}
	log.Flog.StopWait()
	if status {
		log.Flog.Donef("qucheng ready, cost: %v", time.Since(time.Unix(t1, 0)))
	}
	return nil
}
