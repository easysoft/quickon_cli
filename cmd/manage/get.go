// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package manage

import (
	"os"

	"github.com/easysoft/qcadmin/internal/app/debug"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdGetNode(f factory.Factory) *cobra.Command {
	node := &cobra.Command{
		Use:     "node",
		Aliases: []string{"no", "nodes"},
		Short:   "get nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			// node id为0 list
			if len(args) == 0 {
				return qcexec.CommandRun(os.Args[0], "status", "node")
			}
			extargs := []string{"exp", "kubectl", "get", "-o", "wide", "nodes"}
			extargs = append(extargs, args...)
			return qcexec.CommandRun(os.Args[0], extargs...)
		},
	}
	return node
}

func NewCmdGetApp(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "app",
		Aliases: []string{"apps"},
		Short:   "get app",
		Args:    cobra.ExactArgs(1),
		Example: `q get app http://console.efbb.haogs.cn/instance-view-39.html`,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			apidebug := log.GetLevel() == logrus.DebugLevel
			log.Infof("start fetch app: %s", url)
			appdata, err := debug.GetNameByURL(url, apidebug)
			if err != nil {
				return err
			}
			// if apidebug {
			// 	spew.Dump(appdata)
			// }
			extargs := []string{"exp", "kubectl", "get", "-o", "wide", "pods,deploy,pvc,svc,ing", "-l", "release=" + appdata.K8Name}
			// extargs = append(extargs, args...)
			return qcexec.CommandRun(os.Args[0], extargs...)
		},
	}
	return app
}
