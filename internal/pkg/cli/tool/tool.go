// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package tool

import (
	"context"
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/qucheng/domain"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func EmbedCommand() *cobra.Command {
	dns := &cobra.Command{
		Use:    "dns",
		Short:  "dns manager",
		Hidden: true,
	}
	dns.AddCommand(dnsClean())
	dns.AddCommand(dnsAdd())
	return dns
}

func dnsClean() *cobra.Command {
	dns := &cobra.Command{
		Use:    "clean",
		Short:  "clean dns",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig()
			if cfg != nil {
				if !strings.HasSuffix(cfg.Domain, "haogs.cn") {
					return
				}
			}
			kclient, _ := k8s.NewSimpleClient()
			cm, err := kclient.Clientset.CoreV1().ConfigMaps(common.DefaultSystem).Get(context.TODO(), "q-suffix-host", metav1.GetOptions{})
			if err != nil {
				return
			}
			reqbody := domain.ReqBody{
				UUID:      cm.Data["uuid"],
				SecretKey: cm.Data["auth"],
			}
			client := req.C().SetUserAgent(common.GetUG())
			if _, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(&reqbody).
				Delete(common.GetAPI("/api/qdns/oss/record")); err != nil {
				log.Flog.Error("clean dns failed, reason: %v", err)
			}
		},
	}
	return dns
}

func dnsAdd() *cobra.Command {
	dns := &cobra.Command{
		Use:    "init",
		Short:  "init dns",
		Hidden: true,
		Run:    func(cmd *cobra.Command, args []string) {},
	}
	return dns
}
