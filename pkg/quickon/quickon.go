// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package quickon

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/expass"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/imroc/req/v3"
	"golang.org/x/sync/errgroup"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/pkg/util/retry"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	suffixdomain "github.com/easysoft/qcadmin/pkg/qucheng/domain"
	kubeerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Meta struct {
	Domain          string
	IP              string
	Version         string
	ConsolePassword string
	DevopsMode      bool
	OffLine         bool
	DBReplication   bool
	ExtDBHost       string
	ExtDBPort       int
	ExtDBUser       string
	ExtDBPassword   string
	Type            string
	App             string
	kubeClient      *k8s.Client
	Log             log.Logger
	DomainType      string
	UsePHP7         bool
}

func New(f factory.Factory) *Meta {
	return &Meta{
		Log: f.GetLog(),
		// Version:         common.DefaultQuickonOSSVersion,
		ConsolePassword: expass.PwGenAlphaNum(32),
		Type:            common.ZenTaoOSSType.String(),
		DomainType:      "custom",
	}
}

func (m *Meta) GetCustomFlags() []types.Flag {
	return []types.Flag{
		{
			Name:  "domain",
			Usage: "custom domain",
			P:     &m.Domain,
			V:     m.Domain,
		},
		{
			Name:   "ip",
			Usage:  "custom ip",
			P:      &m.IP,
			V:      m.IP,
			Hidden: true,
		},
		{
			Name:   "offline",
			Usage:  "offline install mode, default: false",
			P:      &m.OffLine,
			V:      false,
			Hidden: true,
		},
		{
			Name:  "db-replication",
			Usage: "db use replication mode, default standalone: false",
			P:     &m.DBReplication,
			V:     false,
		},
		{
			Name:  "ext-db-host",
			Usage: "external db host",
			P:     &m.ExtDBHost,
			V:     "",
		},
		{
			Name:  "ext-db-port",
			Usage: "external db port",
			P:     &m.ExtDBPort,
			V:     3306,
		},
		{
			Name:  "ext-db-user",
			Usage: "external db user",
			P:     &m.ExtDBUser,
			V:     "root",
		},
		{
			Name:  "ext-db-password",
			Usage: "external db password",
			P:     &m.ExtDBPassword,
			V:     "",
		},
	}
}

func (m *Meta) GetKubeClient() error {
	kubeClient, err := k8s.NewSimpleClient(common.GetKubeConfig())
	if err != nil {
		return errors.Errorf("kube client create failed: %v", err)
	}
	m.kubeClient = kubeClient
	return nil
}

func (m *Meta) checkIngress() {
	m.Log.StartWait("check default ingress class")
	defaultClass, _ := m.kubeClient.ListDefaultIngressClass(context.Background(), metav1.ListOptions{})
	m.Log.StopWait()
	if defaultClass == nil {
		m.Log.Infof("not found default ingress class, will install nginx ingress")
		m.Log.Debug("start install default ingress: nginx-ingress-controller")
		if err := qcexec.CommandRun(os.Args[0], "platform", "plugins", "enable", "ingress"); err != nil {
			m.Log.Errorf("install ingress failed, reason: %v", err)
		} else {
			m.Log.Done("install ingress: cne-ingress success")
		}
	} else {
		m.Log.Infof("found exist default ingress class: %s", defaultClass.Name)
	}
	m.Log.Done("check default ingress done")
}

