// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package devops

import (
	"fmt"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/providers"
	"github.com/easysoft/qcadmin/pkg/quickon"
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
			Log:        log.GetInstance(),
			DevopsMode: true,
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
	// devops show
	q.MetaData.Show()
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
