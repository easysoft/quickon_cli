// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cmd

import (
	"os"

	"github.com/easysoft/qcadmin/cmd/flags"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/providers"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

var (
	joinCmd = &cobra.Command{
		Use:   "join",
		Short: "Join node(s) to an existing QuCheng cluster",
	}
	jp providers.Provider
)

func init() {
	joinCmd.PersistentFlags().BoolVar(&skip, "skip-precheck", false, "skip precheck")
}

func newCmdJoin(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	// TODO 探测
	name := "native"
	if reg, err := providers.GetProvider(name); err != nil {
		log.Fatalf("failed to get provider: %s", err)
	} else {
		jp = reg
	}
	joinCmd.Flags().AddFlagSet(flags.ConvertFlags(joinCmd, jp.GetJoinFlags()))
	joinCmd.Example = jp.GetUsageExample("join")
	joinCmd.Run = func(cmd *cobra.Command, args []string) {
		if err := jp.PreSystemInit(); err != nil {
			log.Fatalf("presystem init err, reason: %s", err)
		}
		if err := jp.CreateCheck(skip); err != nil {
			log.Fatalf("precheck err, reason: %v", err)
		}
		if err := jp.JoinNode(); err != nil {
			log.Fatal(err)
		}
	}
	joinCmd.AddCommand(newCmdGenJoin(f))
	return joinCmd
}

func newCmdGenJoin(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	genjoin := &cobra.Command{
		Use:   "gen",
		Short: "Generate a join command",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig()
			if cfg == nil {
				log.Fatalf("only support run firstnode")
				return
			}
			log.Infof("\tjoin command: %s join --cne-api %s --cne-token %s", os.Args[0], cfg.InitNode, cfg.Token)
		},
	}
	return genjoin
}
