package cluster

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

// q ext cluster init

func newInit(f factory.Factory) *cobra.Command {
	init := &cobra.Command{
		Use:   "init",
		Short: "init cluster",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return init
}
