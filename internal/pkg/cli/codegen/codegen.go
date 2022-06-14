// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/BeidouCloudPlatform/go-bindata/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	os.Unsetenv("GOPATH")
	bc := &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "hack/bin",
				Recursive: true,
			},
		},
		Package:    "data",
		Arch:       "arm64",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/data/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "hack/bin",
				Recursive: true,
			},
		},
		Package:    "data",
		Arch:       "amd64",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/data/zz_generated_bindata.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "hack/manifests/scripts",
				Recursive: true,
			},
		},
		Package:    "scripts",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/scripts/zz_generated_bindata.go.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
	bc = &bindata.Config{
		Input: []bindata.InputConfig{
			{
				Path:      "hack/manifests/plugins",
				Recursive: true,
			},
		},
		Package:    "deploy",
		NoCompress: true,
		NoMemCopy:  true,
		NoMetadata: true,
		Output:     "internal/static/deploy/zz_generated_bindata.go.go",
	}
	if err := bindata.Translate(bc); err != nil {
		logrus.Fatal(err)
	}
}