func (m *Meta) checkStorage() {
	m.Log.StartWait("check default storage class")
	defaultClass, _ := m.kubeClient.GetDefaultSC(context.Background())
	m.Log.StopWait()
	if defaultClass == nil {
		// default storage
		cfg, _ := config.LoadConfig()
		storage := common.DefaultStorageType
		cfg.Storage.Type = storage
		defer cfg.SaveConfig()
		m.Log.Warnf("not found default storage class, will install %s as default storage", storage)
		m.Log.Debugf("start install default storage: %s", storage)
		if len(cfg.Cluster.InitNode) == 0 {
			cfg.Cluster.InitNode = exnet.LocalIPs()[0]
		}
		// old use nfs as default storage os.Args[0], "cluster", "storage", "nfs", "--ip", cfg.Cluster.InitNode, "--path", common.GetDefaultNFSStoragePath(cfg.DataDir)
		if err := qcexec.CommandRun(os.Args[0], "cluster", "storage", "local"); err != nil {
			m.Log.Errorf("install storage %s failed, reason: %v", storage, err)
		} else {
			m.Log.Donef("install storage %s success", storage)
		}
		// if err := qcexec.CommandRun(os.Args[0], "cluster", "storage", "set-default"); err != nil {
		// 	m.Log.Errorf("set default storageclass failed, reason: %v", err)
		// }
	} else {
		m.Log.Infof("found exist default storage class: %s", defaultClass.Name)
	}
	m.Log.Done("check default storage done")
}

func (m *Meta) CheckInstall() bool {
	_, err := config.LoadConfig()
	if err != nil {
		return false
	}
	_, err = m.kubeClient.GetDeployment(context.Background(), common.GetDefaultSystemNamespace(true), common.GetReleaseName(m.DevopsMode), metav1.GetOptions{})
	if err == nil {
		m.Log.Debug("found exist quickon deployment")
		return true
	}
	return false
}

func (m *Meta) Check() error {
	if err := m.addHelmRepo(); err != nil {
		return err
	}
	if err := m.initNS(); err != nil {
		return err
	}
	m.checkIngress()
	m.checkStorage()
	return nil
}

func (m *Meta) initNS() error {
	m.Log.Debugf("init platform default namespace")
	for _, ns := range common.GetDefaultQuickONNamespace() {
		_, err := m.kubeClient.GetNamespace(context.TODO(), ns, metav1.GetOptions{})
		if err != nil {
			if !kubeerr.IsNotFound(err) {
				return err
			}
			if _, err := m.kubeClient.CreateNamespace(context.TODO(), ns, metav1.CreateOptions{}); err != nil && kubeerr.IsAlreadyExists(err) {
				return err
			}
		}
	}
	m.Log.Donef("init quickon default namespace success")
	return nil
}

func (m *Meta) addHelmRepo() error {
	output, err := qcexec.Command(os.Args[0], "experimental", "helm", "repo-add", "--name", common.DefaultHelmRepoName, "--url", common.GetChartRepo(m.Version)).CombinedOutput()
	if err != nil {
		errmsg := string(output)
		if !strings.Contains(errmsg, "exists") {
			m.Log.Errorf("init quickon helm repo failed, reason: %s", string(output))
			return err
		}
		m.Log.Debugf("quickon helm repo already exists")
	} else {
		m.Log.Donef("add %s channel quickon helm repo success", common.GetChannel(m.Version))
	}
	output, err = qcexec.Command(os.Args[0], "experimental", "helm", "repo-update").CombinedOutput()
	if err != nil {
		m.Log.Errorf("update quickon helm repo failed, reason: %s", string(output))
		return err
	}
	m.Log.Done("update quickon helm repo success")
	return nil
}

