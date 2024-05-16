// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package k3stpl

import (
	"bytes"
	"html/template"
	"os"

	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

type K3sArgs struct {
	TypeMaster   bool
	Master0      bool
	KubeAPI      string
	KubeToken    string
	PodCIDR      string
	ServiceCIDR  string
	DataStore    string
	DataDir      string
	LocalStorage bool
	CNI          string
	OffLine      bool
	Master0IP    string
	Registry     string
}

func render(data K3sArgs, temp string) string {
	var b bytes.Buffer
	t := template.Must(template.New("k3s").Parse(temp))
	_ = t.Execute(&b, &data)
	return b.String()
}

func (k3s K3sArgs) Manifests(template string) string {
	if template == "" {
		template = k3s.Template()
	}
	k3s.DataDir = common.GetDefaultQuickonPlatformDir(k3s.DataDir)
	if k3s.KubeToken == "" {
		k3s.KubeToken = "quickon"
	}
	if k3s.CNI != "" && k3s.CNI != "flannel" {
		if k3s.CNI == "wg" || k3s.CNI == "wireguard-native" || k3s.CNI == "wireguard" {
			k3s.CNI = "wireguard-native"
		} else {
			k3s.CNI = "none"
		}
	} else {
		k3s.CNI = ""
	}
	if k3s.OffLine && len(k3s.Master0IP) == 0 {
		k3s.Master0IP = exnet.LocalIPs()[0]
	}
	if k3s.Registry == "" {
		k3s.Registry = "hub.qucheng.com"
	}
	return render(k3s, template)
}

func (k3s K3sArgs) Template() string {
	return common.K3SServiceTpl
}

func EmbedCommand(f factory.Factory) *cobra.Command {
	var k3sargs K3sArgs
	rootCmd := &cobra.Command{
		Use: "k3stpl",
		Run: func(cmd *cobra.Command, args []string) {
			log := f.GetLog()
			tplfile, _ := os.CreateTemp("/tmp", "")
			log.Infof("file: %s", tplfile.Name())
			file.WriteFile(tplfile.Name(), k3sargs.Manifests(""), true)
		},
	}
	rootCmd.Flags().StringVar(&k3sargs.PodCIDR, "pod-cidr", "10.42.0.0/16", "cluster cidr")
	rootCmd.Flags().StringVar(&k3sargs.ServiceCIDR, "service-cidr", "10.43.0.0/16", "service cidr")
	rootCmd.Flags().StringVar(&k3sargs.DataDir, "data-dir", "", "data dir")
	rootCmd.Flags().StringVar(&k3sargs.DataStore, "data", "", "data type")
	rootCmd.Flags().StringVar(&k3sargs.KubeAPI, "kubeapi", "", "kubeapi")
	rootCmd.Flags().BoolVar(&k3sargs.TypeMaster, "master", true, "master type")
	rootCmd.Flags().BoolVar(&k3sargs.Master0, "master0", false, "master0")
	rootCmd.Flags().BoolVar(&k3sargs.OffLine, "offline", false, "offline")
	rootCmd.Flags().BoolVar(&k3sargs.LocalStorage, "local-storage", true, "local-storage")
	rootCmd.Flags().StringVar(&k3sargs.Master0IP, "master0ip", "", "master0ip, only work offline mode")
	rootCmd.Flags().StringVar(&k3sargs.Registry, "registry", "hub.qucheng.com", "registry")
	return rootCmd
}
