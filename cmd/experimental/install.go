// AGPL License
// Copyright (c) 2021 ysicing <i@ysicing.me>

package experimental

import (
	"fmt"
	"os"
	"runtime"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/color"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/downloader"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

var (
	installExample = templates.Examples(`
		# install tools
		z experimental install helm`)
	installTools = map[string]any{
		"helm":    true,
		"kubectl": true,
		"mc":      true,
		"etcdctl": true,
		"dnsctl":  true,
		"nerdctl": true,
	}
)

// InstallCommand install some tools
func InstallCommand(f factory.Factory) *cobra.Command {
	installCmd := &cobra.Command{
		Use:     "install [flags]",
		Short:   "install tools, like: helm, kubectl,etcdctl,mc,dnsctl,nerdctl",
		Example: installExample,
		Args:    cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				f.GetLog().Fatalf("args error: %v", args)
				return errors.New("missing args: tool name")
			}
			tool := args[0]
			if _, exist := installTools[tool]; !exist {
				return errors.Errorf("not support tool: %s", tool)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			f.GetLog().Debugf("cli args: %v", args)
			tool := args[0]
			remoteURL := fmt.Sprintf("https://pkg.qucheng.com/qucheng/cli/stable/tools/%s-%s-%s", tool, runtime.GOOS, runtime.GOARCH)
			localURL := fmt.Sprintf("%s/qc-%s", common.GetDefaultBinDir(), tool)
			res, err := downloader.Download(remoteURL, localURL)
			if err != nil {
				f.GetLog().Fatalf("download %s error: %v", tool, err)
				return
			}
			f.GetLog().Debugf("download %s result: %v", tool, res.Status)
			_ = os.Chmod(localURL, common.FileMode0755)
			f.GetLog().Donef(fmt.Sprintf("download %s success\n\t usage:   %s", tool, color.SGreen("%s %s", os.Args[0], tool)))
		},
	}
	return installCmd
}