func (m *Meta) Init() error {
	cfg, _ := config.LoadConfig()
	m.Log.Info("executing init logic...")
	ctx := context.Background()
	m.Log.Debug("waiting for storage to be ready...")
	waitsc := time.Now()
	// wait.BackoffUntil TODO
	scName := ""
	for {
		sc, _ := m.kubeClient.GetDefaultSC(ctx)
		if sc != nil {
			m.Log.Donef("default storage %s is ready", sc.Name)
			scName = sc.Name
			break
		}
		time.Sleep(time.Second * 5)
		trywaitsc := time.Now()
		if trywaitsc.Sub(waitsc) > time.Minute*3 {
			m.Log.Warnf("wait storage ready, timeout: %v", trywaitsc.Sub(waitsc).Seconds())
			break
		}
	}

	_, err := m.kubeClient.CreateNamespace(ctx, common.GetDefaultSystemNamespace(true), metav1.CreateOptions{})
	if err != nil {
		if !kubeerr.IsAlreadyExists(err) {
			return err
		}
	}
	installVersion := common.GetVersion(m.DevopsMode, m.Type, m.Version)
	repoChannel := common.GetChannel(m.Version)
	m.Log.Infof("devops: %v, type: %s, version: %s, channel: %s", m.DevopsMode, m.Type, installVersion, repoChannel)
	m.Log.Debugf("start init %s", common.GetInstallType(m.DevopsMode))
	cfg.Quickon.Type = common.QuickonType(m.Type)
	cfg.Quickon.DevOps = m.DevopsMode

	if m.Domain == "" {
		err := retry.Retry(time.Second*1, 3, func() (bool, error) {
			domain, _, err := m.genSuffixHTTPHost(m.IP)
			if err != nil {
				return false, err
			}
			m.Domain = domain
			m.DomainType = "api"
			m.Log.Infof("generate suffix domain: %s, ip: %v", color.SGreen(m.Domain), color.SGreen(m.IP))
			return true, nil
		})
		if err != nil {
			m.Domain = "demo.corp.cc"
			m.Log.Warnf("gen suffix domain failed, reason: %v, use default domain: %s", err, m.Domain)
		}
		if kutil.IsLegalDomain(m.Domain) {
			m.Log.Infof("load %s tls cert", m.Domain)
			defaultTLS := fmt.Sprintf("%s/tls-haogs-cn.yaml", common.GetDefaultCacheDir())
			m.Log.StartWait(fmt.Sprintf("start issuing domain %s certificate, may take 3-5min", m.Domain))
			waittls := time.Now()
			for {
				if file.CheckFileExists(defaultTLS) {
					m.Log.StopWait()
					m.Log.Done("detect tls cert file success")
					if err := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", defaultTLS, "-n", common.GetDefaultSystemNamespace(true), "--kubeconfig", common.GetKubeConfig()).Run(); err != nil {
						m.Log.Warnf("load default tls cert failed, reason: %v", err)
					} else {
						m.Log.Done("load default tls cert success")
					}
					qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", defaultTLS, "-n", "default", "--kubeconfig", common.GetKubeConfig()).Run()
					break
				}
				_, mainDomain := kutil.SplitDomain(m.Domain)
				domainTLS := fmt.Sprintf("https://pkg.zentao.net/ssl/%s/%s/tls.yaml", mainDomain, m.Domain)
				qcexec.Command(os.Args[0], "experimental", "tools", "wget", "-t", domainTLS, "-d", defaultTLS).Run()
				m.Log.Debug("wait for tls cert ready...")
				time.Sleep(time.Second * 5)
				trywaitsc := time.Now()
				if trywaitsc.Sub(waittls) >= time.Minute*5 {
					// TODO  timeout
					m.Log.Warnf("wait tls cert ready, timeout: %v", trywaitsc.Sub(waittls).Seconds())
					cmd := fmt.Sprintf("%s pt tls", os.Args[0])
					m.Log.Warnf("wait cluster install success, visit %s notice 'Your connection to this site isn't secure', please use follow cmd check and fix: %s", color.SGreen(m.Domain), color.SGreen(cmd))
					break
				}
			}
		} else {
			m.Log.Infof("use custom domain %s, you need to add two A records for your domain: %s -> %s, *.%s -> %s", m.Domain, color.SGreen(m.Domain), color.SGreen(m.IP), color.SGreen(m.Domain), color.SGreen(m.IP))
		}
	} else {
		m.Log.Infof("use custom domain %s, you need to add two A records for your domain: %s -> %s, *.%s -> %s", m.Domain, color.SGreen(m.Domain), color.SGreen(m.IP), color.SGreen(m.Domain), color.SGreen(m.IP))
	}
	token := expass.PwGenAlphaNum(32)

	cfg.Domain = m.Domain
	cfg.APIToken = token
	cfg.Quickon.Domain.Name = m.Domain
	cfg.Quickon.Domain.Type = m.DomainType
	cfg.S3.Username = expass.PwGenAlphaNum(8)
	cfg.S3.Password = expass.PwGenAlphaNum(16)
	cfg.SaveConfig()
	m.Log.Info("start deploy custom tools")
	// toolargs := []string{"experimental", "helm", "upgrade", "--name", "selfcert", "--repo", common.DefaultHelmRepoName, "--chart", "selfcert", "--namespace", common.GetDefaultSystemNamespace(true)}
	// if helmstd, err := qcexec.Command(os.Args[0], toolargs...).CombinedOutput(); err != nil {
	// 	m.Log.Warnf("deploy custom tools err: %v, std: %s", err, string(helmstd))
	// } else {
	// 	m.Log.Done("deployed custom tools success")
	// }
	m.Log.Info("start deploy operator")
	operatorargs := []string{"experimental", "helm", "upgrade", "--name", common.DefaultCneOperatorName, "--repo", common.DefaultHelmRepoName, "--chart", common.DefaultCneOperatorName, "--namespace", common.GetDefaultSystemNamespace(true),
		"--set", "minio.ingress.enabled=true",
		"--set", "minio.ingress.host=s3." + m.Domain,
		"--set", "minio.auth.username=" + cfg.S3.Username,
		"--set", "minio.auth.password=" + cfg.S3.Password,
	}
	if helmstd, err := qcexec.Command(os.Args[0], operatorargs...).CombinedOutput(); err != nil {
		m.Log.Warnf("deploy operator err: %v, std: %s", err, string(helmstd))
	} else {
		m.Log.Done("deployed operator success")
	}
	// check operator ready
	if err := m.OperatorReady(); err != nil {
		m.Log.Warnf("check operator ready failed, reason: %v", err)
	} else {
		m.Log.Done("check operator ready success")
	}
	useExtDB := false
	// 如果外部数据库可用，则使用外部数据库
	if m.ExtDBHost != "" && m.ExtDBUser != "" && m.ExtDBPassword != "" {
		m.Log.Infof("detected external db, will use external db as cluster global database instance")
		args := []string{"platform", "crd", "dbsvc", "new", "--host", m.ExtDBHost, "--port", strconv.Itoa(m.ExtDBPort), "--username", m.ExtDBUser, "--password", m.ExtDBPassword}
		if err := qcexec.CommandRun(os.Args[0], args...); err != nil {
			m.Log.Warnf("create external dbservice failed, reason: %v", err)
		} else {
			m.Log.Done("configure external db as cluster global database instance success")
			useExtDB = true
		}
	}

	helmargs := []string{"experimental", "helm", "upgrade", "--name", common.GetReleaseName(m.DevopsMode), "--repo", common.DefaultHelmRepoName, "--chart", common.GetReleaseName(m.DevopsMode), "--namespace", common.GetDefaultSystemNamespace(true), "--set", "env.APP_DOMAIN=" + m.Domain, "--set", "env.CNE_API_TOKEN=" + token, "--set", "cloud.defaultChannel=" + repoChannel}
	hostdomain := m.Domain
	if kutil.IsLegalDomain(hostdomain) && m.DomainType == "api" {
		m.Log.Debugf("use tls cert for domain %s", hostdomain)
		// helmargs = append(helmargs, "--set", "ingress.tls.enabled=true")
		// helmargs = append(helmargs, "--set", "ingress.tls.secretName=tls-haogs-cn")
	} else {
		if !m.DevopsMode {
			hostdomain = fmt.Sprintf("console.%s", hostdomain)
		}
	}
	if len(scName) > 0 {
		helmargs = append(helmargs, "--set", fmt.Sprintf("global.storageClass=%s", scName))
	}

	if m.OffLine {
		helmargs = append(helmargs, "--set", "cloud.host=http://market-cne-market-api.quickon-system.svc:8088")
		helmargs = append(helmargs, "--set", "env.CNE_MARKET_API_SCHEMA=http")
		helmargs = append(helmargs, "--set", "env.CNE_MARKET_API_HOST=market-cne-market-api.quickon-system.svc")
		helmargs = append(helmargs, "--set", "env.CNE_MARKET_API_PORT=8088")
	}
	if m.UsePHP7 {
		helmargs = append(helmargs, "--set", "deploy.phpVersion=php7")
	}
	if useExtDB {
		// create external db
		helmargs = append(helmargs, "--set", "mysql.enabled=false")
		helmargs = append(helmargs, "--set", fmt.Sprintf("mysql.auth.dbservice.name=%s", common.DefaultExternalDBName))
		helmargs = append(helmargs, "--set", fmt.Sprintf("mysql.auth.host=%s.quickon-system.svc", common.DefaultExternalDBName))
	}

	// TODO: 等下个版本禅道正式发版后续删除
	helmargs = append(helmargs, "--set", "env.ZT_SKIP_DEVOPS_INIT=true")

	if m.DBReplication {
		helmargs = append(helmargs, "--set", "mysql.replication.enabled=true")
		helmargs = append(helmargs, "--set", "env.ENABLE_DB_SLAVE=true")
	}

	helmargs = append(helmargs, "--set", fmt.Sprintf("ingress.host=%s", hostdomain))
	if m.DevopsMode {
		// 指定类型
		helmargs = append(helmargs, "--set", fmt.Sprintf("deploy.product=%s", m.Type))
		if repoChannel == "test" {
			helmargs = append(helmargs, "--set", "image.repository=test/zentao")
			helmargs = append(helmargs, "--set", "sidecars.backend.image.tag=dev")
			helmargs = append(helmargs, "--set", "sidecars.backend.image.pullPolicy=Always")
			helmargs = append(helmargs, "--set", "env.PHP_DEBUG=2")
			helmargs = append(helmargs, "--set", "cloud.switchChannel=true")
			helmargs = append(helmargs, "--set", "cloud.selectVersion=true")
		}
		if common.NeedUseCustomVersion(m.Version) {
			deployVersion := fmt.Sprintf("deploy.versions.%s=%s%s.k8s", m.Type, m.Type, installVersion)
			if m.Type == common.ZenTaoOSSType.String() || m.Type == common.ZenTaoOldOSSType.String() {
				deployVersion = fmt.Sprintf("deploy.versions.open=%s", installVersion)
			}
			helmargs = append(helmargs, "--set", deployVersion)
		}
	} else {
		if len(installVersion) > 0 {
			helmargs = append(helmargs, "--version", installVersion)
		}
	}

	output, err := qcexec.Command(os.Args[0], helmargs...).CombinedOutput()
	if err != nil {
		m.Log.Errorf("upgrade install web failed: %s", string(output))
		return err
	}
	m.Log.Donef("install %s success", common.GetReleaseName(m.DevopsMode))
	if m.OffLine {
		output, err := qcexec.Command(os.Args[0], "platform", "market-init", "--host", exnet.LocalIPs()[0]).CombinedOutput()
		if err != nil {
			m.Log.Warnf("upgrade install cloudapp market failed: %s", string(output))
		}
	}
	m.QuickONReady()
	initFile := common.GetCustomConfig(common.InitFileName)
	if err := file.WriteFile(initFile, "init done", true); err != nil {
		m.Log.Warnf("write init done file failed, reason: %v.\n\t please run: touch %s", err, initFile)
	}
	return nil
}

