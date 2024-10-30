// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package httptls

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
)

// CheckReNewCertificate 检查证书是否过期
func CheckReNewCertificate(force bool) (err error) {
	log := log.GetInstance()
	cfg, _ := config.LoadConfig()
	domain := cfg.Domain
	if kutil.IsLegalDomain(domain) {
		needRenew := false
		if force {
			needRenew = true
		} else {
			needRenew, err = checkCertificate(fmt.Sprintf("https://%s", domain))
			if err != nil {
				log.Errorf("check domain %s tls err: %v", domain, err)
				return err
			}
		}
		if needRenew {
			log.Infof("domain %s's certificate need renew now, start to renew", domain)
			return renewCertificate(domain)
		}
		log.Infof("domain %s's certificate has not expired ", domain)
		return nil
	}
	log.Infof("custom domain %s not support", domain)
	return nil
}

func checkCertificate(domain string) (bool, error) {
	log := log.GetInstance()
	log.Infof("start check domain %s certificate", domain)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false}, // nolint:gosec
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}
	resp, err := client.Get(domain)
	if err != nil {
		if strings.Contains(err.Error(), "x509: certificate is valid for ingress.local") {
			log.Warnf("domain %s use self-signed certificate", domain)
			return true, nil
		}
		if strings.Contains(err.Error(), "x509: certificate has expired or is not vet valid") {
			log.Warnf("domain %s tls expired", domain)
			return true, nil
		}
		return false, err
	}
	defer func() { _ = resp.Body.Close() }()
	for _, cert := range resp.TLS.PeerCertificates {
		// 证书过期已过期
		if !cert.NotAfter.After(time.Now()) {
			log.Warnf("domain %s tls expired", domain)
			return true, nil
		}
		// 证书过期时间在7天内过期
		if cert.NotAfter.Sub(time.Now()).Hours() < 7*24 {
			log.Warnf("domain %s tls expire after %fh", domain, cert.NotAfter.Sub(time.Now()).Hours())
			return true, nil
		}
	}
	return false, nil
}

func renewCertificate(domain string) error {
	log := log.GetInstance()
	ds := strings.Split(domain, ".")
	mainDomain := fmt.Sprintf("%s.%s", ds[len(ds)-2], ds[len(ds)-1])
	coreDomain := fmt.Sprintf("%s.%s.%s", ds[len(ds)-3], ds[len(ds)-2], ds[len(ds)-1])
	tlsfile := fmt.Sprintf("https://pkg.qucheng.com/ssl/%s/%s/tls.yaml", mainDomain, coreDomain)
	log.Infof("renew %s tls certificate", domain)
	log.Debugf("renew default tls certificate use %s", tlsfile)
	if err := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", tlsfile, "-n", common.GetDefaultSystemNamespace(true), "--kubeconfig", common.GetKubeConfig()).Run(); err != nil {
		log.Warnf("load renew tls cert for %s failed, reason: %v", common.GetDefaultSystemNamespace(true), err)
		return nil
	}
	if err := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", tlsfile, "--kubeconfig", common.GetKubeConfig()).Run(); err != nil {
		log.Warnf("load renew tls cert for default failed, reason: %v", err)
		return nil
	}
	log.Donef("renew tls certificate success")
	return nil
}
