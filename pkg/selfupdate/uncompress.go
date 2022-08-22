// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ulikunitz/xz"
)

func matchExecutableName(cmd, target string) bool {
	if cmd == target {
		return true
	}

	o, a := runtime.GOOS, runtime.GOARCH

	// When the contained executable name is full name (e.g. foo_darwin_amd64),
	// it is also regarded as a target executable file. (#19)
	for _, d := range []rune{'_', '-'} {
		c := fmt.Sprintf("%s%c%s%c%s", cmd, d, o, d, a)
		if o == "windows" {
			c += ".exe"
		}
		if c == target {
			return true
		}
	}

	return false
}

func unarchiveTar(log log.Logger, src io.Reader, url, cmd string) (io.Reader, error) {
	t := tar.NewReader(src)
	for {
		h, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to unarchive .tar file: %s", err)
		}
		_, name := filepath.Split(h.Name)
		if matchExecutableName(cmd, name) {
			log.Debugf("executable file %s was found in tar archive", h.Name)
			return t, nil
		}
	}

	return nil, fmt.Errorf("file '%s' for the command is not found in %s", cmd, url)
}

// UncompressCommand uncompresses the given source. Archive and compression format is
// automatically detected from 'url' parameter, which represents the URL of asset.
// This returns a reader for the uncompressed command given by 'cmd'. '.zip',
// '.tar.gz', '.tar.xz', '.tgz', '.gz' and '.xz' are supported.
func UncompressCommand(log log.Logger, src io.Reader, url, cmd string) (io.Reader, error) {
	if strings.HasSuffix(url, ".zip") {
		log.Debug("uncompressing zip file", url)

		// Zip format requires its file size for uncompressing.
		// So we need to read the HTTP response into a buffer at first.
		buf, err := io.ReadAll(src)
		if err != nil {
			return nil, fmt.Errorf("failed to create buffer for zip file: %s", err)
		}

		r := bytes.NewReader(buf)
		z, err := zip.NewReader(r, r.Size())
		if err != nil {
			return nil, fmt.Errorf("Failed to uncompress zip file: %s", err)
		}

		for _, file := range z.File {
			_, name := filepath.Split(file.Name)
			if !file.FileInfo().IsDir() && matchExecutableName(cmd, name) {
				log.Debug("executable file", file.Name, "was found in zip archive")
				return file.Open()
			}
		}

		return nil, fmt.Errorf("File '%s' for the command is not found in %s", cmd, url)
	} else if strings.HasSuffix(url, ".tar.gz") || strings.HasSuffix(url, ".tgz") {
		log.Debug("uncompressing tar.gz file", url)

		gz, err := gzip.NewReader(src)
		if err != nil {
			return nil, fmt.Errorf("Failed to uncompress .tar.gz file: %s", err)
		}

		return unarchiveTar(log, gz, url, cmd)
	} else if strings.HasSuffix(url, ".gzip") || strings.HasSuffix(url, ".gz") {
		log.Debug("uncompressing gzip file", url)

		r, err := gzip.NewReader(src)
		if err != nil {
			return nil, fmt.Errorf("Failed to uncompress gzip file downloaded from %s: %s", url, err)
		}

		name := r.Header.Name
		if !matchExecutableName(cmd, name) {
			return nil, fmt.Errorf("File name '%s' does not match to command '%s' found in %s", name, cmd, url)
		}

		log.Debug("executable file", name, "was found in gzip file")
		return r, nil
	} else if strings.HasSuffix(url, ".tar.xz") {
		log.Debug("uncompressing tar.xz file", url)

		xzip, err := xz.NewReader(src)
		if err != nil {
			return nil, fmt.Errorf("failed to uncompress .tar.xz file: %s", err)
		}

		return unarchiveTar(log, xzip, url, cmd)
	} else if strings.HasSuffix(url, ".xz") {
		log.Debug("uncompressing xzip file", url)

		xzip, err := xz.NewReader(src)
		if err != nil {
			return nil, fmt.Errorf("failed to uncompress xzip file downloaded from %s: %s", url, err)
		}

		log.Debug("uncompressed file from xzip is assumed to be an executable", cmd)
		return xzip, nil
	}

	log.Debug("uncompression is not needed", url)
	return src, nil
}
