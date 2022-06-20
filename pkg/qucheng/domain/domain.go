// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package domain

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ergoapi/util/exid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateDomain generate suffix domain
func GenerateDomain(iip, id, secretKey string) (string, error) {
	body := make(url.Values)
	body["ip"] = []string{iip}
	body["uuid"] = []string{id}
	body["secretKey"] = []string{secretKey}
	resp, err := http.PostForm("https://api.qucheng.com/api/qdns", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
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
