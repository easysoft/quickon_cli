// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package kubectl

import (
	"github.com/spf13/cobra"
	"k8s.io/component-base/cli"
	"k8s.io/kubectl/pkg/cmd"
	"k8s.io/kubectl/pkg/cmd/util"
)

func Main() {
	main()
}

func main() {
	command := cmd.NewDefaultKubectlCommand()
	if err := cli.RunNoErrOutput(command); err != nil {
		util.CheckErr(err)
	}
}

// EmbedCommand Used to embed the kubectl command.
func EmbedCommand() *cobra.Command {
	c := cmd.NewDefaultKubectlCommand()
	c.Short = "Kubectl controls the Kubernetes cluster manager"
	return c
}
