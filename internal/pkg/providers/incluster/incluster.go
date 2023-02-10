// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package incluster

import (
	"fmt"
	"os"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/cluster"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/providers"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"
)

// providerName is the name of this provider.
const providerName = "incluster"

type InCluster struct {
	*cluster.Cluster
}

func init() {
	providers.RegisterProvider(providerName, func() (providers.Provider, error) {
		return newProvider(), nil
	})
}

func newProvider() *InCluster {
	c := cluster.NewCluster()
	c.Provider = providerName
	return &InCluster{
		Cluster: c,
	}
}

const createUsageExample = `
	create qucheng cluster:
		q init
`

// GetUsageExample returns native usage example prompt.
func (p *InCluster) GetUsageExample(action string) string {
	switch action {
	case "create":
		return createUsageExample
	default:
		return "not support"
	}
}

// GetCreateFlags returns cluster create flags.
func (p *InCluster) GetCreateFlags() []types.Flag {
	fs := p.GetCreateExtOptions()
	return fs
}

// GetJoinFlags returns cluster join flags.
func (p *InCluster) GetJoinFlags() []types.Flag {
	return nil
}

func (p *InCluster) GetProviderName() string {
	return p.Provider
}

// CreateCluster create cluster.
func (p *InCluster) CreateCluster() (err error) {
	p.Log.Warn("exists cluster, check cluster status")
	return p.AddHelmRepo()
}

// JoinNode join node.
func (p *InCluster) JoinNode() (err error) {
	return nil
}

func (p *InCluster) InitQucheng() (err error) {
	p.KubeClient, err = k8s.NewSimpleClient()
	if err != nil {
		return err
	}
	p.Log.Info("start init qucheng")
	if err := p.InstallQuCheng(); err != nil {
		return err
	}
	file.Writefile(common.GetCustomConfig(common.InitModeCluster), "in cluster ok", true)
	return nil
}

func (p *InCluster) CreateCheck(skip bool) error {
	// load kubeconfig check cluster ready
	return nil
}

func (p *InCluster) PreSystemInit() error {
	// no need to support.
	return nil
}

// GenerateManifest generates manifest deploy command.
func (p *InCluster) GenerateManifest() []string {
	// no need to support.
	return nil
}

// Show show cluster info.
func (p *InCluster) Show() {
	if len(p.Metadata.EIP) <= 0 {
		p.Metadata.EIP = exnet.LocalIPs()[0]
	}
	resetPassArgs := []string{"manage", "reset-password", "--password", p.Metadata.ConsolePassword}
	qcexec.CommandRun(os.Args[0], resetPassArgs...)
	cfg, _ := config.LoadConfig()
	domain := ""
	if cfg != nil {
		cfg.DB = "sqlite"
		cfg.ConsolePassword = p.Metadata.ConsolePassword
		cfg.Token = kutil.GetNodeToken()
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

	p.Log.Info("----------------------------\t")
	if len(domain) > 0 {
		if !kutil.IsLegalDomain(cfg.Domain) {
			domain = fmt.Sprintf("http://console.%s", cfg.Domain)
		} else {
			domain = fmt.Sprintf("https://%s", cfg.Domain)
		}
	} else {
		domain = fmt.Sprintf("http://%s:32379", p.Metadata.EIP)
	}
	p.Log.Donef("console: %s, username: %s, password: %s",
		color.SGreen(domain), color.SGreen(common.QuchengDefaultUser), color.SGreen(p.Metadata.ConsolePassword))

	p.Log.Donef("docs: %s", common.QuchengDocs)
	p.Log.Done("support: 768721743(QQGroup)")
}

func (p *InCluster) SetLog(log log.Logger) {
	p.Log = log
}
