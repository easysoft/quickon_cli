// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package plugin

import (
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	pluginapi "github.com/easysoft/qcadmin/internal/pkg/plugin"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/spf13/cobra"
)

func InstallPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Short:   "install",
		Aliases: []string{"i"},
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := pluginapi.GetMeta(args...)
			if err != nil {
				return err
			}
			c, err := k8s.NewClient("", "")
			if err != nil {
				log.Flog.Fatal("connect k8s failed")
				return nil
			}
			ps.Client = c
			return ps.Install()
		},
	}
	return cmd
}
