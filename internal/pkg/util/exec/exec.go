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
	elog "github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/environ"
	"github.com/sirupsen/logrus"
)

type LogWriter struct {
	logger elog.Logger
	t      string
}

func NewLogWrite(logger elog.Logger, t string) *LogWriter {
	lw := &LogWriter{}
	lw.logger = logger
	return lw
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	if lw.t == "" {
		lw.logger.Debug(string(p))
	} else {
		lw.logger.Error(string(p))
	}
	return len(p), nil
}

func RunCmd(name string, arg ...string) error {
	log := log.GetInstance()
	cmd := sysexec.Command(name, arg[:]...) // #nosec
	// cmd.Dir = common.GetDefaultCacheDir()
	cmd.Stdin = os.Stdin
	cmd.Stderr = NewLogWrite(log, "err")
	cmd.Stdout = NewLogWrite(log, "")
	return cmd.Run()
}

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
	// cmd.Dir = common.GetDefaultCacheDir()
	return cmd.Run()
}

func CommandBashRunWithResp(cmdStr string) (string, error) {
	cmd := sysexec.Command("/bin/bash", "-c", cmdStr) // #nosec
	// cmd.Dir = common.GetDefaultCacheDir()
	trace(cmd)
	result, err := cmd.CombinedOutput()
	return string(result), err
}

// func CommandRespByte(command string, args ...string) ([]byte, error) {
// 	log := log.GetInstance()
// 	c := Command(command, args...)
// 	bytes, err := c.CombinedOutput()
// 	if err != nil {
// 		cmdStr := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
// 		log.Debugf("❌ Unable to execute %q:", cmdStr)
// 		if len(bytes) > 0 {
// 			log.Debugf(" %s", string(bytes))
// 		}
// 		return []byte{}, fmt.Errorf("unable to execute %q: %w", cmdStr, err)
// 	}

// 	return bytes, err
// }
