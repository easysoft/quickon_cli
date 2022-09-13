// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
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
	var useip bool
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "install",
		Short:   "install app",
		Example: `q app install -n zentao-open or q app install zentao-open`,
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
			if useip || apiHost == "" {
				apiHost = fmt.Sprintf("%s:32379", exnet.LocalIPs()[0])
			} else if !strings.HasSuffix(apiHost, "haogs.cn") && !strings.HasSuffix(apiHost, "corp.cc") {
				apiHost = fmt.Sprintf("console.%s", cfg.Domain)
			}
			log.Debugf("install app %s, domain: %s.%s", name, domain, cfg.Domain)
			client := req.C().SetLogger(nil)
			if log.GetLevel() == logrus.DebugLevel {
				client = client.DevMode().EnableDumpAll()
			}
			var result Result
			resp, err := client.R().
				SetHeader("accept", "application/json").
				SetHeader("TOKEN", cfg.APIToken).
				SetBody(&Body{Chart: name, Domain: domain}).
				SetResult(&result).
				Post(fmt.Sprintf("http://%s/instance-apiInstall.html", apiHost))
			if err != nil {
				log.Errorf("install app %s failed, reason: %v", name, err)
				return fmt.Errorf("install app %s failed, reason: %v", name, err)
			}
			if !resp.IsSuccess() {
				log.Errorf("install app %s failed, reason: bad response status %v", name, resp.Status)
				return fmt.Errorf("install app %s failed, reason: bad response status %v", name, resp.Status)
			}
			if result.Code != 200 {
				log.Errorf("install app %s failed, reason: %s", name, result.Message)
				return fmt.Errorf("install app %s failed, reason: %s", name, result.Message)
			}
			log.Donef("app %s install success.", name)
			log.Infof("please wait, the app is starting.")
			hc, err := helm.NewClient(&helm.Config{Namespace: "default"})
			if err != nil {
				log.Debugf("create helm err: %v", err)
				return nil
			}
			release, _, err := hc.List(0, 0, fmt.Sprintf("%s-qadmin-%s", name, time.Now().Format("20060102")))
			if err != nil {
				log.Debugf("helm list %s err: %v", name, err)
				return nil
			}
			releaseValue, err := hc.GetAllValues(release[0].Name)
			if err != nil {
				log.Debugf("helm get allvalues %s err: %v", name, err)
				return nil
			}
			host := getMapValue(getMap(getMap(releaseValue, "global"), "ingress"), "host")
			if len(host) != 0 {
				if strings.HasSuffix(host, "haogs.cn") || strings.HasSuffix(host, "corp.cc") {
					host = fmt.Sprintf("https://%s", host)
				} else {
					host = fmt.Sprintf("http://%s", host)
				}
			}
			// spew.Dump(releaseValue)
			auth := getMap(releaseValue, "auth")
			// spew.Dump(auth)
			if auth != nil {
				authUsername := getMapValue(auth, "username")
				authPassword := getMapValue(auth, "password")
				log.Debugf("authUsername: %s, authPassword: %s", authUsername, authPassword)
				log.Infof("app meta:\n\t   username: %s\n\t   password: %s\n\t   url: %s", color.SBlue(authUsername), color.SBlue(authPassword), color.SBlue(host))
			} else {
				log.Infof("app meta:\n\t   url: %s", color.SBlue(host))
			}
			return nil
		},
	}
	app.Flags().StringVarP(&name, "name", "n", "zentao-open", "app name")
	app.Flags().StringVarP(&domain, "domain", "d", "", "app subdomain")
	app.Flags().BoolVar(&useip, "api-useip", false, "api use ip")
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
