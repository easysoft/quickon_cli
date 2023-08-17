// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package exec

import (
	"fmt"
	"os"
	sysexec "os/exec"
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/environ"
	"github.com/sirupsen/logrus"
)

func trace(cmd *sysexec.Cmd) {
	log := log.GetFileLogger(fmt.Sprintf("trace.%s.log", common.Version))
	if environ.GetEnv("QTRACE_DISABLE", "false") == "false" {
		key := strings.Join(cmd.Args, " ")
		log.Debugf("%s", key)
	}
}

func Command(name string, arg ...string) *sysexec.Cmd {
	if log.GetInstance().GetLevel() == logrus.DebugLevel {
		arg = append(arg, "--debug")
	}
	cmd := sysexec.Command(name, arg...) // #nosec
	// cmd.Dir = common.GetDefaultCacheDir()
	trace(cmd)
	return cmd
}

func CommandRun(name string, arg ...string) error {
	cmd := sysexec.Command(name, arg...) // #nosec
	trace(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func CommandBashRunWithResp(cmdStr string) (string, error) {
	cmd := sysexec.Command("/bin/bash", "-c", cmdStr) // #nosec
	// cmd.Dir = common.GetDefaultCacheDir()
	trace(cmd)
	result, err := cmd.CombinedOutput()
	return string(result), err
}
