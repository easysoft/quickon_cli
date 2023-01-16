// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package main

import (
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/ssh"
)

func main() {
	sshClient := ssh.NewSSHClient(&config.SSH{
		Passwd: "sshutil",
	}, true)
	if err := sshClient.Ping("127.0.0.1:10022"); err != nil {
		panic(err)
	}
	if err := sshClient.CmdAsync("127.0.0.1:10022", "pwd"); err != nil {
		panic(err)
	}
	sshClient1 := ssh.NewSSHClient(&config.SSH{
		Passwd: "sshutil",
		User:   "sshutil",
	}, true)
	if err := sshClient1.Ping("127.0.0.1:10023"); err != nil {
		panic(err)
	}
	if err := sshClient1.CmdAsync("127.0.0.1:10023", "pwd"); err != nil {
		panic(err)
	}
	if err := sshClient.CmdAsync("127.0.0.1:10024", "pwd"); err != nil {
		panic(err)
	}
}
