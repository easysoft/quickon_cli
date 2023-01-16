// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"html/template"
	"os"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/file"
)

type K3sArgs struct {
	TypeMaster  bool
	Master0     bool
	KubeAPI     string
	ClusterCIDR string
	ServiceCIDR string
	DataStore   string
	DataDir     string
	Docker      bool
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
	if k3s.DataDir == "" {
		k3s.DataDir = "/data/k3s"
	}
	return render(k3s, template)
}

func (k3s K3sArgs) Template() string {
	return common.K3SServiceTpl
}

func main() {
	k3sMasterArgs := K3sArgs{
		TypeMaster:  true,
		Master0:     true,
		KubeAPI:     "k.local",
		ClusterCIDR: "10.88.0.0/16",
		ServiceCIDR: "10.89.0.0/16",
		DataStore:   "mysql://",
		DataDir:     "",
		Docker:      true,
	}
	log := log.GetInstance()
	f, _ := os.CreateTemp("/tmp", "")
	log.Infof("file: %s", f.Name())
	file.Writefile(f.Name(), k3sMasterArgs.Manifests(""))
}
