// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package version

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	gv "github.com/blang/semver/v4"
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/color"
	"github.com/imroc/req/v3"
)

var versionTpl = `q version:
 Version:           %v
 Go version:        %v
 Git commit:        %v
 Built:             %v
 OS/Arch:           %v
 Experimental:      true
`

const (
	defaultVersion       = "0.0.0"
	defaultGitCommitHash = "a1b2c3d4"
	defaultBuildDate     = "Mon Aug  3 15:06:50 2020"
)

type versionGet struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
}

// PreCheckLatestVersion 检查最新版本
func PreCheckLatestVersion() (string, error) {
	lastVersion := &versionGet{}
	client := req.C().SetUserAgent(common.GetUG()).SetTimeout(time.Second * 5)
	_, err := client.R().SetResult(lastVersion).Get("http://release.metrics.qucheng.com/last/qcadmin?type=none")
	if err != nil {
		return "", err
	}
	return lastVersion.Version, nil
}

func ShowVersion() {
	// logo.PrintLogo()
	if common.Version == "" {
		common.Version = defaultVersion
	}
	if common.BuildDate == "" {
		common.BuildDate = defaultBuildDate
	}
	if common.GitCommitHash == "" {
		common.GitCommitHash = defaultGitCommitHash
	}
	osarch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	fmt.Printf(versionTpl, common.Version, runtime.Version(), common.GitCommitHash, common.BuildDate, osarch)
	log.Flog.StartWait("check update...")
	lastversion, err := PreCheckLatestVersion()
	log.Flog.StopWait()
	if err != nil {
		log.Flog.Debugf("get update message err: %v", err)
		return
	}
	if lastversion != "" && !strings.Contains(lastversion, defaultVersion) {
		nowversion, _ := gv.New(common.Version)
		needupgrade := nowversion.LT(gv.MustParse(lastversion))
		if needupgrade {
			log.Flog.Infof("当前最新版本 %s, 可以使用 %s 将版本升级到最新版本", color.SGreen(lastversion), color.SGreen("ergo upgrade"))
			return
		}
	}
	log.Flog.Info("current version is the latest")
}
