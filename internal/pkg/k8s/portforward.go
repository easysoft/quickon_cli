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

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"

	"github.com/pkg/browser"
)

func PortForwardCommand(ctx context.Context, ns, svc string, sport, dport int) error {
	log := log.GetInstance()
	args := []string{
		"experimental",
		"kubectl",
		"port-forward",
		"-n", ns,
		fmt.Sprintf("svc/%s", svc),
		"--address", "0.0.0.0",
		"--address", "::",
		fmt.Sprintf("%d:%d", dport, sport)}

	go func() {
		time.Sleep(5 * time.Second)
		url := fmt.Sprintf("http://localhost:%d", dport)

		// avoid cluttering stdout/stderr when opening the browser
		browser.Stdout = io.Discard
		browser.Stderr = io.Discard
		log.Infof("Opening %q in your browser...", url)
		browser.OpenURL(url)
	}()
	_, err := qcexec.CommandRespByte(os.Args[0], args...)
	return err
}
