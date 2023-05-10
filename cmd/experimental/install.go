// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package experimental

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	installExample = templates.Examples(`
		# install tools
		q experimental install helm`)
)

// InstallCommand install some tools
func InstallCommand(f factory.Factory) *cobra.Command {
	installCmd := &cobra.Command{
		Use:     "install [flags]",
		Short:   "install tools",
		Example: installExample,
	}
	return installCmd
}
