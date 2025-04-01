// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package crd

import (
	"context"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"

	quchengv1beta1 "github.com/easysoft/quickon-api/qucheng/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type action struct {
	Name string
}

// cmdDbsList list database
func cmdDbsList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "list",
		Short:   "list database",
		Example: fmt.Sprintf(`%s platform crd db list`, os.Args[0]),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.LoadConfig()
			qclient, err := k8s.NewSimpleQClient()
			if err != nil {
				return err
			}
			dbsList, err := qclient.ListDB(context.TODO(), corev1.NamespaceAll, metav1.ListOptions{})
			if err != nil {
				return err
			}
			if len(dbsList.Items) == 0 {
				log.Warn("no found database")
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
				Label: "select database",
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
					return errors.Errorf("call kube api err: %v", err)
				}
				url := fmt.Sprintf("%s/adminer/?server=%s&username=%s&db=%s&password=%s", kutil.GetConsoleURL(cfg), dbs[it].Status.Address, dbs[it].Spec.Account.User.Value, "", dbs[it].Spec.Account.Password.Value)
				log.Infof("open browser access url: %s", url)
				// if err := browser.OpenURL(url); err == nil {
				// 	log.Done("open browser success")
				// 	return nil
				// }
				// log.Donef("open browser url: %s", url)
				return nil
			}
			return nil
		},
	}
	return app
}

// cmdDbSvcList list dbservice
func cmdDbSvcList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var onlygdb bool
	app := &cobra.Command{
		Use:     "list",
		Short:   "list db instance service",
		Example: fmt.Sprintf(`%s platform crd dbsvc list`, os.Args[0]),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.LoadConfig()
			qclient, err := k8s.NewSimpleQClient()
			if err != nil {
				return err
			}
			dbsvcs, err := qclient.ListDBSvc(context.TODO(), corev1.NamespaceAll, metav1.ListOptions{})
			if err != nil {
				return err
			}
			if len(dbsvcs.Items) == 0 {
				log.Warn("no found database instance service")
				return nil
			}
			var gdbServices []quchengv1beta1.DbService
			for _, dbsvc := range dbsvcs.Items {
				if !*dbsvc.Status.Ready {
					continue
				}
				if onlygdb {
					if vaildGlobalDatabase(dbsvc.Labels) {
						gdbServices = append(gdbServices, dbsvc)
					}
				} else {
					gdbServices = append(gdbServices, dbsvc)
				}
			}
			selectGDB := promptui.Select{
				Label: "select dbservice",
				Items: gdbServices,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Status.Address }})",
					Inactive: "  {{ .Name | cyan }}",
					Selected: "\U0001F389 {{ .Name | green | cyan }} ({{ .Status.Address }})",
				},
				Size: 5,
			}
			it, _, _ := selectGDB.Run()
			name := gdbServices[it].Name
			namespace := gdbServices[it].Namespace
			address := gdbServices[it].Status.Address
			username := gdbServices[it].Spec.Account.User.Value
			if len(username) == 0 && gdbServices[it].Spec.Account.User.ValueFrom != nil {
				// get secret
				secret, err := qclient.GetSecret(context.TODO(), namespace, gdbServices[it].Spec.Account.User.ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				username = string(secret.Data[gdbServices[it].Spec.Account.User.ValueFrom.SecretKeyRef.Key])
			}
			password := gdbServices[it].Spec.Account.Password.Value
			if len(password) == 0 && gdbServices[it].Spec.Account.Password.ValueFrom != nil {
				// get secret
				secret, err := qclient.GetSecret(context.TODO(), namespace, gdbServices[it].Spec.Account.Password.ValueFrom.SecretKeyRef.Name, metav1.GetOptions{})
				if err != nil {
					return err
				}
				password = string(secret.Data[gdbServices[it].Spec.Account.Password.ValueFrom.SecretKeyRef.Key])
			}
			url := fmt.Sprintf("%s/adminer/?server=%s&username=%s&db=%s&password=%s", kutil.GetConsoleURL(cfg), address, username, "", password)
			log.Infof("connect dbservice: %s\taddress: %s\n\tusername:%s\n\tpassword:%s\n\tadminer:%s", name, address, username, password, url)
			return nil
		},
	}
	app.Flags().BoolVar(&onlygdb, "onlygdb", false, "only show db service")
	return app
}
