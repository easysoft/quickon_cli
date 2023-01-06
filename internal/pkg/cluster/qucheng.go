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
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/retry"
	suffixdomain "github.com/easysoft/qcadmin/pkg/qucheng/domain"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/expass"
	"github.com/ergoapi/util/file"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Cluster) genQuChengToken() string {
	// TODO token 生成优化
	return expass.RandomPassword(32)
}

func (p *Cluster) getOrCreateUUIDAndAuth() (auth string, err error) {
	// cm := &corev1.ConfigMap{}
	cm, err := p.KubeClient.Clientset.CoreV1().ConfigMaps(common.DefaultSystem).Get(context.TODO(), "q-suffix-host", metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return "", err
		}
		if errors.IsNotFound(err) {
			p.Log.Debug("q-suffix-host not found, create it")
			cm = suffixdomain.GenerateSuffixConfigMap("q-suffix-host", common.DefaultSystem)
			if _, err := p.KubeClient.Clientset.CoreV1().ConfigMaps(common.DefaultSystem).Create(context.TODO(), cm, metav1.CreateOptions{}); err != nil {
				return "", err
			}
		}
	}
	p.Status.UUID = cm.Data["auth"]
	return p.Status.UUID, nil
}

func (p *Cluster) genSuffixHTTPHost(ip string) (domain, tls string, err error) {
	auth, err := p.getOrCreateUUIDAndAuth()
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

func (p *Cluster) InstallQuCheng() error {
	p.Log.Info("executing init qucheng logic...")
	ctx := context.Background()
	p.Log.Debug("waiting for storage to be ready...")
	waitsc := time.Now()
	// wait.BackoffUntil TODO
	for {
		sc, _ := p.KubeClient.GetDefaultSC(ctx)
		if sc != nil {
			p.Log.Donef("default storage %s is ready", sc.Name)
			break
		}
		time.Sleep(time.Second * 5)
		trywaitsc := time.Now()
		if trywaitsc.Sub(waitsc) > time.Minute*3 {
			p.Log.Warnf("wait storage ready, timeout: %v", trywaitsc.Sub(waitsc).Seconds())
			break
		}
	}

	_, err := p.KubeClient.CreateNamespace(ctx, common.DefaultSystem, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}
	p.Log.Debug("start init qucheng")

	if len(p.Metadata.EIP) <= 0 {
		p.Metadata.EIP = exnet.LocalIPs()[0]
	}
	if p.Domain == "" {
		err := retry.Retry(time.Second*1, 3, func() (bool, error) {
			domain, _, err := p.genSuffixHTTPHost(p.Metadata.EIP)
			if err != nil {
				return false, err
			}
			p.Domain = domain

			p.Log.Infof("generate suffix domain: %s, ip: %v", color.SGreen(p.Domain), color.SGreen(p.Metadata.EIP))
			return true, nil
		})
		if err != nil {
			p.Domain = "demo.corp.cc"
			p.Log.Warnf("gen suffix domain failed, reason: %v, use default domain: %s", err, p.Domain)
		}
		p.Log.Infof("load %s tls cert", p.Domain)
		defaultTLS := fmt.Sprintf("%s/tls-haogs-cn.yaml", common.GetDefaultCacheDir())
		p.Log.StartWait(fmt.Sprintf("start issuing domain %s certificate, may take 3-5min", p.Domain))
		waittls := time.Now()
		for {
			if _, err := os.Stat(defaultTLS); err == nil {
				p.Log.StopWait()
				p.Log.Done("download tls cert success")
				if err := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", defaultTLS, "-n", common.DefaultSystem).Run(); err != nil {
					p.Log.Warnf("load default tls cert failed, reason: %v", err)
				} else {
					p.Log.Done("load default tls cert success")
				}
				qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", defaultTLS, "-n", "default").Run()
				break
			}
			_, mainDomain := kutil.SplitDomain(p.Domain)
			domainTLS := fmt.Sprintf("https://pkg.qucheng.com/ssl/%s/%s/tls.yaml", mainDomain, p.Domain)
			qcexec.Command(os.Args[0], "experimental", "tools", "wget", "-t", domainTLS, "-d", defaultTLS).Run()
			p.Log.Debug("wait for tls cert ready...")
			time.Sleep(time.Second * 5)
			trywaitsc := time.Now()
			if trywaitsc.Sub(waittls) > time.Minute*3 {
				// TODO  timeout
				p.Log.Debugf("wait tls cert ready, timeout: %v", trywaitsc.Sub(waittls).Seconds())
				break
			}
		}
	} else {
		p.Log.Infof("use custom domain %s, you should add dns record to your domain: *.%s -> %s", p.Domain, color.SGreen(p.Domain), color.SGreen(p.Metadata.EIP))
	}
	token := p.genQuChengToken()
	cfg, _ := config.LoadConfig()
	cfg.Domain = p.Domain
	cfg.APIToken = token
	cfg.ClusterID = p.Status.UUID
	cfg.S3.Username = expass.PwGenAlphaNum(8)
	cfg.S3.Password = expass.PwGenAlphaNum(16)
	cfg.DataDir = p.DataDir
	cfg.SaveConfig()
	chartversion := common.GetVersion(p.QuchengVersion)
	p.Log.Info("start deploy cne custom tools")
	toolargs := []string{"experimental", "helm", "upgrade", "--name", "selfcert", "--repo", common.DefaultHelmRepoName, "--chart", "selfcert", "--namespace", common.DefaultSystem}
	if helmstd, err := qcexec.Command(os.Args[0], toolargs...).CombinedOutput(); err != nil {
		p.Log.Warnf("deploy cne custom tools err: %v, std: %s", err, string(helmstd))
	} else {
		p.Log.Done("deployed cne custom tools success")
	}
	p.Log.Info("start deploy cne operator")
	operatorargs := []string{"experimental", "helm", "upgrade", "--name", common.DefaultCneOperatorName, "--repo", common.DefaultHelmRepoName, "--chart", common.DefaultCneOperatorName, "--namespace", common.DefaultSystem,
		"--set", "minio.ingress.enabled=true",
		"--set", "minio.ingress.host=s3." + p.Domain,
		"--set", "minio.auth.username=" + cfg.S3.Username,
		"--set", "minio.auth.password=" + cfg.S3.Password,
	}
	if len(chartversion) > 0 {
		operatorargs = append(operatorargs, "--version", chartversion)
	}
	if helmstd, err := qcexec.Command(os.Args[0], operatorargs...).CombinedOutput(); err != nil {
		p.Log.Warnf("deploy cne-operator err: %v, std: %s", err, string(helmstd))
	} else {
		p.Log.Done("deployed cne-operator success")
	}
	helmchan := common.GetChannel(p.QuchengVersion)
	helmargs := []string{"experimental", "helm", "upgrade", "--name", common.DefaultQuchengName, "--repo", common.DefaultHelmRepoName, "--chart", common.DefaultQuchengName, "--namespace", common.DefaultSystem, "--set", "env.APP_DOMAIN=" + p.Domain, "--set", "env.CNE_API_TOKEN=" + token, "--set", "cloud.defaultChannel=" + helmchan}
	if helmchan != "stable" {
		helmargs = append(helmargs, "--set", "env.PHP_DEBUG=2")
		helmargs = append(helmargs, "--set", "cloud.switchChannel=true")
		helmargs = append(helmargs, "--set", "cloud.selectVersion=true")
	}
	hostdomain := p.Domain
	if kutil.IsLegalDomain(hostdomain) {
		helmargs = append(helmargs, "--set", "ingress.tls.enabled=true")
		helmargs = append(helmargs, "--set", "ingress.tls.secretName=tls-haogs-cn")
	} else {
		hostdomain = fmt.Sprintf("console.%s", hostdomain)
	}
	helmargs = append(helmargs, "--set", fmt.Sprintf("ingress.host=%s", hostdomain))

	if len(chartversion) > 0 {
		helmargs = append(helmargs, "--version", chartversion)
	}
	output, err := qcexec.Command(os.Args[0], helmargs...).CombinedOutput()
	if err != nil {
		p.Log.Errorf("upgrade install qucheng web failed: %s", string(output))
		return err
	}
	p.Log.Done("install qucheng success")
	p.Ready()
	initfile := common.GetCustomConfig(common.InitFileName)
	if err := file.Writefile(initfile, "init done"); err != nil {
		p.Log.Warnf("write init done file failed, reason: %v.\n\t please run: touch %s", err, initfile)
	}
	if !p.DisableInstallApp {
		p.defaultAppInstall()
	}
	return nil
}
