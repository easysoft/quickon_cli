// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package debug

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func IngressNoHostCommand(f factory.Factory) *cobra.Command {
	ingress := &cobra.Command{
		Use:   "ingress-no-host",
		Short: "ingress not listen on host",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO 需要二次确认
		},
		Hidden: true,
	}
	return ingress
}
