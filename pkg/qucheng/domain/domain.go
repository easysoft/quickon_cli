// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package domain

import (
	"fmt"

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
	IP         string `json:"ip"`
	SecretKey  string `json:"secretKey"`
	SubDomain  string `json:"sub,omitempty"`
	MainDomain string `json:"domain,omitempty"`
}

type RespBody struct {
	Code int `json:"code"`
	Data struct {
		Domain      string `json:"domain,omitempty"`
		K8sTLS      string `json:"k8s-tls,omitempty"`
		TLSCertPath string `json:"tls_cert_path,omitempty"`
		TLSKeyPath  string `json:"tls_key_path,omitempty"`
	} `json:"data"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

func SearchCustomDomain(iip, secretKey, domain string) string {
	var respbody RespBody
	subDomain, mainDomain := kutil.SplitDomain(domain)
	reqbody := ReqBody{
		IP:         iip,
		SecretKey:  secretKey,
		MainDomain: mainDomain,
		SubDomain:  subDomain,
	}
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG())
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetSuccessResult(&respbody).
		SetBody(&reqbody).
		Post(common.GetAPI("/api/qdnsv2/oss/custom"))

	if err != nil {
		return ""
	}
	return respbody.Data.Domain
}

// UpgradeTLSDDomain tls domain
func UpgradeTLSDDomain(iip, secretKey, domain string) error {
	var respbody RespBody
	if !kutil.IsLegalDomain(domain) {
		return fmt.Errorf("domain not allow")
	}
	subDomain, mainDomain := kutil.SplitDomain(domain)
	reqbody := ReqBody{
		IP:         iip,
		SecretKey:  secretKey,
		MainDomain: mainDomain,
		SubDomain:  subDomain,
	}
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG())
	_, err := client.R().
		SetSuccessResult(&respbody).
		SetBody(&reqbody).
		Post(common.GetAPI("/api/qdnsv2/oss/tls"))
	return err
}

// GenerateDomain generate suffix domain
func GenerateDomain(iip, secretKey, domain string) (string, string, error) {
	log := log.GetInstance()
	var respbody RespBody
	if !kutil.IsLegalDomain(domain) {
		return "", "", fmt.Errorf("domain not allow")
	}
	subDomain, mainDomain := kutil.SplitDomain(domain)
	reqbody := ReqBody{
		IP:         iip,
		SecretKey:  secretKey,
		SubDomain:  subDomain,
		MainDomain: mainDomain,
	}
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG())
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetSuccessResult(&respbody).
		SetBody(&reqbody).
		Post(common.GetAPI("/api/qdnsv2/oss/record"))

	if err != nil {
		return "", "", err
	}
	if len(respbody.Data.Domain) == 0 {
		if len(reqbody.SubDomain) > 0 {
			log.Warnf("current domain %s is unavailable, please try again", color.SRed("%s.%s", reqbody.SubDomain, reqbody.MainDomain))
			return GenerateDomain(iip, secretKey, GenCustomDomain(domain))
		}
	}
	return respbody.Data.Domain, respbody.Data.K8sTLS, nil
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

func GenCustomDomain(domain string) string {
	log := log.GetInstance()
	subDomain, mainDomain := kutil.SplitDomain(domain)
	prompt := promptui.Prompt{
		Label:   fmt.Sprintf("configure custom domain, like: %s.\t", domain),
		Default: domain,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }}",
			Valid:   "{{ . | green }}",
			Invalid: "{{ . | red }}",
			Success: "{{ . | bold }}",
		},
		Validate: func(input string) error {
			if !kutil.IsLegalDomain(input) {
				input = fmt.Sprintf("%s.%s", input, mainDomain)
			}
			if len(input) < 12 {
				return fmt.Errorf("subdomain must be at least 4 characters, like %s", subDomain)
			}
			if msgs := validation.NameIsDNSSubdomain(input, false); len(msgs) != 0 {
				return fmt.Errorf("%s", msgs[0])
			}
			return nil
		},
	}
	result, _ := prompt.Run()
	if result == "" {
		log.Donef("use default domain: %s", color.SGreen(domain))
		result = domain
	}
	// result = strings.ReplaceAll(result, fmt.Sprintf(".%s", mainDomain), "")
	log.Infof("check subdomain %s availability", color.SGreen(result))
	return result
}
