// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package manage

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/httptls"
)

func NewRenewTLS(f factory.Factory) *cobra.Command {
	var force bool
	tlsCmd := &cobra.Command{
		Use:     "tls",
		Short:   "check and renew tls",
		Version: "1.2.11",
		RunE: func(_ *cobra.Command, _ []string) error {
			return httptls.CheckReNewCertificate(force)
		},
	}
	tlsCmd.Flags().BoolVarP(&force, "force", "f", false, "force renew tls")
	return tlsCmd
}
