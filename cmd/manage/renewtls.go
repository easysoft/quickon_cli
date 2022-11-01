package manage

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/httptls"
	"github.com/spf13/cobra"
)

func NewRenewTLS(f factory.Factory) *cobra.Command {
	rtls := &cobra.Command{
		Use:     "renewtls",
		Short:   "renew qucheng tls domain",
		Aliases: []string{"rtls"},
		Version: "1.2.11",
		RunE: func(cmd *cobra.Command, args []string) error {
			return httptls.CheckReNewCertificate()
		},
	}
	return rtls
}
