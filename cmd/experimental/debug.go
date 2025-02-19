package experimental

import (
	"github.com/easysoft/qcadmin/internal/pkg/cli/debug"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

// DebugCommand debug command.
func DebugCommand(f factory.Factory) *cobra.Command {
	debugCmd := &cobra.Command{
		Use:   "debug",
		Short: "debug, not a stable interface, contains misc debug facilities",
	}
	debugCmd.AddCommand(debug.GOPSCommand(f))
	debugCmd.AddCommand(debug.HostInfoCommand(f))
	debugCmd.AddCommand(debug.NetcheckCommand(f))
	debugCmd.AddCommand(debug.PortForwardCommand(f))
	return debugCmd
}
