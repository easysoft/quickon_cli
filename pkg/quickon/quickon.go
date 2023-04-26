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
	"strings"
	"time"

	"github.com/easysoft/qcadmin/internal/pkg/types"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/pkg/util/retry"
	suffixdomain "github.com/easysoft/qcadmin/pkg/qucheng/domain"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/expass"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/ztime"
	"github.com/imroc/req/v3"
	"golang.org/x/sync/errgroup"
	kubeerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Meta struct {
	Domain          string
	IP              string
	Version         string
	ConsolePassword string
	QuickonOSS      bool
	QuickonType     common.QuickonType
	kubeClient      *k8s.Client
	log             log.Logger
}

func New(f factory.Factory) *Meta {
	return &Meta{
		log: f.GetLog(),
		// Version:         common.DefaultQuickonOssVersion,
		ConsolePassword: expass.PwGenAlphaNum(32),
		QuickonType:     common.QuickonOSSType,
	}
}

func (m *Meta) GetFlags() []types.Flag {
	return []types.Flag{
		{
			Name:  "domain",
			Usage: "quickon domain",
			P:     &m.Domain,
			V:     m.Domain,
		},
		{
			Name:  "ip",
			Usage: "quickon ip",
			P:     &m.IP,
			V:     m.IP,
		},
		{
			Name:  "version",
			Usage: fmt.Sprintf("quickon version(oss: %s/ee: %s)", common.DefaultQuickonOssVersion, common.DefaultQuickonEEVersion),
			P:     &m.Version,
			V:     m.Version,
		},
		{
			Name:  "password",
			Usage: "quickon console password",
			P:     &m.ConsolePassword,
			V:     m.ConsolePassword,
		},
		{
			Name:  "oss",
			Usage: "quickon mode",
			P:     &m.QuickonOSS,
			V:     m.QuickonType == common.QuickonOSSType,
		},
	}
}

func (m *Meta) GetKubeClient() error {
	kubeClient, err := k8s.NewSimpleClient(common.GetKubeConfig())
	if err != nil {
		return errors.Errorf("load k8s client failed, reason: %v", err)
	}
	m.kubeClient = kubeClient
	return nil
}

func (m *Meta) checkIngress() {
	m.log.StartWait("check default ingress class")
	defaultClass, _ := m.kubeClient.ListDefaultIngressClass(context.Background(), metav1.ListOptions{})
	m.log.StopWait()
	if defaultClass == nil {
		m.log.Infof("not found default ingress class, will install nginx ingress")
		m.log.Debug("start install default ingress: nginx-ingress-controller")
		if err := qcexec.CommandRun(os.Args[0], "quickon", "plugins", "enable", "ingress"); err != nil {
			m.log.Errorf("install ingress failed, reason: %v", err)
		} else {
			m.log.Done("install ingress: cne-ingress success")
		}
	} else {
		m.log.Infof("found exist default ingress class: %s", defaultClass.Name)
	}
	m.log.Done("check default ingress done")
}

func (m *Meta) checkStorage() {
	m.log.StartWait("check default storage class")
	defaultClass, _ := m.kubeClient.GetDefaultSC(context.Background())
	m.log.StopWait()
	if defaultClass == nil {
		m.log.Infof("not found default storage class, will install default storage")
		m.log.Debug("start install default storage: local-storage")
		if err := qcexec.CommandRun(os.Args[0], "quickon", "plugins", "enable", "storage"); err != nil {
			m.log.Errorf("install storage failed, reason: %v", err)
		} else {
			m.log.Done("install storage: local-storage success")
		}
	} else {
		m.log.Infof("found exist default storage class: %s", defaultClass.Name)
	}
	m.log.Done("check default storage done")
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
	m.log.Debugf("init quickon default namespace.")
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
	m.log.Donef("init quickon default namespace success.")
	return nil
}

func (m *Meta) addHelmRepo() error {
	output, err := qcexec.Command(os.Args[0], "experimental", "helm", "repo-add", "--name", common.DefaultHelmRepoName, "--url", common.GetChartRepo(m.Version)).CombinedOutput()
	if err != nil {
		errmsg := string(output)
		if !strings.Contains(errmsg, "exists") {
			m.log.Errorf("init quickon helm repo failed, reason: %s", string(output))
			return err
		}
		m.log.Debugf("quickon helm repo already exists")
	} else {
		m.log.Done("add quickon helm repo success")
	}
	output, err = qcexec.Command(os.Args[0], "experimental", "helm", "repo-update").CombinedOutput()
	if err != nil {
		m.log.Errorf("update quickon helm repo failed, reason: %s", string(output))
		return err
	}
	m.log.Done("update quickon helm repo success")
	return nil
}

