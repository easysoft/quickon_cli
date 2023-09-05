// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Account string `json:"account"`
	} `json:"data"`
}

type Body struct {
	Chart  string `json:"chart"`
	Domain string `json:"domain"`
}

func NewCmdAppInstall(f factory.Factory) *cobra.Command {
	var name, domain string
	var useIP bool
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "install",
		Short:   "install app",
		Example: `q app install -n zentao or q app install zentao`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}
			if len(domain) == 0 {
				domain = name
			}
			cfg, _ := config.LoadConfig()
			apiHost := cfg.Domain
			if useIP || apiHost == "" {
				apiHost = fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0])
			} else if !kutil.IsLegalDomain(apiHost) {
				apiHost = fmt.Sprintf("http://console.%s", cfg.Domain)
			} else {
				apiHost = fmt.Sprintf("https://%s", apiHost)
			}
			log.Debugf("install app %s, domain: %s.%s", name, domain, cfg.Domain)
			client := req.C().SetLogger(nil).SetUserAgent(common.GetUG())
			if log.GetLevel() == logrus.DebugLevel {
				client = client.DevMode().EnableDumpAll()
			}
			var result Result
			resp, err := client.R().
				SetHeader("accept", "application/json").
				SetHeader("TOKEN", cfg.APIToken).
				SetBody(&Body{Chart: name, Domain: domain}).
				SetSuccessResult(&result).
				Post(fmt.Sprintf("%s/instance-apiInstall.html", apiHost))
			if err != nil {
				log.Errorf("install app %s failed, reason: %v", name, err)
				return errors.Errorf("install app %s failed, reason: %v", name, err)
			}
			if !resp.IsSuccessState() {
				log.Errorf("install app %s failed, reason: bad response status %v", name, resp.Status)
				return errors.Errorf("install app %s failed, reason: bad response status %v", name, resp.Status)
			}
			if result.Code != 200 {
				log.Errorf("install app %s failed, reason: %s", name, result.Message)
				return errors.Errorf("install app %s failed, reason: %s", name, result.Message)
			}
			log.Donef("app %s install success.", name)
			log.Infof("please wait, the app is starting.")
			hc, err := helm.NewClient(&helm.Config{Namespace: common.DefaultAppNamespace})
			if err != nil {
				log.Debugf("create helm err: %v", err)
				return nil
			}
			release, _, err := hc.List(0, 0, fmt.Sprintf("%s-qadmin-%s", name, time.Now().Format("20060102")))
			if err != nil {
				log.Debugf("helm list %s err: %v", name, err)
				return nil
			}
			if len(release) == 0 {
				// 2.6版本不支持
				log.Infof("please login console check it.")
				return nil
			}
			releaseValue, err := hc.GetAllValues(release[0].Name)
			if err != nil {
				log.Debugf("helm get all values %s err: %v", name, err)
				return nil
			}
			host := getMapValue(getMap(getMap(releaseValue, "global"), "ingress"), "host")
			if len(host) != 0 {
				if kutil.IsLegalDomain(host) {
					host = fmt.Sprintf("https://%s", host)
				} else {
					host = fmt.Sprintf("http://%s", host)
				}
			}
			auth := getMap(releaseValue, "auth")
			if auth != nil {
				authUsername := getMapValue(auth, "username")
				authPassword := getMapValue(auth, "password")
				log.Debugf("authUsername: %s, authPassword: %s", authUsername, authPassword)
				log.Infof("app meta:\n\t   username: %s\n\t   password: %s\n\t   url: %s", color.SBlue(authUsername), color.SBlue(authPassword), color.SBlue(host))
			} else {
				log.Infof("app meta:\n\turl: %s", color.SBlue(host))
			}
			return nil
		},
	}
	app.Flags().StringVarP(&name, "name", "n", "zentao", "app name")
	app.Flags().StringVarP(&domain, "domain", "d", "", "app subdomain")
	app.Flags().BoolVar(&useIP, "api-useip", true, "api use ip")
	return app
}

func getMap(mapValues map[string]interface{}, key string) map[string]interface{} {
	if key == "" || mapValues == nil {
		return nil
	}

	if v, ok := mapValues[key]; ok {
		return v.(map[string]interface{})
	}
	return nil
}

func getMapValue(mapValues map[string]interface{}, key string) string {
	if key == "" || mapValues == nil {
		return "-"
	}

	if v, ok := mapValues[key]; ok {
		return v.(string)
	}
	return "-"
}
