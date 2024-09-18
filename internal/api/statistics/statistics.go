// Copyright (c) 2021-2024 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package statistics

import (
	"fmt"

	"github.com/imroc/req/v3"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
)

type CollectData struct {
	CliVersion string `json:"cli_version"`
	ClusterID  string `json:"cluster_id"`
	Domain     string `json:"domain,omitempty"`
	Type       string `json:"type"`
	Action     string `json:"action"`
	Devops     bool   `json:"devops,omitempty"`
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendStatistics(action string) error {
	cfg, _ := config.LoadConfig()

	data := &CollectData{
		CliVersion: common.Version,
		ClusterID:  cfg.Cluster.ID,
		Domain:     cfg.Domain,
		Type:       cfg.Quickon.Type.String(),
		Action:     action,
		Devops:     cfg.Quickon.DevOps,
	}

	// send statistics
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG())
	var result Result
	resp, err := client.R().
		SetHeader("accept", "application/json").
		SetHeader("cluster.id", cfg.Cluster.ID).
		SetBody(data).
		SetSuccessResult(&result).
		Post("https://api.qucheng.com/api/qoss/collect")
	if err != nil {
		return fmt.Errorf("send failed, reason: %v", err)
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("send failed, reason: bad response status %v", resp.Status)
	}
	if result.Code != 200 {
		return fmt.Errorf("send failed, reason: %s", result.Message)
	}
	return nil
}
