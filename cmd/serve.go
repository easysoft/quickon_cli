// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/easysoft/qcadmin/internal/app/serve"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/spf13/cobra"
)

func newCmdServe() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve daemon",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
			go func() {
				<-ctx.Done()
				stop()
			}()

			if err := serve.Serve(ctx); err != nil {
				log.Flog.Fatalf("run serve: %v", err)
			}
		},
	}
	return cmd
}
