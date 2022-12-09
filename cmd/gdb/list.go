// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package gdb

import (
	"context"
	"fmt"

	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	quchengv1beta1 "github.com/easysoft/quickon-api/qucheng/v1beta1"
	"github.com/ergoapi/util/exmap"
	"github.com/manifoldco/promptui"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type action struct {
	Name string
}

func NewCmdGdbList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var name string
	app := &cobra.Command{
		Use:     "list",
		Short:   "list gdb",
		Example: `g gdb list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := config.LoadConfig()
			qclient, err := k8s.NewSimpleQClient()
			if err != nil {
				return err
			}
			dbsvcs, err := qclient.ListQuchengDBSvc(context.TODO(), name, metav1.ListOptions{})
			if err != nil {
				return err
			}
			if len(dbsvcs.Items) == 0 {
				log.Warn("no found global database service")
				return nil
			}
			var gdbServices []quchengv1beta1.DbService
			for _, dbsvc := range dbsvcs.Items {
				if vaildGlobalDatabase(dbsvc.Labels) {
					gdbServices = append(gdbServices, dbsvc)
				}
			}
			selectGDB := promptui.Select{
				Label: "select global db service",
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
			actions := []action{{"manage"}, {"backup"}}
			selectDBAction := promptui.Select{
				Label: "select action",
				Items: actions,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Name | cyan }}",
					Inactive: "  {{ .Name | cyan }}",
					Selected: fmt.Sprintf("\U0001F389 {{ .Name | green | cyan }} (%s)", gdbServices[it].Status.Address),
				},
			}
			iac, _, _ := selectDBAction.Run()
			if actions[iac].Name == "manage" {
				// https://console.dev.haogs.cn/adminer/?server=10.10.16.15%3A3306&username=root&db=ysicing&password=password123
				if err := fakeUserInfo(qclient, &gdbServices[it]); err != nil {
					return fmt.Errorf("call k8s api err: %v", err)
				}
				url := fmt.Sprintf("http://%s/adminer/?server=%s&username=%s&db=%s&password=%s", cfg.Domain, gdbServices[it].Status.Address, gdbServices[it].Spec.Account.User.Value, "", gdbServices[it].Spec.Account.Password.Value)
				if err := browser.OpenURL(url); err != nil {
					log.Warnf("try open browser err: %v", err)
					log.Infof("open browser access url: %s", url)
					return nil
				}
				log.Done("open browser")
				return nil
			} else if actions[iac].Name == "backup" {
				log.Warn("not implement")
			}
			return nil
		},
	}
	return app
}

func vaildGlobalDatabase(l map[string]string) bool {
	if exmap.CheckLabel(l, "easycorp.io/global_database") {
		return exmap.GetLabelValue(l, "easycorp.io/global_database") == "true"
	}
	return false
}

func fakeUserInfo(qclient *k8s.Client, dbsvc *quchengv1beta1.DbService) error {
	if dbsvc.Spec.Account.User.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), dbsvc.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: dbsvc.Spec.Account.User.ValueFrom.SecretKeyRef.Name,
			},
			Key: dbsvc.Spec.Account.User.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		dbsvc.Spec.Account.User.Value = string(user)
	}
	if dbsvc.Spec.Account.Password.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), dbsvc.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: dbsvc.Spec.Account.Password.ValueFrom.SecretKeyRef.Name,
			},
			Key: dbsvc.Spec.Account.Password.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		dbsvc.Spec.Account.Password.Value = string(user)
	}
	return nil
}
