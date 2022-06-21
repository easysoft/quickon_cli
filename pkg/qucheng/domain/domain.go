// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package domain

import (
	"github.com/easysoft/qcadmin/common"
	"github.com/ergoapi/util/exid"
	"github.com/imroc/req/v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReqBody struct {
	IP        string `json:"ip"`
	UUID      string `json:"uuid"`
	SecretKey string `json:"secretKey"`
}

// GenerateDomain generate suffix domain
func GenerateDomain(iip, id, secretKey string) (string, error) {
	var respbody struct {
		Code int `json:"code"`
		Data struct {
			Domain string `json:"domain"`
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
	}
	client := req.C().SetUserAgent(common.GetUG())
	_, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&respbody).
		SetBody(&reqbody).
		Post("https://api.qucheng.com/api/qdns/oss/record")

	if err != nil {
		return "", err
	}
	return respbody.Data.Domain, nil
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
