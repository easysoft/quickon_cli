// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package quickon

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/providers"
	"github.com/easysoft/qcadmin/pkg/quickon"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/expass"
	"github.com/sirupsen/logrus"
	"k8s.io/kubectl/pkg/util/i18n"
	"k8s.io/kubectl/pkg/util/templates"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
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
			Type:            common.ZenTaoOSSType.String(),
			Version:         common.DefaultQuickonOSSVersion,
			App:             "zentao",
			DomainType:      "custom",
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
		Usage: "quickon version",
		P:     &q.MetaData.Version,
		V:     q.MetaData.Version,
	}, types.Flag{
		Name:  "app",
		Usage: "default quickon install app",
		P:     &q.MetaData.App,
		V:     q.MetaData.App,
	}, types.Flag{
		Name:      "type",
		Usage:     "quickon type",
		P:         &q.MetaData.Type,
		V:         common.ZenTaoOSSType.String(),
		ShortHand: "t",
	})
	return fs
}

func (q *Quickon) Install() error {
	if err := q.MetaData.Init(); err != nil {
		return errors.Errorf("init quickon data error: %v", err)
	}
	return qcexec.CommandRun(os.Args[0], "quickon", "app", "install", "--name", q.MetaData.App, "--api-useip", fmt.Sprintf("--debug=%v", q.MetaData.Log.GetLevel() == logrus.DebugLevel))
}

func (q *Quickon) Show() {
	if len(q.MetaData.IP) <= 0 {
		q.MetaData.IP = exnet.LocalIPs()[0]
	}
	resetPassArgs := []string{"quickon", "reset-password", "--password", q.MetaData.ConsolePassword}
	qcexec.CommandRun(os.Args[0], resetPassArgs...)
	cfg, _ := config.LoadConfig()
	cfg.ConsolePassword = q.MetaData.ConsolePassword
	cfg.SaveConfig()
	domain := cfg.Domain

	q.MetaData.Log.Info("----------------------------\t")
	if len(domain) > 0 {
		if !kutil.IsLegalDomain(cfg.Domain) {
			domain = fmt.Sprintf("http://console.%s", cfg.Domain)
		} else {
			domain = fmt.Sprintf("https://%s", cfg.Domain)
		}
	} else {
		domain = fmt.Sprintf("http://%s:32379", q.MetaData.IP)
	}
	q.MetaData.Log.Donef("console: %s, username: %s, password: %s",
		color.SGreen(domain), color.SGreen(common.QuchengDefaultUser), color.SGreen(q.MetaData.ConsolePassword))
	q.MetaData.Log.Donef("docs: %s", common.QuchengDocs)
	q.MetaData.Log.Done("support: 768721743(QQGroup)")
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

// GetUsageExample returns quickon usage example prompt.
func (q *Quickon) GetUsageExample() string {
	return templates.Examples(i18n.T(`
	# init quickon platform use example domain your.example.quickon.domain
	q init --provider quickon --domain your.example.quickon.domain`))
}