// QuickONReady 渠成Ready
func (m *Meta) QuickONReady() {
	clusterWaitGroup, ctx := errgroup.WithContext(context.Background())
	clusterWaitGroup.Go(func() error {
		return m.readyQuickON(ctx)
	})
	if err := clusterWaitGroup.Wait(); err != nil {
		m.Log.Error(err)
	}
}

// OperatorReady OperatorReady
func (m *Meta) OperatorReady() error {
	m.Log.Info("waiting for operator ready")
	for i := 1; i <= 60; i++ {
		deploy, err := m.kubeClient.GetDeployment(context.Background(), common.GetDefaultSystemNamespace(true), common.DefaultCneOperatorName, metav1.GetOptions{})
		if err != nil {
			time.Sleep(4 * time.Second)
			continue
		}
		ready := deploy.Status.ReadyReplicas == *deploy.Spec.Replicas
		if ready {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return errors.Errorf("operator not ready")
}

func (m *Meta) readyQuickON(ctx context.Context) error {
	t1 := ztime.NowUnix()
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG()).SetTimeout(time.Second * 1)
	m.Log.StartWait("waiting for ready")
	status := false
	for {
		t2 := ztime.NowUnix() - t1
		if t2 > 180 {
			m.Log.Warnf("waiting for ready 3min timeout: check your network or storage. after install you can run: %s status", os.Args[0])
			break
		}
		_, err := client.R().Get(fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0]))
		if err == nil {
			status = true
			break
		}
		time.Sleep(time.Second * 5)
	}
	m.Log.StopWait()
	if status {
		m.Log.Donef("%s ready, cost: %v", common.GetInstallType(m.DevopsMode), time.Since(time.Unix(t1, 0)))
	}
	return nil
}

