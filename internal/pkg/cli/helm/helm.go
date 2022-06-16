// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package helm

import (
	"fmt"

	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/spf13/cobra"
)

func EmbedCommand() *cobra.Command {
	helm := &cobra.Command{
		Use:   "helm",
		Short: "The Kubernetes package manager",
	}
	helm.AddCommand(repoUpdate())
	helm.AddCommand(repoAdd())
	helm.AddCommand(chartUpgrade())
	helm.AddCommand(chartUninstall())
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
			hc.UpdateRepo()
			log.Flog.Done("Update Complete. ⎈ Happy Helming!⎈ ")
			return nil
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
			if len(name) == 0 || len(url) == 0 {
				return fmt.Errorf("name or url is empty")
			}
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			return hc.AddRepo(name, url, username, password)
		},
	}
	helm.Flags().StringVar(&name, "name", "", "repo name")
	helm.Flags().StringVar(&url, "url", "", "repo url")
	helm.Flags().StringVar(&username, "username", "", "private repo username")
	helm.Flags().StringVar(&password, "password", "", "private repo password")
	return helm
}

func chartUpgrade() *cobra.Command {
	var ns, name, repoName, chartName, chartVersion string
	var p []string
	helm := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade a release",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(name) == 0 || len(repoName) == 0 || len(chartName) == 0 {
				return fmt.Errorf("name or repoName or chartName is empty")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if ns == "" {
				ns = "default"
			}
			hc, err := helm.NewClient(&helm.Config{Namespace: ns})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			values, _ := helm.MergeValues(p)
			_, err = hc.Upgrade(name, repoName, chartName, chartVersion, values)
			return err
		},
	}
	helm.Flags().StringVarP(&ns, "namespace", "n", "", "namespace")
	helm.Flags().StringVar(&name, "name", "", "release name")
	helm.Flags().StringVar(&repoName, "repo", "", "repo name")
	helm.Flags().StringVar(&chartName, "chart", "", "chart name")
	helm.Flags().StringVar(&chartVersion, "version", "", "chart version")
	helm.Flags().StringArrayVar(&p, "set", []string{}, "set values on the command line (e.g. '--set key1=value1,key2=value2')")
	return helm
}

func chartUninstall() *cobra.Command {
	var ns, name string
	helm := &cobra.Command{
		Use:     "uninstall",
		Aliases: []string{"un", "del", "delete"},
		Short:   "uninstall a chart",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(name) == 0 {
				return fmt.Errorf("name is empty")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if ns == "" {
				ns = "default"
			}
			hc, err := helm.NewClient(&helm.Config{Namespace: ns})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			_, err = hc.Uninstall(name)
			return err
		},
	}
	helm.Flags().StringVarP(&ns, "namespace", "n", "", "namespace")
	helm.Flags().StringVar(&name, "name", "", "release name")
	return helm
}
