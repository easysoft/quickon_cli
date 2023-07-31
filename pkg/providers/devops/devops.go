// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package devops

import (
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/pkg/providers"
)

const providerName = "devops"

type Devops struct {
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *Devops {
	return &Devops{}
}

func (q *Devops) GetProviderName() string {
	return providerName
}

func (q *Devops) GetFlags() []types.Flag {
	return nil
}

func (q *Devops) Install() error {
	return nil
}

func (q *Devops) Show() error {
	return nil
}
