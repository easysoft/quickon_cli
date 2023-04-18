// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package hostinfo

import (
	"fmt"

	"github.com/easysoft/qcadmin/common"
	hinfo "tailscale.com/hostinfo"
	"tailscale.com/tailcfg"
)

// New returns a partially populated Hostinfo for the current host.
func New() *tailcfg.Hostinfo {
	t := hinfo.New()
	t.IPNVersion = fmt.Sprintf("%s-%s", common.Version, common.GitCommitHash)
	return t
}
