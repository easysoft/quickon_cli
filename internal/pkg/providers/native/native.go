// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package native

import (
	"fmt"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/cluster"
	"github.com/easysoft/qcadmin/internal/pkg/providers"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/pkg/util/preflight"
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
	log.Flog.Info("start init cluster")
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
		log.Flog.Warn("skip precheck")
		return nil
	}
	log.Flog.Info("start pre-flight checks")
	if err := preflight.RunInitNodeChecks(utilsexec.New(), &p.Metadata, false); err != nil {
		return err
	}
	log.Flog.Done("pre-flight checks passed")
	return nil
}

func (p *Native) PreSystemInit() error {
	log.Flog.Info("start system init")
	if err := p.SystemInit(); err != nil {
		return err
	}
	log.Flog.Done("system init passed")
	return nil
}

// GenerateManifest generates manifest deploy command.
func (p *Native) GenerateManifest() []string {
	// no need to support.
	return nil
}

// Show show cluster info.
func (p *Native) Show() {
	loginip := p.Metadata.EIP
	if len(loginip) <= 0 {
		loginip = exnet.LocalIPs()[0]
	}
	cfg, _ := config.LoadConfig()
	if cfg != nil {
		cfg.DB = "sqlite"
		cfg.Token = kutil.GetNodeToken()
		cfg.InitNode = loginip
		cfg.Master = []config.Node{
			{
				Name: zos.GetHostname(),
				Host: loginip,
				Init: true,
			},
		}
		cfg.SaveConfig()
	}

	log.Flog.Info("----------------------------")
	log.Flog.Donef("web:: %s", fmt.Sprintf("http://%s:32379", loginip))
	log.Flog.Donef("docs: %s", common.QuchengDocs)
}