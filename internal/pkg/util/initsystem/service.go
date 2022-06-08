// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package initsystem

import (
	"context"
	"os"

	"github.com/kardianos/service"
)

// Config configures the service.
type Config struct {
	Name string // service name
	Desc string // service description
	Dir  string
	Exec string
	Args []string
	Env  []string

	Stderr, Stdout string
}

type DaemonService struct {
	cancel context.CancelFunc
}

var nocontext = context.Background()

func (es *DaemonService) Start(s service.Service) error {
	_, cancel := context.WithCancel(nocontext)
	es.cancel = cancel
	return nil
}

func (es *DaemonService) Stop(s service.Service) error {
	es.cancel()
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}
