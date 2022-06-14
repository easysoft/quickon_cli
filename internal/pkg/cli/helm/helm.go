// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package helm

import (
	"fmt"

	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/spf13/cobra"
)

func EmbedCommand() *cobra.Command {
	helm := &cobra.Command{
		Use:   "helm",
		Short: "The Kubernetes package manager",
	}
	helm.AddCommand(repoUpdate())
	helm.AddCommand(repoAdd())
	return helm
}

func repoUpdate() *cobra.Command {
	helm := &cobra.Command{
		Use:   "repo-update",
		Short: "update information of available charts locally from chart repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			return hc.UpdateRepo()
		},
	}
	return helm
}

func repoAdd() *cobra.Command {
	var name, url, username, password string
	helm := &cobra.Command{
		Use:   "repo-add",
		Short: "add a chart repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			return hc.AddRepo(name, url, username, password)
		},
	}
	helm.Flags().StringVar(&name, "name", "", "")
	return helm
}
