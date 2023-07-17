// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package precheck

import (
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/internal/pkg/util/preflight"
	utilsexec "k8s.io/utils/exec"
)

type PreCheck struct {
	IgnorePreflightErrors bool
}

func (pc PreCheck) Run() error {
	log := log.GetInstance()
	log.Info("start pre-flight checks")
	if err := preflight.RunInitNodeChecks(utilsexec.New(), &types.Metadata{
		ServiceCidr: "10.42.0.0/16",
		ClusterCidr: "10.43.0.0/16",
	}, pc.IgnorePreflightErrors); err != nil {
		return err
	}
	return nil
}
