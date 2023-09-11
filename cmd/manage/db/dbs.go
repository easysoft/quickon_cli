// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package db

import (
	"context"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	quchengv1beta1 "github.com/easysoft/quickon-api/qucheng/v1beta1"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func cmdDbsList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "dbs",
		Short:   "list dbs",
		Example: fmt.Sprintf(`%s platform db list dbs`, os.Args[0]),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.LoadConfig()
			qclient, err := k8s.NewSimpleQClient()
			if err != nil {
				return err
			}
			dbsList, err := qclient.ListQuchengDB(context.TODO(), corev1.NamespaceAll, metav1.ListOptions{})
			if err != nil {
				return err
			}
			if len(dbsList.Items) == 0 {
				log.Warn("no found database service")
				return nil
			}
			var dbs []quchengv1beta1.Db
			for _, db := range dbsList.Items {
				if !*db.Status.Ready {
					continue
				}
				dbs = append(dbs, db)
			}
			selectGDB := promptui.Select{
				Label: "select db",
				Items: dbs,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Status.Address }})",
					Inactive: "  {{ .Name | cyan }}",
					Selected: "\U0001F389 {{ .Name | green | cyan }} ({{ .Status.Address }})",
				},
				Size: 5,
			}
			it, _, _ := selectGDB.Run()
			actions := []action{{"manage"}}
			selectDBAction := promptui.Select{
				Label: "select action",
				Items: actions,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Name | cyan }}",
					Inactive: "  {{ .Name | cyan }}",
					Selected: fmt.Sprintf("\U0001F389 {{ .Name | green | cyan }} (%s)", dbs[it].Status.Address),
				},
			}
			iac, _, _ := selectDBAction.Run()
			if actions[iac].Name == "manage" {
				// https://console.example.corp.cc/adminer/?server=10.10.16.15%3A3306&username=root&db=ysicing&password=password123
				if err := fakeDbUserInfo(qclient, &dbs[it]); err != nil {
					return errors.Errorf("call k8s api err: %v", err)
				}
				url := fmt.Sprintf("%s/adminer/?server=%s&username=%s&db=%s&password=%s", kutil.GetConsoleURL(cfg), dbs[it].Status.Address, dbs[it].Spec.Account.User.Value, "", dbs[it].Spec.Account.Password.Value)
				log.Infof("open browser access url: %s", url)
				if err := browser.OpenURL(url); err != nil {
					log.Warnf("try open browser failed: %v", err)
					return nil
				}
				log.Done("open browser success")
				return nil
			}
			return nil
		},
	}
	return app
}

func fakeDbUserInfo(qclient *k8s.Client, db *quchengv1beta1.Db) error {
	if db.Spec.Account.User.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), db.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: db.Spec.Account.User.ValueFrom.SecretKeyRef.Name,
			},
			Key: db.Spec.Account.User.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		db.Spec.Account.User.Value = string(user)
	}
	if db.Spec.Account.Password.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), db.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: db.Spec.Account.Password.ValueFrom.SecretKeyRef.Name,
			},
			Key: db.Spec.Account.Password.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		db.Spec.Account.Password.Value = string(user)
	}
	return nil
}
