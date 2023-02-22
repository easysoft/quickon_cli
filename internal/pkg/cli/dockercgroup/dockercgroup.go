// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package main

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/autodetect"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
)

func main() {
	log := log.GetInstance()
	// if err := autodetect.VerifyDockerDaemon(); err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// log.Done("docker daemon setting is ok")
	status := autodetect.VerifyCgroupDriverSystemd()
	if status {
		log.Done("docker daemon cgroup driver is systemd")
		return
	}
	log.Done("docker daemon cgroup driver is cgroupfs")
}
