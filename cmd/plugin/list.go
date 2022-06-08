// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package plugin

import (
	"os"
	"strings"

	pluginapi "github.com/easysoft/qcadmin/internal/pkg/plugin"
	"github.com/easysoft/qcadmin/internal/pkg/util/output"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
)

var show string

func ListPluginCmd() *cobra.Command {

	listcmd := &cobra.Command{
		Use:     "list",
		Short:   "list",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ps, err := pluginapi.GetAll()
			if err != nil {
				return err
			}
			switch strings.ToLower(show) {
			case "json":
				return output.EncodeJSON(os.Stdout, ps)
			case "yaml":
				return output.EncodeYAML(os.Stdout, ps)
			default:
				table := uitable.New()
				table.MaxColWidth = 80
				table.Wrap = true
				for _, d := range ps {
					table.AddRow("type: ", d.Type)
					table.AddRow("default: ", d.Default)
					if len(d.Item) > 1 {
						str := []string{}
						for _, v := range d.Item {
							if v.Name != d.Default {
								str = append(str, v.Name)
							}
						}
						table.AddRow("optional: ", strings.Join(str, ", "))
					}
					table.AddRow("------------")
				}
				return output.EncodeTable(os.Stdout, table)
			}
		},
	}
	listcmd.Flags().StringVarP(&show, "output", "o", "", "prints the output in the specified format. Allowed values: table, json, yaml (default table)")
	return listcmd
}
