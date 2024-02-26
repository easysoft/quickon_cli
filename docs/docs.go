// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"strings"

	"github.com/easysoft/qcadmin/cmd"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/github"
	"github.com/ergoapi/util/version"
	doc "github.com/ysicing/cobra2vitepress"
)

type versionInfo struct {
	Latest string `json:"latest"`
	Stable string `json:"stable"`
	Dev    string `json:"dev"`
}

func main() {
	f := factory.DefaultFactory()
	q := cmd.BuildRoot(f)
	err := doc.GenMarkdownTree(q, "./docs")
	if err != nil {
		panic(err)
	}
	pkg := github.Pkg{
		Owner: "easysoft",
		Repo:  "quickon_cli",
	}
	tag, err := pkg.LastTag()
	if err != nil {
		return
	}
	nextVersion := strings.TrimPrefix(version.Next(tag.Name, false, false, true), "v")
	file.WriteFile("VERSION", nextVersion, true)
	v := versionInfo{
		Latest: nextVersion,
		Stable: nextVersion,
		Dev:    nextVersion,
	}
	jsonData, _ := json.MarshalIndent(v, "", "    ")
	file.WriteFile("version.json", string(jsonData), true)
}
