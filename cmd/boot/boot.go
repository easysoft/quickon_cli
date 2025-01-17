// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package boot

import (
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/zos"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/static"
)

var rootDirs = []string{
	common.DefaultLogDir,
	common.DefaultDataDir,
	common.DefaultBinDir,
	common.DefaultCfgDir,
	common.DefaultCacheDir,
}

var qDirs = []string{
	common.GetDefaultQuickonPlatformDir(""),
	common.GetDefaultQuickonBackupDir(""),
	common.DefaultNerdctlDir,
}

func initRootDirectory() error {
	home := zos.GetHomeDir()
	for _, dir := range rootDirs {
		if err := os.MkdirAll(home+"/"+dir, common.FileMode0755); err != nil {
			return errors.Errorf("failed to mkdir %s, err: %s", dir, err)
		}
	}
	for _, dir := range qDirs {
		if err := os.MkdirAll(common.GetCustomQuickonDir(dir), common.FileMode0777); err != nil {
			return errors.Errorf("failed to mkdir %s, err: %s", dir, err)
		}
	}

	// TODO 自定义目录可能有问题
	os.Chmod(common.GetDefaultQuickonBackupDir(""), common.FileMode0777)

	if err := static.StageFiles(); err != nil {
		return errors.Errorf("failed to stage files, err: %s", err)
	}
	if !file.CheckFileExists(common.DefaultNerdctlConfig) {
		file.Copy(common.GetCustomFile("hack/manifests/hub/nerdctl.toml"), common.DefaultNerdctlConfig, true)
	}
	return nil
}

func OnBoot() error {
	return initRootDirectory()
}
