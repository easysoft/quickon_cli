// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/autodetect"
	"github.com/easysoft/qcadmin/internal/pkg/util/binfile"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/initsystem"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/excmd"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/expass"
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
	KubeClient     *k8s.Client
	Log            log.Logger
}

func NewCluster() *Cluster {
	return &Cluster{
		Metadata: types.Metadata{
			ClusterCidr:      "10.42.0.0/16",
			ServiceCidr:      "10.43.0.0/16",
			Network:          "flannel",
			QuchengVersion:   common.DefaultQuchengVersion,
			DisableIngress:   false,
			ImportDefaultApp: "zentao-open",
			ConsolePassword:  expass.PwGenAlphaNumSymbols(16),
		},
		M: new(syncmap.Map),
	}
}

func (p *Cluster) GetCreateNativeOptions() []types.Flag {
	return []types.Flag{
		{
			Name:  "kube-token",
			P:     &p.KubeToken,
			V:     p.KubeToken,
			Usage: "token to use for cluster authentication",
		},
		{
			Name:  "cluster-data-source",
			P:     &p.DataSource,
			V:     p.DataSource,
			Usage: "data source for cluster, default is sqlite",
		},
		{
			Name:  "console-password",
			P:     &p.ConsolePassword,
			V:     p.ConsolePassword,
			Usage: "qucheng console default password",
		},
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
			Usage:     "qucheng version",
		},
		{
			Name:  "domain",
			P:     &p.Domain,
			V:     p.Domain,
			Usage: "application custom domain name",
		},
		{
			Name:  "app",
			P:     &p.ImportDefaultApp,
			V:     p.ImportDefaultApp,
			Usage: "install default app",
		},
	}
}

func (p *Cluster) AddHelmRepo() error {
	output, err := qcexec.Command(os.Args[0], "experimental", "helm", "repo-add", "--name", common.DefaultHelmRepoName, "--url", common.GetChartRepo(p.QuchengVersion)).CombinedOutput()
	if err != nil {
		errmsg := string(output)
		if !strings.Contains(errmsg, "exists") {
			p.Log.Errorf("init qucheng install repo failed: %s", string(output))
			return err
		}
		p.Log.Warn("qucheng install repo  already exists")
	} else {
		p.Log.Done("init qucheng install repo success")
	}

	output, err = qcexec.Command(os.Args[0], "experimental", "helm", "repo-update").CombinedOutput()
	if err != nil {
		p.Log.Errorf("update qucheng install repo failed: %s", string(output))
		return err
	}
	p.Log.Done("update qucheng install repo success")
	return nil
}

func (p *Cluster) InitCluster() error {
	p.Status.Status = common.StatusCreating
	if err := p.InitK3sCluster(); err != nil {
		return err
	}
	p.Status.Status = common.StatusRunning
	// dataDir := common.GetDefaultDataDir()
	// templateVars := map[string]string{
	// 	"%{NAMESPACE}%": common.DefaultSystem,
	// }
	// if err := deploy.StageFunc(dataDir, templateVars); err != nil {
	// 	return err
	// }
	if p.Metadata.DisableIngress {
		p.Log.Warn("disable ingress controller")
	} else {
		p.Log.Debug("start deploy ingress plugins: nginx-ingress-controller")
		// localp, _ := pluginapi.GetMeta("ingress", "nginx-ingress-controller")
		// localp.Client = p.KubeClient
		// if err := localp.Install(); err != nil {
		// 	p.Log.Warnf("deploy ingress plugins: nginx-ingress-controller failed, reason: %v", err)
		// } else {
		// 	p.Log.Done("deployed ingress plugins: nginx-ingress-controller success")
		// }
		if err := qcexec.CommandRun(os.Args[0], "manage", "plugins", "enable", "ingress"); err != nil {
			p.Log.Errorf("deploy plugin ingress err: %v", err)
		} else {
			p.Log.Done("deployed operator plugins: cne-ingress success")
		}
	}
	return nil
}