func (m *Meta) Init() error {
	m.log.Info("executing init quickon logic...")
	ctx := context.Background()
	m.log.Debug("waiting for storage to be ready...")
	waitsc := time.Now()
	// wait.BackoffUntil TODO
	for {
		sc, _ := m.kubeClient.GetDefaultSC(ctx)
		if sc != nil {
			m.log.Donef("default storage %s is ready", sc.Name)
			break
		}
		time.Sleep(time.Second * 5)
		trywaitsc := time.Now()
		if trywaitsc.Sub(waitsc) > time.Minute*3 {
			m.log.Warnf("wait storage ready, timeout: %v", trywaitsc.Sub(waitsc).Seconds())
			break
		}
	}

	_, err := m.kubeClient.CreateNamespace(ctx, common.GetDefaultSystemNamespace(true), metav1.CreateOptions{})
	if err != nil {
		if !kubeerr.IsAlreadyExists(err) {
			return err
		}
	}
	chartVersion := common.GetVersion(m.Version, m.QuickonType)
	m.log.Debugf("start init quickon %v, version: %s", m.QuickonType, chartVersion)
	if m.Domain == "" {
		err := retry.Retry(time.Second*1, 3, func() (bool, error) {
			domain, _, err := m.genSuffixHTTPHost(m.IP)
			if err != nil {
				return false, err
			}
			m.Domain = domain
			m.log.Infof("generate suffix domain: %s, ip: %v", color.SGreen(m.Domain), color.SGreen(m.IP))
			return true, nil
		})
		if err != nil {
			m.Domain = "demo.corp.cc"
			m.log.Warnf("gen suffix domain failed, reason: %v, use default domain: %s", err, m.Domain)
		}
		m.log.Infof("load %s tls cert", m.Domain)
		defaultTLS := fmt.Sprintf("%s/tls-haogs-cn.yaml", common.GetDefaultCacheDir())
		m.log.StartWait(fmt.Sprintf("start issuing domain %s certificate, may take 3-5min", m.Domain))
		waittls := time.Now()
		for {
			if file.CheckFileExists(defaultTLS) {
				m.log.StopWait()
				m.log.Done("download tls cert success")
				if err := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", defaultTLS, "-n", common.GetDefaultSystemNamespace(true), "--kubeconfig", common.GetKubeConfig()).Run(); err != nil {
					m.log.Warnf("load default tls cert failed, reason: %v", err)
				} else {
					m.log.Done("load default tls cert success")
				}
				qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", defaultTLS, "-n", "default", "--kubeconfig", common.GetKubeConfig()).Run()
				break
			}
			_, mainDomain := kutil.SplitDomain(m.Domain)
			domainTLS := fmt.Sprintf("https://pkg.qucheng.com/ssl/%s/%s/tls.yaml", mainDomain, m.Domain)
			qcexec.Command(os.Args[0], "experimental", "tools", "wget", "-t", domainTLS, "-d", defaultTLS).Run()
			m.log.Debug("wait for tls cert ready...")
			time.Sleep(time.Second * 5)
			trywaitsc := time.Now()
			if trywaitsc.Sub(waittls) > time.Minute*3 {
				// TODO  timeout
				m.log.Debugf("wait tls cert ready, timeout: %v", trywaitsc.Sub(waittls).Seconds())
				break
			}
		}
	} else {
		m.log.Infof("use custom domain %s, you should add dns record to your domain: *.%s -> %s", m.Domain, color.SGreen(m.Domain), color.SGreen(m.IP))
	}
	token := expass.PwGenAlphaNum(32)
	cfg, _ := config.LoadConfig()
	cfg.Domain = m.Domain
	cfg.APIToken = token
	cfg.S3.Username = expass.PwGenAlphaNum(8)
	cfg.S3.Password = expass.PwGenAlphaNum(16)
	cfg.Quickon.Type = m.QuickonType
	cfg.SaveConfig()
	m.log.Info("start deploy cne custom tools")
	toolargs := []string{"experimental", "helm", "upgrade", "--name", "selfcert", "--repo", common.DefaultHelmRepoName, "--chart", "selfcert", "--namespace", common.GetDefaultSystemNamespace(true)}
	if helmstd, err := qcexec.Command(os.Args[0], toolargs...).CombinedOutput(); err != nil {
		m.log.Warnf("deploy cne custom tools err: %v, std: %s", err, string(helmstd))
	} else {
		m.log.Done("deployed cne custom tools success")
	}
	m.log.Info("start deploy cne operator")
	operatorargs := []string{"experimental", "helm", "upgrade", "--name", common.DefaultCneOperatorName, "--repo", common.DefaultHelmRepoName, "--chart", common.DefaultCneOperatorName, "--namespace", common.GetDefaultSystemNamespace(true),
		"--set", "minio.ingress.enabled=true",
		"--set", "minio.ingress.host=s3." + m.Domain,
		"--set", "minio.auth.username=" + cfg.S3.Username,
		"--set", "minio.auth.password=" + cfg.S3.Password,
	}
	//if len(chartversion) > 0 {
	//	operatorargs = append(operatorargs, "--version", chartversion)
	//}
	if helmstd, err := qcexec.Command(os.Args[0], operatorargs...).CombinedOutput(); err != nil {
		m.log.Warnf("deploy cne-operator err: %v, std: %s", err, string(helmstd))
	} else {
		m.log.Done("deployed cne-operator success")
	}
	helmchan := common.GetChannel(m.Version)
	helmargs := []string{"experimental", "helm", "upgrade", "--name", common.DefaultQuchengName, "--repo", common.DefaultHelmRepoName, "--chart", common.GetQuickONName(m.QuickonType), "--namespace", common.GetDefaultSystemNamespace(true), "--set", "env.APP_DOMAIN=" + m.Domain, "--set", "env.CNE_API_TOKEN=" + token, "--set", "cloud.defaultChannel=" + helmchan}
	if helmchan != "stable" {
		helmargs = append(helmargs, "--set", "env.PHP_DEBUG=2")
		helmargs = append(helmargs, "--set", "cloud.switchChannel=true")
		helmargs = append(helmargs, "--set", "cloud.selectVersion=true")
	}
	hostdomain := m.Domain
	if kutil.IsLegalDomain(hostdomain) {
		helmargs = append(helmargs, "--set", "ingress.tls.enabled=true")
		helmargs = append(helmargs, "--set", "ingress.tls.secretName=tls-haogs-cn")
	} else {
		hostdomain = fmt.Sprintf("console.%s", hostdomain)
	}
	helmargs = append(helmargs, "--set", fmt.Sprintf("ingress.host=%s", hostdomain))

	if len(chartVersion) > 0 {
		helmargs = append(helmargs, "--version", chartVersion)
	}
	output, err := qcexec.Command(os.Args[0], helmargs...).CombinedOutput()
	if err != nil {
		m.log.Errorf("upgrade install quickon web failed: %s", string(output))
		return err
	}
	m.log.Done("install quickon success")
	m.QuickONReady()
	initFile := common.GetCustomConfig(common.InitFileName)
	if err := file.WriteFile(initFile, "init done", true); err != nil {
		m.log.Warnf("write init done file failed, reason: %v.\n\t please run: touch %s", err, initFile)
	}
	m.Show()
	return nil
}

