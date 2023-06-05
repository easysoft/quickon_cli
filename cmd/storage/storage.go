// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package storage

import (
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log/survey"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	nfsExample = templates.Examples(`
		# deploy local nfs storage
		q cluster storage nfs
		# deploy qcloud cfs v3
		q cluster storage nfs --ip cfsip --path cfspath
`)
)

// NewCmdStorage returns a cobra command for `storage` subcommands
func NewCmdStorage(f factory.Factory) *cobra.Command {
	s := &cobra.Command{
		Use:   "storage",
		Short: "storage commands",
		Long:  "install cluster storage",
	}
	s.AddCommand(nfs(f))
	return s
}

func nfs(f factory.Factory) *cobra.Command {
	var ip, path string
	logpkg := f.GetLog()
	cmd := &cobra.Command{
		Use:     "nfs",
		Short:   "deploy nfs storage",
		Example: nfsExample,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(ip) == 0 {
				an, err := logpkg.Question(&survey.QuestionOptions{
					Question:     "nfs server ip is empty, install local nfs",
					DefaultValue: "yes",
					Options:      []string{"yes", "no"},
				})
				if err != nil {
					return err
				}
				logpkg.Infof("answer: %s", an)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	cmd.Flags().StringVar(&ip, "ip", "", "cloud cfs/nas ip")
	cmd.Flags().StringVar(&path, "path", "", "cloud cfs/nas path")
	return cmd
}
