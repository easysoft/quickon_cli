// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package backup

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

func NewCmdBackupCluster(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	bc := &cobra.Command{
		Use:     "cluster",
		Short:   "backup cluster",
		Long:    "backup cluster",
		Aliases: []string{"snapshot"},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return errors.Errorf("found config err, reason: %v", err)
			}
			if cfg.DB != "" && cfg.DB != "etcd" {
				return errors.Errorf("not support datastore %s", cfg.DB)
			}
			log.Info("start backup cluster datastore etcd")
			return exec.CommandRun("bash", "-c", fmt.Sprintf("%s %s %s", common.GetCustomFile("hack/manifests/scripts/etcd-snapshot.sh"), cfg.DataDir, common.GetDefaultLogDir()))
		},
	}
	return bc
}