// QuickONReady 渠成Ready
func (m *Meta) QuickONReady() {
	clusterWaitGroup, ctx := errgroup.WithContext(context.Background())
	clusterWaitGroup.Go(func() error {
		return m.readyQuickON(ctx)
	})
	if err := clusterWaitGroup.Wait(); err != nil {
		m.log.Error(err)
	}
}

func (m *Meta) readyQuickON(ctx context.Context) error {
	t1 := ztime.NowUnix()
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG()).SetTimeout(time.Second * 1)
	m.log.StartWait("waiting for qucheng ready")
	status := false
	for {
		t2 := ztime.NowUnix() - t1
		if t2 > 180 {
			m.log.Warnf("waiting for qucheng ready 3min timeout: check your network or storage. after install you can run: q status")
			break
		}
		_, err := client.R().Get(fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0]))
		if err == nil {
			status = true
			break
		}
		time.Sleep(time.Second * 10)
	}
	m.log.StopWait()
	if status {
		m.log.Donef("qucheng ready, cost: %v", time.Since(time.Unix(t1, 0)))
	}
	return nil
}

func (m *Meta) getOrCreateUUIDAndAuth() (auth string, err error) {
	// cm := &corev1.ConfigMap{}
	cm, err := m.kubeClient.Clientset.CoreV1().ConfigMaps(common.GetDefaultSystemNamespace(true)).Get(context.TODO(), "q-suffix-host", metav1.GetOptions{})
	if err != nil {
		if !kubeerr.IsNotFound(err) {
			return "", err
		}
		if kubeerr.IsNotFound(err) {
			m.log.Debug("q-suffix-host not found, create it")
			cm = suffixdomain.GenerateSuffixConfigMap("q-suffix-host", common.GetDefaultSystemNamespace(true))
			if _, err := m.kubeClient.Clientset.CoreV1().ConfigMaps(common.GetDefaultSystemNamespace(true)).Create(context.TODO(), cm, metav1.CreateOptions{}); err != nil {
				return "", err
			}
		}
	}
	return cm.Data["auth"], nil
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

func (m *Meta) Show() {
	if len(m.IP) <= 0 {
		m.IP = exnet.LocalIPs()[0]
	}
	resetPassArgs := []string{"quickon", "reset-password", "--password", m.ConsolePassword}
	qcexec.CommandRun(os.Args[0], resetPassArgs...)
	cfg, _ := config.LoadConfig()
	cfg.ConsolePassword = m.ConsolePassword
	cfg.SaveConfig()
	domain := cfg.Domain

	m.log.Info("----------------------------\t")
	if len(domain) > 0 {
		if !kutil.IsLegalDomain(cfg.Domain) {
			domain = fmt.Sprintf("http://console.%s", cfg.Domain)
		} else {
			domain = fmt.Sprintf("https://%s", cfg.Domain)
		}
	} else {
		domain = fmt.Sprintf("http://%s:32379", m.IP)
	}
	m.log.Donef("console: %s, username: %s, password: %s",
		color.SGreen(domain), color.SGreen(common.QuchengDefaultUser), color.SGreen(m.ConsolePassword))
	m.log.Donef("docs: %s", common.QuchengDocs)
	m.log.Done("support: 768721743(QQGroup)")
}

func (m *Meta) UnInstall() error {
	m.log.Warnf("start clean quickon.")
	// 清理helm安装应用
	m.log.Info("start uninstall cne custom tools")
	toolArgs := []string{"experimental", "helm", "uninstall", "--name", "selfcert", "--namespace", common.GetDefaultSystemNamespace(true)}
	if cleanStd, err := qcexec.Command(os.Args[0], toolArgs...).CombinedOutput(); err != nil {
		m.log.Warnf("uninstall cne custom tools err: %v, std: %s", err, string(cleanStd))
	} else {
		m.log.Done("uninstall cne custom tools success")
	}
	m.log.Info("start uninstall cne operator")
	operatorArgs := []string{"experimental", "helm", "uninstall", "--name", common.DefaultCneOperatorName, "--namespace", common.GetDefaultSystemNamespace(true)}
	if cleanStd, err := qcexec.Command(os.Args[0], operatorArgs...).CombinedOutput(); err != nil {
		m.log.Warnf("uninstall cne-operator err: %v, std: %s", err, string(cleanStd))
	} else {
		m.log.Done("uninstall cne-operator success")
	}
	m.log.Info("start uninstall cne quickon")
	quickonCleanArgs := []string{"experimental", "helm", "uninstall", "--name", common.DefaultQuchengName, "--namespace", common.GetDefaultSystemNamespace(true)}
	if cleanStd, err := qcexec.Command(os.Args[0], quickonCleanArgs...).CombinedOutput(); err != nil {
		m.log.Warnf("uninstall quickon err: %v, std: %s", err, string(cleanStd))
	} else {
		m.log.Done("uninstall quickon success")
	}
	m.log.Info("start uninstall helm repo")
	repoCleanArgs := []string{"experimental", "helm", "repo-del"}
	_ = qcexec.Command(os.Args[0], repoCleanArgs...).Run()
	m.log.Done("uninstall helm repo success")
	f := common.GetCustomConfig(common.InitFileName)
	if file.CheckFileExists(f) {
		os.Remove(f)
	}
	return nil
}
