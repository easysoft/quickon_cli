// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

//go:generate hack/scripts/getbin.sh
//go:generate go run internal/pkg/cli/codegen/codegen.go

package main

import (
	"github.com/docker/docker/pkg/reexec"
	"github.com/easysoft/qcadmin/cmd"
	"github.com/easysoft/qcadmin/cmd/boot"
	"github.com/easysoft/qcadmin/internal/pkg/cli/kubectl"
)

func init() {
	reexec.Register("kubectl", kubectl.Main)
}

func main() {
	if err := boot.OnBoot(); err != nil {
		panic(err)
	}
	cmd.Execute()
}