func (p *Cluster) InitK3sCluster() error {
	p.Log.Debug("executing init k3s cluster logic...")
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
		p.Log.Errorf("create k3s service failed: %s", err)
		return err
	}
	// TODO fix reinstall typo
	_ = s.Uninstall()
	if err := s.Install(); err != nil {
		p.Log.Errorf("install k3s service failed: %s", err)
		return err
	}
	p.Log.Done("installed k3s service success")
	// Start k3s service.
	if err := s.Start(); err != nil {
		p.Log.Errorf("start k3s service failed: %s", err)
		return err
	}
	p.Log.Done("started k3s service success")
	if !excmd.CheckBin("kubectl") {
		os.Symlink(k3sbin, common.KubectlBinPath)
		p.Log.Done("create kubectl soft link")
	}
	p.Log.StartWait("waiting for k3s cluster to be ready...")
	t1 := time.Now()
	for {
		if file.CheckFileExists(common.K3sKubeConfig) {
			break
		}
		time.Sleep(time.Second * 5)
		p.Log.Info(".")
	}
	p.Log.StopWait()
	t2 := time.Now()
	p.Log.Donef("k3s cluster ready, cost: %v", t2.Sub(t1))
	d := common.GetDefaultKubeConfig()
	os.Symlink(common.K3sKubeConfig, d)
	p.Log.Donef("create kubeconfig soft link %v ---> %v", common.K3sKubeConfig, d)
	os.Rename(common.K3sDefaultDir, common.K3sDefaultDir+"_"+time.Now().Format("20060102150405"))
	os.Symlink(common.DefaultQuickonPlatformDir, common.K3sDefaultDir)
	p.Log.Donef("create kubeconfig soft link %v ---> %v", common.DefaultQuickonPlatformDir, common.K3sDefaultDir)
	kclient, _ := k8s.NewSimpleClient()
	if kclient != nil {
		_, err = kclient.CreateNamespace(context.TODO(), common.DefaultSystem, metav1.CreateOptions{})
		if err == nil {
			p.Log.Donef("create namespace %s", common.DefaultSystem)
		}
		p.KubeClient = kclient
	}
	return nil
}

func (p *Cluster) configCommonOptions() []string {
	var args []string
	if excmd.CheckBin("docker") {
		args = append(args, "--docker")
		// check docker  cgroup
		if autodetect.VerifyCgroupDriverSystemd() {
			args = append(args, "--kubelet-arg=cgroup-driver=systemd")
		}
	}
	// if len(p.EIP) != 0 {
	// 	args = append(args, fmt.Sprintf("--node-external-ip=%v", p.EIP))
	// }
	if len(p.KubeToken) != 0 {
		args = append(args, "--token="+p.KubeToken)
	}
	args = append(args, "--kubelet-arg=max-pods=220",
		"--kube-proxy-arg=proxy-mode=ipvs",
		"--kube-proxy-arg=masquerade-all=true",
		"--kube-proxy-arg=metrics-bind-address=0.0.0.0",
		"--data-dir="+common.DefaultQuickonPlatformDir,
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

	if len(p.DataSource) != 0 {
		if p.DataSource == "init-etcd" {
			args = append(args, "--cluster-init")
		} else if strings.HasPrefix(p.DataSource, "postgres") || strings.HasPrefix(p.DataSource, "mysql") {
			args = append(args, "--datastore-endpoint="+p.DataSource)
		} else if strings.HasPrefix(p.DataSource, "etcd") {
			args = append(args, "--datastore-endpoint="+strings.ReplaceAll(p.DataSource, "etcd://", "http://"))
		}
	}
	return args
}

func (p *Cluster) JoinCluster() error {
	p.Log.Debug("executing init k3s cluster logic...")
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
		p.Log.Errorf("create k3s agent failed: %s", err)
		return err
	}
	if err := s.Install(); err != nil {
		p.Log.Errorf("install k3s agent failed: %s", err)
		return err
	}
	p.Log.Done("installed k3s agent success")
	// Start k3s service.
	if err := s.Start(); err != nil {
		p.Log.Errorf("start k3s agent failed: %s", err)
		return err
	}
	p.Log.Done("started k3s agent success")
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
	p.Log.Debug("agent: %s, %s", p.CNEAPI, sever)
	if len(sever) > 0 {
		args = append(args, fmt.Sprintf("--server=https://%s:6443", sever))
	}
	token := p.getEnv(p.CNEToken, "CNE_TOKEN", "")
	p.Log.Debug("agent: %s, %s", p.CNEToken, token)
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
		p.Log.Error(err)
	}
}

func (p *Cluster) ready(ctx context.Context) error {
	t1 := ztime.NowUnix()
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG()).SetTimeout(time.Second * 1)
	p.Log.StartWait("waiting for qucheng ready")
	status := false
	for {
		t2 := ztime.NowUnix() - t1
		if t2 > 180 {
			p.Log.Warnf("waiting for qucheng ready 3min timeout: check your network or storage. after install you can run: q status")
			break
		}
		_, err := client.R().Get(fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0]))
		if err == nil {
			status = true
			break
		}
		time.Sleep(time.Second * 10)
	}
	p.Log.StopWait()
	if status {
		p.Log.Donef("qucheng ready, cost: %v", time.Since(time.Unix(t1, 0)))
	}
	return nil
}
