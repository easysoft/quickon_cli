// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package market

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ergoapi/util/exnet"
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewMarketCmd create market command offline mode
func NewMarketCmd(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var helmhost string
	cmd := &cobra.Command{
		Use:    "market-init",
		Short:  "offline mode init market",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeClient, err := k8s.NewSimpleClient()
			if err != nil {
				log.Errorf("kube client create err: %v", err)
				return nil
			}
			log.Infof("start deploy cloudapp market")

			// install cne-market
			marketargs := []string{"experimental", "helm", "upgrade", "--name", "market", "--repo", common.DefaultHelmRepoName, "--chart", "cne-market-api", "--namespace", common.GetDefaultSystemNamespace(true)}
			output, err := qcexec.Command(os.Args[0], marketargs...).CombinedOutput()
			if err != nil {
				log.Warnf("upgrade install cloudapp market failed: %s", string(output))
			}
			// patch quickon
			cmfileName := fmt.Sprintf("%s-files", common.GetReleaseName(true))
			log.Debugf("fetch helm cm %s", cmfileName)
			for i := 0; i < 20; i++ {
				time.Sleep(5 * time.Second)
				foundRepofiles, _ := kubeClient.GetConfigMap(context.Background(), common.GetDefaultSystemNamespace(true), cmfileName, metav1.GetOptions{})
				if foundRepofiles != nil {
					foundRepofiles.Data["repositories.yaml"] = fmt.Sprintf(`apiVersion: ""
generated: "0001-01-01T00:00:00Z"
repositories:
- caFile: ""
  certFile: ""
  insecure_skip_tls_verify: true
  keyFile: ""
  name: qucheng-stable
  pass_credentials_all: false
  password: ""
  url: http://%s:32377
  username: ""
`, helmhost)
					_, err := kubeClient.UpdateConfigMap(context.Background(), foundRepofiles, metav1.UpdateOptions{})
					if err != nil {
						log.Warnf("patch offline repo file, check: kubectl get cm/%s  -n %s", cmfileName, common.GetDefaultSystemNamespace(true))
					}
					// 重建pod
					pods, _ := kubeClient.ListPods(context.Background(), common.GetDefaultSystemNamespace(true), metav1.ListOptions{})
					for _, pod := range pods.Items {
						if strings.HasPrefix(pod.Name, common.GetReleaseName(true)) {
							if err := kubeClient.DeletePod(context.Background(), pod.Name, common.GetDefaultSystemNamespace(true), metav1.DeleteOptions{}); err != nil {
								log.Warnf("recreate %s pods", common.GetReleaseName(true))
							}
						}
					}
					break
				}
			}
			genmarket := common.GetCustomFile("hack/manifests/scripts/market.sh")
			output, err = qcexec.Command(genmarket, helmhost).CombinedOutput()
			if err != nil {
				log.Errorf("generate market script failed: %s", string(output))
				return err
			}
			log.Done("run market script success")
			return nil
		},
	}
	cmd.Flags().StringVar(&helmhost, "host", exnet.LocalIPs()[0], "helm host")
	return cmd
}
