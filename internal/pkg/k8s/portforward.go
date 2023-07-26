// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package k8s

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/easysoft/qcadmin/common"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/zos"

	"github.com/pkg/browser"
)

func PortForwardCommand(ctx context.Context, ns, svc string, sport int) error {
	log := log.GetInstance()
	dport := exnet.GetFreePort()
	args := []string{
		"experimental",
		"kubectl",
		"port-forward",
		"-n", ns,
		fmt.Sprintf("svc/%s", svc),
		"--address", "0.0.0.0",
		"--address", "::",
		"--kubeconfig", common.GetKubeConfig(),
		fmt.Sprintf("%d:%d", dport, sport)}

	url := fmt.Sprintf("%s:%d", exnet.LocalIPs()[0], dport)
	log.Infof("listen to: %s", url)

	go func() {
		time.Sleep(5 * time.Second)

		// avoid cluttering stdout/stderr when opening the browser
		browser.Stdout = io.Discard
		browser.Stderr = io.Discard
		if zos.IsMacOS() {
			browser.OpenURL(url)
		}
	}()
	return qcexec.CommandRun(os.Args[0], args...)
}
