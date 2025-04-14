// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"context"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/confirm"
	"github.com/ergoapi/util/exnet"
	"github.com/ergoapi/util/file"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/cmd/precheck"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/api/statistics"
	"github.com/easysoft/qcadmin/internal/pkg/status"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/pkg/cluster"
)

var (
	initExample = templates.Examples(`
		# init default cluster
		z cluster init

		# init cluster with custom cidr
		z cluster init --pod-cidr 10.100.0.0/16 --service-cidr 10.200.0.0/16

		# init cluster use mysql as datastore
		z cluster init --datastore "mysql://root:pass4Zenta0Pass@tcp(192.168.99.31:3306)/"

		# init cluster use postgres as datastore
		z cluster init --datastore "postgres://postgres:pass4Zenta0Pass@192.168.99.31:5432/"

		# more args
		z cluster init --help
	`)
)

func InitCommand(f factory.Factory) *cobra.Command {
	var preCheck precheck.PreCheck
	myCluster := cluster.NewCluster(f)
	init := &cobra.Command{
		Use:     "init",
		Short:   "init cluster",
		Example: initExample,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// 禁止重复初始化
			if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
				return errors.New("cluster is already initialized")
			}
			if len(myCluster.MasterIPs) == 0 {
				myCluster.MasterIPs = append(myCluster.MasterIPs, exnet.LocalIPs()[0])
			}
			preCheck.IgnorePreflightErrors = myCluster.IgnorePreflightErrors
			preCheck.OffLine = myCluster.OffLine
			if err := preCheck.Run(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := myCluster.InitNode(); err != nil {
				return err
			}
			// statistics.SendStatistics("install")
			return nil
		},
	}
	init.Flags().AddFlagSet(flags.ConvertFlags(init, myCluster.GetInitFlags()))
	return init
}

func CleanCommand(f factory.Factory) *cobra.Command {
	myCluster := cluster.NewCluster(f)
	log := f.GetLog()
	clean := &cobra.Command{
		Use:     "clean",
		Short:   "clean cluster",
		Version: "2.0.4",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, _ := confirm.Confirm("Are you sure to clean cluster")
			if status {
				if err := myCluster.Clean(); err != nil {
					return err
				}
				log.Donef("uninstall cluster success")
				statistics.SendStatistics("uninstall-cluster")
				return nil
			}
			log.Donef("cancel clean cluster")
			return nil
		},
	}
	return clean
}

func StatusCommand(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var params = status.K8sStatusOption{
		Log:         log,
		OnlyCluster: true,
	}
	status := &cobra.Command{
		Use:   "status",
		Short: "status cluster",
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
			s.Format()
			if err != nil {
				log.Fatalf("Unable to determine status:  %s", err)
			}
			return err
		},
	}
	status.Flags().StringVarP(&params.KubeConfig, "kubeconfig", "c", common.GetKubeConfig(), "Kubernetes configuration file")
	status.Flags().BoolVar(&params.Wait, "wait", false, "Wait for status to report success (no errors and warnings)")
	status.Flags().DurationVar(&params.WaitDuration, "wait-duration", common.StatusWaitDuration, "Maximum time to wait for status")
	status.Flags().BoolVar(&params.IgnoreWarnings, "ignore-warnings", false, "Ignore warnings when waiting for status to report success")
	status.Flags().StringVarP(&params.ListOutput, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	return status
}

func StopCommand(f factory.Factory) *cobra.Command {
	myCluster := cluster.NewCluster(f)
	log := f.GetLog()
	stop := &cobra.Command{
		Use:     "stop",
		Short:   "stop cluster",
		Version: "3.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, _ := confirm.Confirm("Are you sure to stop cluster")
			if status {
				return myCluster.Stop()
			}
			log.Donef("cancel stop cluster")
			return nil
		},
	}
	return stop
}

func StartUPCommand(f factory.Factory) *cobra.Command {
	myCluster := cluster.NewCluster(f)
	log := f.GetLog()
	stop := &cobra.Command{
		Use:     "start",
		Short:   "startup cluster",
		Aliases: []string{"startup"},
		Version: "4.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, _ := confirm.Confirm("Are you sure to start cluster")
			if status {
				return myCluster.StartUP()
			}
			log.Donef("cancel start cluster")
			return nil
		},
	}
	return stop
}
