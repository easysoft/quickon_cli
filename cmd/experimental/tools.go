package experimental

import (
	"github.com/easysoft/qcadmin/internal/pkg/cli/tool"
	"github.com/spf13/cobra"
)

// ToolsCommand helm command.
func ToolsCommand() *cobra.Command {
	return tool.EmbedCommand()
}
