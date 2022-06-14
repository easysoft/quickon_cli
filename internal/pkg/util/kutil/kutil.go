// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package kutil

import "io/ioutil"

const (
	NodeToken = "/var/lib/rancher/k3s/server/node-token"
)

func GetNodeToken() string {
	b, err := ioutil.ReadFile(NodeToken)
	if err != nil {
		return ""
	}
	return string(b)
}
