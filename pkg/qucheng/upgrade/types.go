// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

type Version struct {
	Components []ComponentVersion `json:",omitempty"`
}

type ComponentVersion struct {
	Name           string
	Deploy         CVersion
	Remote         CVersion
	CanUpgrade     bool   `json:",omitempty"`
	UpgradeMessage string `json:",omitempty"`
}

type CVersion struct {
	AppVersion   string
	ChartVersion string
}
