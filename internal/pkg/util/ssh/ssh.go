// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh

import (
	"net"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/wait"
)

var defaultBackoff = wait.Backoff{
	Duration: 15 * time.Second,
	Factor:   1,
	Steps:    5,
}

type Interface interface {
	// Copy is copy local files to remote host
	// scp -r /tmp root@192.168.0.2:/root/tmp => Copy("192.168.0.2","tmp","/root/tmp")
	// need check md5sum
	Copy(host, srcFilePath, dstFilePath string) error
	// CmdAsync is exec command on remote host, and asynchronous return logs
	CmdAsync(host string, cmd ...string) error
	// Cmd is exec command on remote host, and return combined standard output and standard error
	Cmd(host, cmd string) ([]byte, error)
	//CmdToString is exec command on remote host, and return spilt standard output and standard error
	CmdToString(host, cmd, spilt string) (string, error)
	Ping(host string) error
}

type SSH struct {
	isStdout   bool
	User       string
	Password   string
	PkFile     string
	PkData     string
	PkPassword string
	Timeout    time.Duration

	// private properties
	localAddress *[]net.Addr
	clientConfig *ssh.ClientConfig
	log          log.Logger
}

func NewSSHClient(ssh *config.SSH, isStdout bool) Interface {
	log := log.GetInstance()
	if ssh.User == "" {
		ssh.User = common.DefaultOSUserRoot
	}
	address, err := listLocalHostAddrs()
	// todo: return error?
	if err != nil {
		log.Warnf("failed to get local address, %v", err)
	}
	return &SSH{
		isStdout:     isStdout,
		User:         ssh.User,
		Password:     ssh.Passwd,
		PkFile:       ssh.Pk,
		PkData:       ssh.PkData,
		PkPassword:   ssh.PkPasswd,
		localAddress: address,
		log:          log,
	}
}

type Client struct {
	SSH  Interface
	Host string
}
