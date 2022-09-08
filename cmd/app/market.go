// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exhash"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdAppMarket(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	appMarket := &cobra.Command{
		Use:   "market",
		Short: "app market",
		Long:  "app market, you can choose app to install",
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := helm.NewClient(&helm.Config{Namespace: ""})
			if err != nil {
				return err
			}
			repos, err := hc.ListRepo()
			if err != nil {
				return err
			}
			if len(repos) == 0 {
				return fmt.Errorf("not found qucheng market, you should: %s exp helm repo-add --name qucheng --url https://hub.qucheng.com/chartrepo/stable", os.Args[0])
			}
			quchengRepoName := ""
			for _, repo := range repos {
				if strings.Contains(repo.URL, "qucheng") {
					quchengRepoName = repo.Name
					break
				}
			}
			if len(quchengRepoName) == 0 {
				return fmt.Errorf("not found qucheng market, you should: %s exp helm repo-add --name qucheng --url https://hub.qucheng.com/chartrepo/stable", os.Args[0])
			}
			charts, err := hc.ListCharts(quchengRepoName, "", false)
			if err != nil {
				return fmt.Errorf("fetch remote market failed, reason: %v", err)
			}
			selectInstallApp := promptui.Select{
				Label: "select install app",
				Items: charts,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . }}?",
					Active:   "\U0001F449 {{ .Chart.Metadata.Name | cyan }} ({{ .Chart.AppVersion }}) {{ .Chart.Metadata.Description }}",
					Inactive: "  {{ .Chart.Metadata.Name | cyan }}",
					Selected: "\U0001F389 {{ .Chart.Metadata.Name | red | cyan }} ({{ .Chart.AppVersion }})",
				},
				Size: 5,
			}
			it, _, _ := selectInstallApp.Run()
			//  spew.Dump(charts[it])
			log.Infof("select install app: %s, version: %s", color.SGreen(charts[it].Chart.Name), color.SGreen(charts[it].Chart.AppVersion))
			subdomain := fmt.Sprintf("%s-%s", charts[it].Chart.Name, exhash.MD5(time.Now().String())[:8])
			installArgs := []string{"app", "install", charts[it].Chart.Name, "--domain", subdomain, "--api-useip"}
			if log.GetLevel() == logrus.DebugLevel {
				installArgs = append(installArgs, "--debug")
			}
			if err := qcexec.CommandRun(os.Args[0], installArgs...); err != nil {
				log.Errorf("install default app %s err, reason: %v", charts[it].Chart.Name, err)
				return err
			}
			return nil
		},
	}
	return appMarket
}
