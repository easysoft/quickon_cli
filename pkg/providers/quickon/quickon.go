// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package quickon

import (
	"fmt"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/providers"
	"github.com/easysoft/qcadmin/pkg/quickon"
	"github.com/ergoapi/util/expass"
)

const providerName = "quickon"

type Quickon struct {
	MetaData *quickon.Meta
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *Quickon {
	return &Quickon{
		MetaData: &quickon.Meta{
			Log:             log.GetInstance(),
			DevopsMode:      false,
			ConsolePassword: expass.PwGenAlphaNum(32),
			QuickonType:     common.QuickonOSSType,
		},
	}
}

func (q *Quickon) GetProviderName() string {
	return providerName
}

func (q *Quickon) GetFlags() []types.Flag {
	fs := q.MetaData.GetCustomFlags()
	fs = append(fs, types.Flag{
		Name:  "password",
		Usage: "quickon console password",
		P:     &q.MetaData.ConsolePassword,
		V:     q.MetaData.ConsolePassword,
	}, types.Flag{
		Name:  "version",
		Usage: fmt.Sprintf("quickon version(oss: %s/ee: %s)", common.DefaultQuickonOSSVersion, common.DefaultQuickonEEVersion),
		P:     &q.MetaData.Version,
		V:     q.MetaData.Version,
	})
	return fs
}

func (q *Quickon) Install() error {
	return q.MetaData.Init()
}

func (q *Quickon) Show() {
	// quickon show
	q.MetaData.Show()
}

func (q *Quickon) GetKubeClient() error {
	return q.MetaData.GetKubeClient()
}

func (q *Quickon) Check() error {
	return q.MetaData.Check()
}

func (q *Quickon) GetMeta() *quickon.Meta {
	return q.MetaData
}
