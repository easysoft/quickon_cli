// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package boot

import (
	"os"

	"github.com/ergoapi/util/zos"
	"github.com/pkg/errors"

	"github.com/easysoft/qcadmin/common"
)

var rootDirs = []string{
	common.DefaultLogDir,
	common.DefaultDataDir,
	common.DefaultBinDir,
	common.DefaultCfgDir,
	common.DefaultCacheDir,
}

func initRootDirectory() error {
	home := zos.GetHomeDir()
	for _, dir := range rootDirs {
		err := os.MkdirAll(home+"/"+dir, common.FileMode0755)
		if err != nil {
			return errors.Errorf("failed to mkdir %s, err: %s", dir, err)
		}
	}
	if err := os.MkdirAll(common.GetDefaultQuickonDir(), common.FileMode0777); err != nil {
		return errors.Errorf("failed to mkdir %s, err: %s", common.GetDefaultQuickonDir(), err)
	}
	return nil
}

func OnBoot() error {
	return initRootDirectory()
}
