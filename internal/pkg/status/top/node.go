// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package top

import (
	"fmt"
	"os"
	"strings"

	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/output"
	"github.com/gosuri/uitable"
	"k8s.io/apimachinery/pkg/labels"
)

type NodeOption struct {
	Selector   string
	SortBy     string
	QPS        float32
	Burst      int
	KubeCtx    string
	KubeConfig string
	Output     string
}

func (o *NodeOption) Validate() {
	if len(o.SortBy) > 0 {
		if o.SortBy != "cpu" {
			o.SortBy = "memory"
		}
	}
}

func (o *NodeOption) RunResourceNode() error {
	selector := labels.Everything()
	var err error
	if len(o.Selector) > 0 {
		selector, err = labels.Parse(o.Selector)
		if err != nil {
			return err
		}
	}
	k, err := k8s.NewClient(o.KubeCtx, o.KubeConfig)
	if err != nil {
		return err
	}
	data, err := k.GetNodeResources(o.SortBy, selector)
	if err != nil {
		return err
	}
	switch strings.ToLower(o.Output) {
	case "json":
		return output.EncodeJSON(os.Stdout, data)
	case "yaml":
		return output.EncodeYAML(os.Stdout, data)
	default:
		table := uitable.New()
		table.AddRow("Name", "IP", "CPU Usage", "CPU Request", "CPU Limit", "CPU Capacity", "Memory Usage", "Memory Request", "Memory Limit", "Memory Capacity", "Pod Usage", "Pod Capacity", "Age")
		for _, d := range data {
			table.AddRow(d.NodeName, d.NodeIP,
				d.CPUUsages, fmt.Sprintf("%v(%v)", d.CPURequests, d.CPURequestsFraction), fmt.Sprintf("%v(%v)", d.CPULimits, d.CPULimitsFraction), d.CPUCapacity,
				d.MemoryUsages, fmt.Sprintf("%v(%v)", d.MemoryRequests, d.MemoryRequestsFraction), fmt.Sprintf("%v(%v)", d.MemoryLimits, d.MemoryLimitsFraction), d.MemoryCapacity,
				fmt.Sprintf("%v(%v)", d.AllocatedPods, d.PodFraction), d.PodCapacity, d.Age)
		}
		return output.EncodeTable(os.Stdout, table)
	}
}
