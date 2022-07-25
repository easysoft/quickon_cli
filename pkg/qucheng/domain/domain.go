// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package domain

import (
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exid"
	"github.com/imroc/req/v3"
	"github.com/manifoldco/promptui"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReqBody struct {
	IP        string `json:"ip"`
	UUID      string `json:"uuid"`
	SecretKey string `json:"secretKey"`
	Domain    string `json:"domain"`
}

// GenerateDomain generate suffix domain
func GenerateDomain(iip, id, secretKey, domain string) (string, string, error) {
	log := log.GetInstance()
	var respbody struct {
		Code int `json:"code"`
		Data struct {
			Domain string `json:"domain"`
			TLS    string `json:"tls"`
		} `json:"data"`
		Message   string `json:"message"`
		Timestamp int    `json:"timestamp"`
	}
	// reqbody := struct {
	// 	IP        string `json:"ip"`
	// 	UUID      string `json:"uuid"`
	// 	SecretKey string `json:"secretKey"`
	// }{
	// 	IP:        iip,
	// 	UUID:      id,
	// 	SecretKey: secretKey,
	// }
	reqbody := ReqBody{
		IP:        iip,
		UUID:      id,
		SecretKey: secretKey,
		Domain:    domain,
	}
	client := req.C().SetUserAgent(common.GetUG())
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&respbody).
		SetBody(&reqbody).
		Post(common.GetAPI("/api/qdns/oss/record"))

	if err != nil {
		return "", "", err
	}
	if len(respbody.Data.Domain) == 0 {
		if len(reqbody.Domain) > 0 {
			log.Warnf("current domain %s is unavailable, please try again", color.SRed("%s.haogs,cn", reqbody.Domain))
			return GenerateDomain(iip, id, secretKey, GenCustomDomain())
		}
	}
	return respbody.Data.Domain, respbody.Data.TLS, nil
}

// GenerateSuffixConfigMap -
func GenerateSuffixConfigMap(name, namespace string) *corev1.ConfigMap {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"uuid": exid.GenUUID(),
			"auth": exid.GenUUID(),
		},
	}
	return cm
}

func GenCustomDomain() string {
	log := log.GetInstance()
	prompt := promptui.Prompt{
		Label:   "config custom alias domain for haogs.cn, like: <custom>.haogs.cn",
		Default: "",
	}
	result, _ := prompt.Run()
	if result == "" {
		log.Info("Platform random generation")
		return ""
	}
	log.Info("Check domain name availability")
	return result
}
