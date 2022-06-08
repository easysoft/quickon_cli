// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package selfupdate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/inconshreveable/go-update"
)

type Updater struct{}

func DefaultUpdater() *Updater {
	return &Updater{}
}

func UpdateTo(assetURL, cmdPath string) error {
	up := DefaultUpdater()
	src, err := up.downloadDirectlyFromURL(assetURL)
	if err != nil {
		return err
	}
	defer src.Close()
	return uncompressAndUpdate(src, assetURL, cmdPath)
}

func (up *Updater) downloadDirectlyFromURL(assetURL string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", assetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create HTTP request to %s: %s", assetURL, err)
	}

	req.Header.Add("Accept", "application/octet-stream")
	req = req.WithContext(context.Background())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to download a release file from %s: %s", assetURL, err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to download a release file from %s: Not successful status %d", assetURL, res.StatusCode)
	}

	return res.Body, nil
}

func uncompressAndUpdate(src io.Reader, assetURL, cmdPath string) error {
	_, cmd := filepath.Split(cmdPath)
	asset, err := UncompressCommand(src, assetURL, cmd)
	if err != nil {
		return err
	}

	log.Flog.Debugf("Will update %s to the latest downloaded from", cmdPath, assetURL)
	return update.Apply(asset, update.Options{
		TargetPath: cmdPath,
	})
}
