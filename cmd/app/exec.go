// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package app

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/internal/app/debug"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func NewCmdAppExec(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var useip bool
	app := &cobra.Command{
		Use:     "exec",
		Short:   "exec app",
		Args:    cobra.ExactArgs(1),
		Example: `q app exec http://console.example.corp.cc/instance-view-39.html`,
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			apidebug := log.GetLevel() == logrus.DebugLevel
			log.Infof("start exec app: %s", url)
			appdata, err := debug.GetNameByURL(url, apidebug, useip)
			if err != nil {
				return err
			}
			k8sClient, err := k8s.NewSimpleClient()
			if err != nil {
				log.Errorf("k8s client err: %v", err)
				return err
			}
			ctx := context.Background()
			podlist, _ := k8sClient.ListPods(ctx, "default", metav1.ListOptions{
				LabelSelector: labels.SelectorFromSet(map[string]string{
					"release": appdata.K8Name,
				}).String(),
			})
			if len(podlist.Items) < 1 {
				return errors.Errorf("podnum %d, app maybe not running", len(podlist.Items))
			}
			templates := &promptui.SelectTemplates{
				Label:    "{{ . }}?",
				Active:   "\U0001F449 {{ .Name | cyan }}",
				Inactive: "  {{ .Name | cyan }}",
				Selected: "\U0001F389 {{ .Name | red | cyan }}",
			}

			prompt := promptui.Select{
				Label:     "select pod",
				Items:     podlist.Items,
				Templates: templates,
				Size:      5,
			}
			it, _, _ := prompt.Run()

			return k8sClient.ExecPodWithTTY(ctx, "default", podlist.Items[it].Name, podlist.Items[it].Spec.Containers[0].Name, []string{"/bin/sh", "-c", "sh"})
		},
	}
	app.Flags().BoolVar(&useip, "api-useip", false, "api use ip")
	return app
}
