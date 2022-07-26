// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package tool

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/easysoft/qcadmin/pkg/qucheng/domain"
	suffixdomain "github.com/easysoft/qcadmin/pkg/qucheng/domain"
	"github.com/ergoapi/util/exmap"
	"github.com/ergoapi/util/exnet"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func EmbedDomainCommand(f factory.Factory) *cobra.Command {
	domain := &cobra.Command{
		Use:    "domain",
		Short:  "domain manager",
		Hidden: true,
	}
	domain.AddCommand(domainClean(f))
	domain.AddCommand(domainAdd(f))
	return domain
}

func domainClean(f factory.Factory) *cobra.Command {
	dns := &cobra.Command{
		Use:    "clean",
		Short:  "clean domain",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig()
			if cfg != nil {
				if !strings.HasSuffix(cfg.Domain, "haogs.cn") {
					return
				}
			}
			kclient, _ := k8s.NewSimpleClient()
			cm, err := kclient.Clientset.CoreV1().ConfigMaps(common.DefaultSystem).Get(context.TODO(), "q-suffix-host", metav1.GetOptions{})
			if err != nil {
				return
			}
			reqbody := domain.ReqBody{
				UUID:      cm.Data["uuid"],
				SecretKey: cm.Data["auth"],
			}
			client := req.C().SetUserAgent(common.GetUG())
			if _, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(&reqbody).
				Delete(common.GetAPI("/api/qdns/oss/record")); err != nil {
				f.GetLog().Error("clean dns failed, reason: %v", err)
			}
		},
	}
	return dns
}

func domainAdd(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var customdomain string
	dns := &cobra.Command{
		Use:    "init",
		Short:  "init domain",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			// load config
			domain := ""
			cfg, _ := config.LoadConfig()
			if cfg != nil {
				domain = cfg.Domain
			}
			if len(domain) > 0 {
				return
			}
			if len(customdomain) == 0 {
				kclient, _ := k8s.NewSimpleClient()
				cm, err := kclient.Clientset.CoreV1().ConfigMaps(common.DefaultSystem).Get(context.TODO(), "q-suffix-host", metav1.GetOptions{})
				if err != nil {
					if errors.IsNotFound(err) {
						log.Debug("q-suffix-host not found, create it")
						cm = suffixdomain.GenerateSuffixConfigMap("q-suffix-host", common.DefaultSystem)
						if _, err := kclient.Clientset.CoreV1().ConfigMaps(common.DefaultSystem).Create(context.TODO(), cm, metav1.CreateOptions{}); err != nil {
							log.Errorf("k8s api err: %v", err)
							return
						}
					} else {
						log.Errorf("conn k8s err: %v", err)
						return
					}
				}
				id := cm.Data["uuid"]
				auth := cm.Data["auth"]
				ip := exnet.LocalIPs()[0]
				domain, _, err = suffixdomain.GenerateDomain(ip, id, auth, suffixdomain.GenCustomDomain(suffixdomain.SearchCustomDomain(ip, id, auth)))
				if len(domain) == 0 {
					log.Warnf("gen domain failed: %v, use default domain: demo.haogs.cn", err)
					domain = "demo.haogs.cn"
				}
				cfg.Domain = domain
			} else {
				cfg.Domain = customdomain
			}
			// save config
			cfg.SaveConfig()
			// upgrade qucheng
			helmClient, _ := helm.NewClient(&helm.Config{Namespace: common.DefaultSystem})
			if err := helmClient.UpdateRepo(); err != nil {
				log.Warn("update repo failed, reason: %v", err)
			}
			if err := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", fmt.Sprintf("%s/hack/haogstls/haogs.yaml", common.GetDefaultDataDir()), "-n", common.DefaultSystem).Run(); err != nil {
				log.Warn("load default tls cert failed, reason: %v", err)
			} else {
				log.Done("load default tls cert success")
			}
			defaultValue, _ := helmClient.GetValues(common.DefaultQuchengName)
			var values []string
			host := cfg.Domain
			if strings.HasSuffix(host, "haogs.cn") {
				values = append(values, "ingress.tls.enabled=true")
				values = append(values, "ingress.tls.secretName=tls-haogs-cn")
			} else {
				host = fmt.Sprintf("console.%s", host)
			}
			values = append(values, fmt.Sprintf("ingress.host=%s", host), fmt.Sprintf("env.APP_DOMAIN=%s", cfg.Domain))
			base := map[string]interface{}{}
			for _, value := range values {
				strvals.ParseInto(value, base)
			}
			defaultValue = exmap.MergeMaps(defaultValue, base)
			if _, err := helmClient.Upgrade(common.DefaultQuchengName, common.DefaultHelmRepoName, common.DefaultQuchengName, "", defaultValue); err != nil {
				log.Warnf("upgrade %s failed, reason: %v", common.DefaultQuchengName, err)
			} else {
				log.Donef("upgrade %s success", common.DefaultQuchengName)
			}
		},
	}
	dns.Flags().StringVar(&customdomain, "domain", "", "app custom domain")
	return dns
}
