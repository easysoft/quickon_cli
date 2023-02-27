// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"os"

	statussubcmd "github.com/easysoft/qcadmin/cmd/status"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/status"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"
)

func newCmdStatus(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var params = status.K8sStatusOption{
		Log: log,
	}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display status",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			defaultArgs := os.Args
			if !file.CheckFileExists(params.KubeConfig) {
				log.Warnf("not found cluster. just run %s init cluster", color.SGreen("%s init", defaultArgs[0]))
				os.Exit(0)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			collector, err := status.NewK8sStatusCollector(params)
			if err != nil {
				return err
			}
			s, err := collector.Status(context.Background())
			// Report the most recent status even if an error occurred.
			s.Format()
			if err != nil {
				log.Fatalf("Unable to determine status:  %s", err)
			}
			return err
		},
	}
	cmd.Flags().StringVarP(&params.KubeConfig, "kubeconfig", "c", common.GetKubeConfig(), "Kubernetes configuration file")
	cmd.Flags().BoolVar(&params.Wait, "wait", false, "Wait for status to report success (no errors and warnings)")
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", common.StatusWaitDuration, "Maximum time to wait for status")
	cmd.Flags().BoolVar(&params.IgnoreWarnings, "ignore-warnings", false, "Ignore warnings when waiting for status to report success")
	cmd.Flags().StringVarP(&params.ListOutput, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	cmd.AddCommand(statussubcmd.TopNodeCmd())
	return cmd
}
