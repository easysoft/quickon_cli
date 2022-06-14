// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"fmt"

	"github.com/easysoft/qcadmin/common"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
)

func (p *Cluster) SystemInit() (err error) {
	initShell := fmt.Sprintf("%s/hack/manifests/scripts/init.sh", common.GetDefaultDataDir())
	log.Flog.Debugf("gen init shell: %v", initShell)
	if err := qcexec.RunCmd("/bin/bash", initShell); err != nil {
		return err
	}
	return nil
}
