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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewCmdAppList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "list",
		Short:   "list app",
		Example: `z app list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return err
			}
			release, _, err := hc.List(0, 0, "")
			if err != nil {
				return err
			}

			if len(release) == 0 {
				log.Warn("no found app")
				return nil
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
			podName := podlist.Items[podit].Name
			log.Infof("select app %s pod %s", release[it].Name, podName)
			selectAction := promptui.Select{
				Label: "select action",
				Items: []string{"logs", "exec"},
			}
			_, action, _ := selectAction.Run()
			podNamespace := podlist.Items[podit].Namespace
			selectPodContainer := promptui.Select{
				Label: fmt.Sprintf("select %s's container", podName),
				Items: podlist.Items[podit].Spec.Containers,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Name | cyan }}",
					Inactive: "  {{ .Name | cyan }}",
					Selected: "\U0001F389 {{ .Name | red | cyan }}",
				},
				Size: 5,
			}
			podContainerit, _, _ := selectPodContainer.Run()
			containerName := podlist.Items[podit].Spec.Containers[podContainerit].Name
			log.Infof("select app %s pod %s's container %s", release[it].Name, podName, containerName)
			if action == "logs" {
				return k8sClient.GetFollowLogs(ctx, podNamespace, podName, containerName, false)
			}
			return k8sClient.ExecPodWithTTY(ctx, podNamespace, podName, containerName, []string{"/bin/sh", "-c", "sh"})
		},
	}
	return app
}
