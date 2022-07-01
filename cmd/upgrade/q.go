// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/easysoft/qcadmin/cmd/version"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/selfupdate"
	"github.com/spf13/cobra"
)

type option struct {
	log log.Logger
}

func NewUpgradeQ(f factory.Factory) *cobra.Command {
	up := option{
		log: f.GetLog(),
	}
	upq := &cobra.Command{
		Use:     "q",
		Aliases: []string{"qcadmin"},
		Short:   "upgrade qcadmin(q) to the newest version",
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			up.DoQcadmin()
		},
	}
	return upq
}

func (up option) DoQcadmin() {
	up.log.StartWait("fetch latest version from remote...")
	lastversion, err := version.PreCheckLatestVersion()
	up.log.StopWait()
	if err != nil {
		up.log.Errorf("fetch latest version err, reason: %v", err)
		return
	}
	if lastversion == "" || lastversion == common.Version || strings.Contains(common.Version, lastversion) {
		up.log.Infof("The current version %s is the latest version", common.Version)
		return
	}
	cmdPath, err := os.Executable()
	if err != nil {
		up.log.Errorf("q executable err:%v", err)
		return
	}
	up.log.StartWait(fmt.Sprintf("downloading version %s...", lastversion))
	assetURL := fmt.Sprintf("https://pkg.qucheng.com/qucheng/cli/stable/qcadmin_%s_%s", runtime.GOOS, runtime.GOARCH)
	err = selfupdate.UpdateTo(up.log, assetURL, cmdPath)
	up.log.StopWait()
	if err != nil {
		up.log.Errorf("upgrade failed, reason: %v", err)
		return
	}
	up.log.Donef("Successfully updated ergo to version %s", lastversion)
	up.log.Infof("Release note: \n\trelease %s ", lastversion)
}
