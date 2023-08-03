// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package quickon

import (
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/pkg/providers"
	"github.com/easysoft/qcadmin/pkg/quickon"
)

const providerName = "quickon"

type Quickon struct {
	quickon.Meta
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *Quickon {
	return &Quickon{}
}

func (q *Quickon) GetProviderName() string {
	return providerName
}

func (q *Quickon) GetFlags() []types.Flag {
	fs := q.GetCustomFlags()
	fs = append(fs, types.Flag{
		Name:  "password",
		Usage: "quickon console password",
		P:     &q.ConsolePassword,
		V:     q.ConsolePassword,
	})
	return fs
}

func (q *Quickon) Install() error {
	return q.Init()
}

func (q *Quickon) Show() error {
	return q.Show()
}

func (q *Quickon) GetKubeClient() error {
	return q.GetKubeClient()
}

func (q *Quickon) Check() error {
	return q.Check()
}

func (q *Quickon) GetMeta() *quickon.Meta {
	return &q.Meta
}
