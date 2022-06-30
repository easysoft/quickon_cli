// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package binfile

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/easysoft/qcadmin/common"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/ergoapi/util/file"
)

type Meta struct{}

// func (p *Meta) unpack(binName string) error {
// 	log.Debugf("unpacking %s-linux-%s", binName, runtime.GOARCH)
// 	sourcefile, err := bin.BinFS.ReadFile(fmt.Sprintf("%s-linux-%s", binName, runtime.GOARCH))
// 	if err != nil {
// 		return err
// 	}
// 	installFile, err := os.OpenFile(fmt.Sprintf("/usr/local/bin/%s", binName), os.O_CREATE|os.O_RDWR|os.O_TRUNC, common.FileMode0755)
// 	defer func() { _ = installFile.Close() }()
// 	if err != nil {
// 		return err
// 	}
// 	if _, err := io.Copy(installFile, bytes.NewReader(sourcefile)); err != nil {
// 		return err
// 	}
// 	log.Donef("unpack %s complete", binName)
// 	return nil
// }

func (p *Meta) LoadLocalBin(binName string) (string, error) {
	filebin, err := exec.LookPath(binName)
	if err != nil {
		sourcebin := fmt.Sprintf("%s/hack/bin/%s-%s-%s", common.GetDefaultDataDir(), binName, runtime.GOOS, runtime.GOARCH)
		filebin = fmt.Sprintf("/usr/local/bin/%s", binName)
		if file.CheckFileExists(sourcebin) {
			if err := qcexec.Command("cp", "-a", sourcebin, filebin).Run(); err != nil {
				return "", err
			}
		}
	}
	output, err := qcexec.Command(filebin, "--help").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("seems like there are issues with your %s client: \n\n%s", binName, output)
	}
	return filebin, nil
}

// func (p *Meta) download(binName string) error {
// 	log.Debugf("unpack %s bin failed, will download from remote.", binName)
// 	binPath := fmt.Sprintf("/usr/local/bin/%s", binName)
// 	if _, err := downloader.Download(common.GetBinURL(binName), binPath); err != nil {
// 		return err
// 	}
// 	os.Chmod(binPath, common.FileMode0755)
// 	log.Donef("download %s complete", binName)
// 	return nil
// }
