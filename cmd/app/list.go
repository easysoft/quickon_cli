// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/ergoapi/util/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func NewCmdAppList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "list",
		Short:   "list app",
		Example: `q app list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return err
			}
			release, _, err := hc.List(0, 0, "")
			if err != nil {
				return err
			}

			selectApp := promptui.Select{
				Label: "select app",
				Items: release,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Chart.Metadata.Name }})",
					Inactive: "  {{ .Name | cyan }}",
					Selected: "\U0001F389 {{ .Name | red | cyan }} ({{ .Chart.Metadata.Name }})",
				},
				Size: 5,
			}
			it, _, _ := selectApp.Run()
			log.Infof("select app: %s", release[it].Name)

			selectInfo := promptui.Select{
				Label: "select action: meta get url info, svc get container service message",
				Items: []string{"meta", "svc"},
			}
			_, infoAction, _ := selectInfo.Run()
			if infoAction == "meta" {
				values, err := hc.GetAllValues(release[it].Name)
				if err != nil {
					return err
				}
				host := getMapValue(getMap(getMap(values, "global"), "ingress"), "host")
				if len(host) != 0 {
					if kutil.IsLegalDomain(host) {
						host = fmt.Sprintf("https://%s", host)
					} else {
						host = fmt.Sprintf("http://%s", host)
					}
				}
				auth := getMap(values, "auth")
				if auth != nil {
					authUsername := getMapValue(auth, "username")
					authPassword := getMapValue(auth, "password")
					log.Debugf("authUsername: %s, authPassword: %s", authUsername, authPassword)
					log.Infof("app meta:\n\t   username: %s\n\t   password: %s\n\t   url: %s", color.SBlue(authUsername), color.SBlue(authPassword), color.SBlue(host))
				} else {
					log.Infof("app meta:\n\t url: %s", color.SBlue(host))
				}
				return nil
			}

			k8sClient, err := k8s.NewSimpleClient()
			if err != nil {
				log.Errorf("k8s client err: %v", err)
				return err
			}
			ctx := context.Background()
			podlist, _ := k8sClient.ListPods(ctx, release[it].Namespace, metav1.ListOptions{
				LabelSelector: labels.SelectorFromSet(map[string]string{
					"release": release[it].Name,
				}).String(),
			})
			if len(podlist.Items) < 1 {
				return errors.Errorf("podnum %d, app maybe not running", len(podlist.Items))
			}

			selectPod := promptui.Select{
				Label: "select pod",
				Items: podlist.Items,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Name | cyan }}",
					Inactive: "  {{ .Name | cyan }}",
					Selected: "\U0001F389 {{ .Name | red | cyan }}",
				},
				Size: 5,
			}
			podit, _, _ := selectPod.Run()
			log.Infof("select app %s pod %s", release[it].Name, podlist.Items[podit].Name)
			selectAction := promptui.Select{
				Label: "select action",
				Items: []string{"logs", "exec"},
			}
			_, action, _ := selectAction.Run()
			if action == "logs" {
				return k8sClient.GetFollowLogs(ctx, podlist.Items[podit].Namespace, podlist.Items[podit].Name, podlist.Items[podit].Spec.Containers[0].Name, false)
			}
			return k8sClient.ExecPodWithTTY(ctx, podlist.Items[podit].Namespace, podlist.Items[podit].Name, podlist.Items[podit].Spec.Containers[0].Name, []string{"/bin/sh", "-c", "sh"})
		},
	}
	return app
}
