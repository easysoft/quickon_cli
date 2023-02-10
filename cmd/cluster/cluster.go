package cluster

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func InitCommand(f factory.Factory) *cobra.Command {
	init := &cobra.Command{
		Use:   "init",
		Short: "init cluster",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return init
}

func JoinCommand(f factory.Factory) *cobra.Command {
	join := &cobra.Command{
		Use:   "join",
		Short: "join cluster",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return join
}

func CleanCommand(f factory.Factory) *cobra.Command {
	clean := &cobra.Command{
		Use:   "clean",
		Short: "clean cluster",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return clean
}
