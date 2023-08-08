// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package devops

import (
	"fmt"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/providers"
	"github.com/easysoft/qcadmin/pkg/quickon"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
)

const providerName = "devops"

type Devops struct {
	MetaData *quickon.Meta
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *Devops {
	return &Devops{
		MetaData: &quickon.Meta{
			Log:         log.GetInstance(),
			DevopsMode:  true,
			QuickonType: common.QuickonOSSType,
			Version:     common.GetZenTaoVersion("", common.QuickonOSSType),
		},
	}
}

// GetProviderName returns the name of the provider
func (q *Devops) GetProviderName() string {
	return providerName
}

// GetFlags returns the flags of the provider
func (q *Devops) GetFlags() []types.Flag {
	fs := q.MetaData.GetCustomFlags()
	fs = append(fs, types.Flag{
		Name:  "version",
		Usage: fmt.Sprintf("zentao devops version %s", common.DefaultZentaoDevOPSOSSVersion),
		P:     &q.MetaData.Version,
		V:     q.MetaData.Version,
	})
	return fs
}

// Install installs the provider
func (q *Devops) Install() error {
	return q.MetaData.Init()
}

func (q *Devops) Show() {
	if len(q.MetaData.IP) <= 0 {
		q.MetaData.IP = exnet.LocalIPs()[0]
	}
	cfg, _ := config.LoadConfig()
	domain := cfg.Domain

	q.MetaData.Log.Info("----------------------------\t")
	if len(domain) > 0 {
		if !kutil.IsLegalDomain(cfg.Domain) {
			domain = fmt.Sprintf("http://zentao.%s", cfg.Domain)
		} else {
			domain = fmt.Sprintf("https://%s", cfg.Domain)
		}
	} else {
		domain = fmt.Sprintf("http://%s:32379", q.MetaData.IP)
	}
	q.MetaData.Log.Donef("console: %s", color.SGreen(domain))
	q.MetaData.Log.Donef("docs: %s", common.ZentaoDocs)
	q.MetaData.Log.Done("support: 768721743(QQGroup)")
}

func (q *Devops) GetKubeClient() error {
	return q.MetaData.GetKubeClient()
}

func (q *Devops) Check() error {
	return q.MetaData.Check()
}

func (q *Devops) GetMeta() *quickon.Meta {
	return q.MetaData
}
