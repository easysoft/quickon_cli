// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"os"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/sirupsen/logrus"
)

func (p *Cluster) defaultAppInstall() {
	args := []string{"app", "install", p.ImportDefaultApp, "--api-useip"}
	if p.Log.GetLevel() == logrus.DebugLevel {
		args = append(args, "--debug")
	}
	if err := qcexec.CommandRun(os.Args[0], args...); err != nil {
		p.Log.Errorf("install default app %s err, reason: %v", p.ImportDefaultApp, err)
	}
}
