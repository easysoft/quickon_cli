// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package gdb

import (
	"context"

	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/ergoapi/util/exmap"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewCmdGdbList(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	var name string
	app := &cobra.Command{
		Use:     "list",
		Short:   "list gdb",
		Example: `g gdb list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			qclient, err := k8s.NewSimpleQClient()
			if err != nil {
				return err
			}
			dbsvcs, err := qclient.ListQuchengDBSvc(context.TODO(), name, metav1.ListOptions{})
			if err != nil {
				return err
			}
			if len(dbsvcs.Items) == 0 {
				log.Warn("no found global database service")
				return nil
			}
			for _, dbsvc := range dbsvcs.Items {
				if vaildGlobalDatabase(dbsvc.Labels) {
					log.Infof("%s", dbsvc.Name)
				}
			}
			return nil
		},
	}
	return app
}

func vaildGlobalDatabase(l map[string]string) bool {
	if exmap.CheckLabel(l, "easycorp.io/global_database") {
		return exmap.GetLabelValue(l, "easycorp.io/global_database") == "true"
	}
	return false
}
