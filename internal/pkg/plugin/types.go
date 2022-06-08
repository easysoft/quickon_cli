// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package plugin

import "github.com/easysoft/qcadmin/internal/pkg/k8s"

type Meta struct {
	Type    string `json:"type"`
	Default string `json:"default"`
	Item    []Item `json:"item"`
}

type Item struct {
	Client      *k8s.Client `json:"-"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Home        string      `json:"home"`
	Appversion  string      `json:"appversion"`
	Type        string      `json:"type"`
	Path        string      `json:"path"`
	Tool        string      `json:"tool"`
}

type List []Meta
