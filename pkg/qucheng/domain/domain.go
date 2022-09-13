// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package domain

import (
	"fmt"
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exid"
	"github.com/imroc/req/v3"
	"github.com/manifoldco/promptui"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReqBody struct {
	IP        string `json:"ip"`
	UUID      string `json:"uuid"`
	SecretKey string `json:"secretKey"`
	Domain    string `json:"domain,omitempty"`
}

type RespBody struct {
	Code int `json:"code"`
	Data struct {
		Domain string `json:"domain,omitempty"`
		TLS    string `json:"tls,omitempty"`
	} `json:"data"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

func SearchCustomDomain(iip, id, secretKey string) string {
	var respbody RespBody
	reqbody := ReqBody{
		IP:        iip,
		UUID:      id,
		SecretKey: secretKey,
	}
	client := req.C().SetUserAgent(common.GetUG())
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&respbody).
		SetBody(&reqbody).
		Post(common.GetAPI("/api/qdns/oss/custom"))

	if err != nil {
		return ""
	}
	return respbody.Data.Domain
}

// UpgradeTLSDDomain tls domain
func UpgradeTLSDDomain(iip, id, secretKey, domain string) error {
	var respbody RespBody
	if kutil.IsLegalDomain(domain) {
		domain = strings.TrimSuffix(domain, ".haogs.cn")
		domain = strings.TrimSuffix(domain, ".corp.cc")
	}
	reqbody := ReqBody{
		IP:        iip,
		UUID:      id,
		SecretKey: secretKey,
		Domain:    domain,
	}
	client := req.C().SetUserAgent(common.GetUG())
	_, err := client.R().
		SetResult(&respbody).
		SetBody(&reqbody).
		Post(common.GetAPI("/api/qdns/oss/tls"))
	return err
}

// GenerateDomain generate suffix domain
func GenerateDomain(iip, id, secretKey, domain string) (string, string, error) {
	log := log.GetInstance()
	var respbody RespBody
	if kutil.IsLegalDomain(domain) {
		domain = strings.TrimSuffix(domain, ".haogs.cn")
		domain = strings.TrimSuffix(domain, ".corp.cc")
	}
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
			log.Warnf("current domain %s is unavailable, please try again", color.SRed("%s.haogs.cn", reqbody.Domain))
			return GenerateDomain(iip, id, secretKey, GenCustomDomain(domain))
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

func GenCustomDomain(defaultDomain string) string {
	log := log.GetInstance()
	prompt := promptui.Prompt{
		Label:   "config custom subdomain for haogs.cn, like: <custom>.haogs.cn.\t",
		Default: defaultDomain,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }}",
			Valid:   "{{ . | green }}",
			Invalid: "{{ . | red }}",
			Success: "{{ . | bold }}",
		},
		Validate: func(input string) error {
			if !kutil.IsLegalDomain(input) {
				input = fmt.Sprintf("%s.haogs.cn", input)
			}
			if len(input) < 13 {
				return fmt.Errorf("subdomain must be at least 4 characters, like %s", defaultDomain)
			}
			if msgs := validation.NameIsDNSSubdomain(input, false); len(msgs) != 0 {
				return fmt.Errorf("%s", msgs[0])
			}
			return nil
		},
	}
	result, _ := prompt.Run()
	if result == "" {
		log.Donef("use default domain: %s", color.SGreen(defaultDomain))
		result = defaultDomain
	}
	result = strings.ReplaceAll(result, ".haogs.cn", "")
	log.Infof("check subdomain %s availability", result)
	return result
}
