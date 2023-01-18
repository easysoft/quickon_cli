package cluster

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

// ClusterCommand cluster command.
func ClusterCommand(f factory.Factory) *cobra.Command {
	c := &cobra.Command{
		Use:     "cluster",
		Short:   "cluster tools",
		Aliases: []string{"c"},
	}
	c.AddCommand(newInit(f))
	return c
}
