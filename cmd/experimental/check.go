package experimental

import (
	"github.com/easysoft/qcadmin/internal/pkg/cli/check"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func CheckCommand(f factory.Factory) *cobra.Command {
	tCmd := &cobra.Command{
		Use:     "check",
		Short:   "check tools",
		Long:    "检查数据库等中间件是否可用",
		Version: "4.0.0",
	}
	tCmd.AddCommand(check.CheckMySQLCommand(f))
	return tCmd
}
