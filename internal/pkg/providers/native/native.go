// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package native

import (
	"fmt"
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/cluster"
	"github.com/easysoft/qcadmin/internal/pkg/providers"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/pkg/util/preflight"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/zos"

	utilsexec "k8s.io/utils/exec"
)

// providerName is the name of this provider.
const providerName = "native"

type Native struct {
	*cluster.Cluster
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *Native {
	c := cluster.NewCluster()
	c.Provider = providerName
	return &Native{
		Cluster: c,
	}
}

const (
	createUsageExample = `
	create default cluster:
		q init

	create custom cluster
		q init --podsubnet "10.42.0.0/16" \
 			--svcsubnet "10.43.0.0/16" \
			--eip 1.1.1.1  \
			--san kubeapi.k8s.io

	create ha cluster with mysql 5.7
		q init --podsubnet "10.42.0.0/16" \
			--svcsubnet "10.43.0.0/16" \
			--eip 1.1.1.1  \
			--san kubeapi.k8s.io \
			--cluster-data-source "mysql://username:password@tcp(hostname:3306)/database-name"
`
	joinUsageExample = `
	join node to cluster:

		# use k3s api & k3s nodetoken
		q join --cne-api <api address> --cne-token <api token>
`
)

// GetUsageExample returns native usage example prompt.
func (p *Native) GetUsageExample(action string) string {
	switch action {
	case "create":
		return createUsageExample
	case "join":
		return joinUsageExample
	default:
		return "not support"
	}
}

// GetCreateFlags returns native create flags.
func (p *Native) GetCreateFlags() []types.Flag {
	fs := p.GetCreateOptions()
	fs = append(fs, p.GetCreateExtOptions()...)
	fs = append(fs, p.GetCreateNativeOptions()...)
	return fs
}

// GetJoinFlags returns native join flags.
func (p *Native) GetJoinFlags() []types.Flag {
	return p.GetJoinOptions()
}

func (p *Native) GetProviderName() string {
	return p.Provider
}

// CreateCluster create cluster.
func (p *Native) CreateCluster() (err error) {
	if p.AddHelmRepo() != nil {
		return err
	}
	return p.InitCluster()
}

// JoinNode join node.
func (p *Native) JoinNode() (err error) {
	return p.JoinCluster()
}

func (p *Native) InitQucheng() error {
	return p.InstallQuCheng()
}

func (p *Native) CreateCheck(skip bool) error {
	if skip {
		p.Log.Warn("skip precheck")
		return nil
	}
	p.Log.Info("start pre-flight checks")
	if err := preflight.RunInitNodeChecks(utilsexec.New(), &p.Metadata, false); err != nil {
		return err
	}
	p.Log.Done("pre-flight checks passed")
	return nil
}

func (p *Native) PreSystemInit() error {
	p.Log.Info("start system init")
	if err := p.SystemInit(); err != nil {
		return err
	}
	p.Log.Done("system init success")
	return nil
}

// GenerateManifest generates manifest deploy command.
func (p *Native) GenerateManifest() []string {
	// no need to support.
	return nil
}

// Show show cluster info.
func (p *Native) Show() {
	if len(p.Metadata.EIP) <= 0 {
		p.Metadata.EIP = exnet.LocalIPs()[0]
	}
	cfg, _ := config.LoadConfig()
	domain := ""
	if cfg != nil {
		cfg.DB = "sqlite"
		cfg.Token = kutil.GetNodeToken()
		cfg.InitNode = p.Metadata.EIP
		cfg.Master = []config.Node{
			{
				Name: zos.GetHostname(),
				Host: p.Metadata.EIP,
				Init: true,
			},
		}
		domain = cfg.Domain
		cfg.SaveConfig()
	}

	p.Log.Info("----------------------------")
	if len(domain) > 0 {
		if !strings.HasSuffix(cfg.Domain, "haogs.cn") {
			domain = fmt.Sprintf("console.%s", cfg.Domain)
		} else {
			domain = fmt.Sprintf("https://%s", cfg.Domain)
		}
	} else {
		domain = fmt.Sprintf("http://%s:32379", p.Metadata.EIP)
	}
	p.Log.Donef("web: %s, username: %s, password: %s",
		color.SGreen(domain), color.SGreen(common.QuchengDefaultUser), color.SGreen(common.QuchengDefaultPass))
	p.Log.Donef("docs: %s", common.QuchengDocs)
	p.Log.Done("support: 768721743(QQGroup)")
}

func (p *Native) SetLog(log log.Logger) {
	p.Log = log
}
