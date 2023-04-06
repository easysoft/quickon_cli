// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package helm

import (
	"fmt"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/spf13/cobra"
)

func EmbedCommand(f factory.Factory) *cobra.Command {
	helm := &cobra.Command{
		Use:   "helm",
		Short: "The Kubernetes package manager",
	}
	helm.AddCommand(repoUpdate(f))
	helm.AddCommand(repoAdd(f))
	helm.AddCommand(repoDel(f))
	helm.AddCommand(chartUpgrade(f))
	helm.AddCommand(chartUninstall(f))
	helm.AddCommand(chartClean(f))
	helm.AddCommand(chartList(f))
	return helm
}

func repoUpdate(f factory.Factory) *cobra.Command {
	helm := &cobra.Command{
		Use:   "repo-update",
		Short: "update information of available charts locally from chart repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			hc.UpdateRepo()
			return nil
		},
	}
	return helm
}

func repoAdd(f factory.Factory) *cobra.Command {
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

func repoDel(f factory.Factory) *cobra.Command {
	var name string
	helm := &cobra.Command{
		Use:     "repo-del",
		Aliases: []string{"repo-remove"},
		Short:   "del a chart repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			if err := hc.RemoveRepo(name); err != nil {
				return err
			}
			return nil
		},
	}
	helm.Flags().StringVar(&name, "name", "install", "repo name")
	return helm
}

func chartUpgrade(f factory.Factory) *cobra.Command {
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

func chartUninstall(f factory.Factory) *cobra.Command {
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
			release, _ := hc.GetDetail(name)
			if release != nil {
				_, err = hc.Uninstall(name)
				return err
			}
			return nil
		},
	}
	helm.Flags().StringVarP(&ns, "namespace", "n", "", "namespace")
	helm.Flags().StringVar(&name, "name", "", "release name")
	return helm
}

// chartClean clean all chart installed by quickon
func chartClean(f factory.Factory) *cobra.Command {
	var ns, name string
	helm := &cobra.Command{
		Use:   "clean",
		Short: "clean all chart installed by quickon",
		RunE: func(cmd *cobra.Command, args []string) error {
			if ns == "" {
				ns = "default"
			}
			hc, err := helm.NewClient(&helm.Config{Namespace: ns})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			release, _ := hc.GetDetail(name)
			if release != nil {
				// _, err = hc.Clean(name)
				return err
			}
			return nil
		},
	}
	helm.Flags().StringVarP(&ns, "namespace", "n", "", "namespace")
	helm.Flags().StringVar(&name, "name", "", "release name")
	return helm
}

// chartList list all chart installed by quickon
func chartList(f factory.Factory) *cobra.Command {
	var ns string
	var page, limit int
	log := f.GetLog()
	helm := &cobra.Command{
		Use:   "list",
		Short: "list all chart installed by quickon",
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ns})
			if err != nil {
				return fmt.Errorf("helm create go client err: %v", err)
			}
			releases, _, err := hc.List(page, limit, "")
			if err != nil {
				return err
			}
			for _, release := range releases {
				log.Infof("name: %s, namespace: %s", release.Name, release.Namespace)
			}
			return nil
		},
	}
	helm.Flags().StringVarP(&ns, "namespace", "n", "", "namespace")
	helm.Flags().IntVar(&page, "page", 1, "page")
	helm.Flags().IntVar(&limit, "limit", 100, "limit")
	return helm
}