func (m *Meta) getOrCreateUUIDAndAuth() (auth string, err error) {
	ns, err := m.kubeClient.GetNamespace(context.TODO(), common.DefaultKubeSystem, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(ns.GetUID()), nil
}

func (m *Meta) genSuffixHTTPHost(ip string) (domain, tls string, err error) {
	auth, err := m.getOrCreateUUIDAndAuth()
	if err != nil {
		return "", "", err
	}
	defaultDomain := suffixdomain.SearchCustomDomain(ip, auth, "")
	domain, tls, err = suffixdomain.GenerateDomain(ip, auth, suffixdomain.GenCustomDomain(defaultDomain))
	if err != nil {
		return "", "", err
	}
	return domain, tls, nil
}

func (m *Meta) UnInstall() error {
	m.Log.Warnf("start clean platform")
	cfg, _ := config.LoadConfig()
	// 清理helm安装应用
	// m.Log.Info("start uninstall cne custom tools")
	// toolArgs := []string{"experimental", "helm", "uninstall", "--name", "selfcert", "--namespace", common.GetDefaultSystemNamespace(true)}
	// if cleanStd, err := qcexec.Command(os.Args[0], toolArgs...).CombinedOutput(); err != nil {
	// 	m.Log.Warnf("uninstall cne custom tools err: %v, std: %s", err, string(cleanStd))
	// } else {
	// 	m.Log.Done("uninstall cne custom tools success")
	// }
	m.Log.Info("start uninstall cne operator")
	operatorArgs := []string{"experimental", "helm", "uninstall", "--name", common.DefaultCneOperatorName, "--namespace", common.GetDefaultSystemNamespace(true)}
	if cleanStd, err := qcexec.Command(os.Args[0], operatorArgs...).CombinedOutput(); err != nil {
		m.Log.Warnf("uninstall cne-operator err: %v, std: %s", err, string(cleanStd))
	} else {
		m.Log.Done("uninstall cne-operator success")
	}
	m.Log.Info("start uninstall platform")
	quickonCleanArgs := []string{"experimental", "helm", "uninstall", "--name", common.DefaultQuchengName, "--namespace", common.GetDefaultSystemNamespace(true)}
	if cleanStd, err := qcexec.Command(os.Args[0], quickonCleanArgs...).CombinedOutput(); err != nil {
		m.Log.Warnf("uninstall platform err: %v, std: %s", err, string(cleanStd))
	} else {
		m.Log.Done("uninstall platform success")
	}
	m.Log.Info("start uninstall helm repo")
	repoCleanArgs := []string{"experimental", "helm", "repo-del"}
	_ = qcexec.Command(os.Args[0], repoCleanArgs...).Run()
	m.Log.Done("uninstall helm repo success")
	if kutil.IsLegalDomain(cfg.Domain) {
		m.Log.Infof("clean domain %s", cfg.Domain)
		if err := qcexec.Command(os.Args[0], "exp", "tools", "domain", "clean", cfg.Domain).Run(); err != nil {
			m.Log.Warnf("clean domain %s failed, reason: %v", cfg.Domain, err)
		}
	}
	f := common.GetCustomConfig(common.InitFileName)
	if file.CheckFileExists(f) {
		os.Remove(f)
	}
	return nil
}
