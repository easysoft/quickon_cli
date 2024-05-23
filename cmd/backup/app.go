// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package backup

import (
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/api/cne"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/output"
)

func NewCmdBackupApp(f factory.Factory) *cobra.Command {
	bc := &cobra.Command{
		Use:   "app",
		Short: "backup app",
		Long:  "backup app",
	}
	bc.AddCommand(newCmdBackupAppCreate(f))
	bc.AddCommand(newCmdBackupAppList(f))
	return bc
}

func newCmdBackupAppCreate(f factory.Factory) *cobra.Command {
	var app, ns, backupName string
	var err error
	log := f.GetLog()
	bc := &cobra.Command{
		Use:   "app",
		Short: "backup app",
		Long:  "backup app",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(app) == 0 || len(ns) == 0 {
				return errors.New("missing app or ns")
			}
			if !kutil.ValidNamespace(ns) {
				return errors.New("allow support quickon prefix namespace")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("start backup app: %s", app)
			cneClient := cne.NewCneAPI()
			if backupName == "" {
				backupName, err = cneClient.CreateAppBackUP(ns, app)
				if err != nil {
					log.Errorf("backup app %s failed, reason: %v", app, err)
					return
				}
			}

			timeout := 5 * time.Minute
			costtime := time.Now()
			deadline := time.Now().Add(timeout)
			for {
				backupStatus, err := cneClient.AppBackUPStatus(ns, app, backupName)
				if err != nil {
					log.Errorf("backup app %s failed, reason: %v", app, err)
					return
				}
				if strings.ToLower(backupStatus.Status) == "completed" {
					log.Infof("backup app %s(%s) success, cost: %vs", app, backupName, time.Since(costtime).Seconds())
					return
				}
				if strings.ToLower(backupStatus.Status) == "failed" {
					log.Errorf("backup app %s(%s) failed, reason: %s", app, backupName, backupStatus.Reason)
					return
				}
				if time.Now().After(deadline) {
					log.Errorf("backup app %s(%s) timeout", app, backupName)
					return
				}
				log.Infof("backup app %s(%s) status: %s", app, backupName, backupStatus.Status)
				time.Sleep(5 * time.Second)
			}
		},
	}
	bc.Flags().StringVar(&app, "app", "", "app chart name")
	bc.Flags().StringVarP(&ns, "namespace", "n", "", "namespace")
	bc.Flags().StringVarP(&backupName, "backupName", "b", "", "existing backup name")
	return bc
}

func newCmdBackupAppList(f factory.Factory) *cobra.Command {
	var app, ns, show string
	log := f.GetLog()
	bc := &cobra.Command{
		Use:   "list",
		Short: "list app backup",
		Long:  "list app backup",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(app) == 0 || len(ns) == 0 {
				return errors.New("missing app or ns")
			}
			if !kutil.ValidNamespace(ns) {
				return errors.New("allow support quickon prefix namespace")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cneClient := cne.NewCneAPI()
			data, err := cneClient.ListAppBackUP(ns, app)
			if err != nil {
				log.Errorf("backup app %s failed, reason: %v", app, err)
				return nil
			}
			switch strings.ToLower(show) {
			case "json":
				return output.EncodeJSON(os.Stdout, data)
			case "yaml":
				return output.EncodeYAML(os.Stdout, data)
			default:
				log.Infof("list backup app: %s", app)
				table := uitable.New()
				table.MaxColWidth = 80
				table.Wrap = true
				for _, d := range data {
					table.AddRow("name: ", d.Name)
					table.AddRow("status: ", d.Status)
					table.AddRow("restore: ", len(d.Restores))
					table.AddRow("dbs: ", len(d.BackupDetails.DBs))
					table.AddRow("volumes: ", len(d.BackupDetails.Volumes))
					table.AddRow("------------")
				}
				return output.EncodeTable(os.Stdout, table)
			}
		},
	}
	bc.Flags().StringVar(&app, "app", "", "app chart name")
	bc.Flags().StringVarP(&ns, "namespace", "n", "", "namespace")
	bc.Flags().StringVarP(&show, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	return bc
}
