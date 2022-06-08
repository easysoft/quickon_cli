// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package exec

import (
	"os"
	sysexec "os/exec"
	"strings"

	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	elog "github.com/ergoapi/log"
	"github.com/ergoapi/util/environ"
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
	cmd := sysexec.Command(name, arg[:]...) // #nosec
	cmd.Stdin = os.Stdin
	cmd.Stderr = NewLogWrite(log.Flog, "err")
	cmd.Stdout = NewLogWrite(log.Flog, "")
	return cmd.Run()
}

func Trace(cmd *sysexec.Cmd) {
	if environ.GetEnv("TRACE", "false") == "true" {
		key := strings.Join(cmd.Args, " ")
		log.Flog.Debugf("+ %s\n", key)
	}
}

func Command(name string, arg ...string) *sysexec.Cmd {
	cmd := sysexec.Command(name, arg...) // #nosec
	Trace(cmd)
	return cmd
}
