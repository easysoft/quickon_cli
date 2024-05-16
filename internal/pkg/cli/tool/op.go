// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package tool

import (
	"github.com/spf13/cobra"

	"github.com/easysoft/qcadmin/internal/pkg/util/downloader"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
)

func EmbedWgetCommand(f factory.Factory) *cobra.Command {
	var target, dst string
	wget := &cobra.Command{
		Use:   "wget",
		Short: "wget",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := downloader.Download(target, dst)
			return err
		},
	}
	wget.Flags().StringVarP(&target, "target", "t", "", "target url")
	wget.Flags().StringVarP(&dst, "dst", "d", "", "dst file")
	return wget
}
