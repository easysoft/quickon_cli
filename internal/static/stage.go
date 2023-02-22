// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package static

import (
	"fmt"
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/static/data"
	"github.com/easysoft/qcadmin/internal/static/deploy"
	"github.com/easysoft/qcadmin/internal/static/haogstls"
	"github.com/easysoft/qcadmin/internal/static/scripts"
	"github.com/ergoapi/util/file"
)

func StageFiles() error {
	dataDir := common.GetDefaultDataDir()
	if err := data.Stage(dataDir); err != nil {
		return err
	}
	if err := deploy.Stage(dataDir); err != nil {
		return err
	}
	if err := scripts.Stage(dataDir); err != nil {
		return err
	}
	if err := haogstls.Stage(dataDir); err != nil {
		return err
	}
	if err := initInternalCommand(dataDir); err != nil {
		return err
	}
	return nil
}

func initInternalCommand(dataDir string) error {
	// cp -a /root/.qc/data/hack/manifests/scripts/qc-* /root/.qc/bin/
	// cp -a /root/.qc/data/hack/manifests/scripts/qcadmin-* /root/.qc/bin/
	sourcePath := fmt.Sprintf("%s/hack/manifests/scripts", dataDir)
	files, err := file.DirFilesList(sourcePath, common.ValidPrefixes, nil)
	if err != nil {
		return err
	}
	for _, f := range files {
		s := strings.Split(f, "/")
		targetfile := fmt.Sprintf("%s/%s", common.GetDefaultBinDir(), s[len(s)-1])
		sourcefile := fmt.Sprintf("%s/%s", sourcePath, f)
		if err := file.Copy(sourcefile, targetfile, true); err != nil {
			return err
		}
	}
	return nil
}
