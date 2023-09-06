// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package manage

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/httptls"
	"github.com/spf13/cobra"
)

func NewRenewTLS(f factory.Factory) *cobra.Command {
	var force bool
	rtls := &cobra.Command{
		Use:     "renewtls",
		Short:   "renew tls domain",
		Aliases: []string{"rtls", "rt"},
		Version: "1.2.11",
		RunE: func(cmd *cobra.Command, args []string) error {
			return httptls.CheckReNewCertificate(force)
		},
	}
	rtls.Flags().BoolVarP(&force, "force", "f", false, "force renew tls")
	return rtls
}
