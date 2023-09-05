// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package debug

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/ergoapi/util/exnet"
	"github.com/imroc/req/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Result struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    AppData `json:"data"`
}

type AppData struct {
	ID         string `json:"id"`
	Space      string `json:"space"`
	Name       string `json:"name"`
	AppID      string `json:"appID"`
	AppName    string `json:"appName"`
	AppVersion string `json:"appVersion"`
	Chart      string `json:"chart"`
	Logo       string `json:"logo"`
	Version    string `json:"version"`
	Source     string `json:"source"`
	K8Name     string `json:"k8name"`
	Status     string `json:"status"`
	Domain     string `json:"domain"`
	CreatedBy  string `json:"createdBy"`
	CreatedAt  string `json:"createdAt"`
	Deleted    string `json:"deleted"`
}

func GetNameByURL(url string, debug, useip bool) (*AppData, error) {
	// 获取ID
	k := strings.Split(url, "-")
	if len(k) < 3 {
		return nil, errors.Errorf("url %s err", url)
	}
	key := k[2]

	cfg, _ := config.LoadConfig()
	if cfg.APIToken == "" {
		k8sClient, err := k8s.NewSimpleClient()
		if err != nil {
			return nil, err
		}
		cneapiDeploy, err := k8sClient.GetDeployment(context.Background(), common.GetDefaultSystemNamespace(true), "qucheng", metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		for _, e := range cneapiDeploy.Spec.Template.Spec.Containers[0].Env {
			if e.Name == "CNE_API_TOKEN" {
				cfg.APIToken = e.Value
				break
			}
		}
		cfg.SaveConfig()
	}

	apiHost := cfg.Domain
	if useip || apiHost == "" {
		apiHost = fmt.Sprintf("http://%s:32379", exnet.LocalIPs()[0])
	} else if !kutil.IsLegalDomain(apiHost) {
		apiHost = fmt.Sprintf("http://console.%s", cfg.Domain)
	} else {
		apiHost = fmt.Sprintf("https://%s", apiHost)
	}

	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG())
	if debug {
		client = client.DevMode().EnableDumpAll()
	}
	var result Result
	resp, err := client.R().
		SetHeader("accept", "application/json").
		SetHeader("TOKEN", cfg.APIToken).
		Get(fmt.Sprintf("%s/instance-apidetail-%s.html", apiHost, key))
	if err != nil {
		return nil, errors.Errorf("fetch api failed, reason: %v", err)
	}
	if !resp.IsSuccessState() {
		return nil, errors.Errorf("fetch api failed, reason: bad response status %v", resp.Status)
	}
	json.Unmarshal([]byte(resp.String()), &result)
	return &result.Data, nil
}
