package apply

import (
	"github.com/easysoft/qcadmin/internal/pkg/cli/helm"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

// ApplyCommand apply command.
func ApplyCommand(f factory.Factory) *cobra.Command {
	return helm.EmbedCommand(f)
}
