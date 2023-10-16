// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

import "github.com/easysoft/qcadmin/common"

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

type ZtUpgrade struct {
	Name    string
	Key     common.QuickonType
	Version string
}

var selectItems = []ZtUpgrade{
	{
		Name:    "开源版",
		Key:     common.ZenTaoOSSType,
		Version: common.DefaultZentaoDevOPSOSSVersion,
	},
	{
		Name:    "企业版",
		Key:     common.ZenTaoBizType,
		Version: common.DefaultZentaoDevOPSBizVersion,
	},
	{
		Name:    "旗舰版",
		Key:     common.ZenTaoMaxType,
		Version: common.DefaultZentaoDevOPSMaxVersion,
	},
	{
		Name:    "IPD",
		Key:     common.ZenTaoIPDType,
		Version: common.DefaultZentaoDevOPSIPDVersion,
	},
}
