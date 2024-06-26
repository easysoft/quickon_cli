// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cne

type AppBackUPReq struct {
	Cluster   string `json:"cluster,omitempty"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`               // chartname
	UserName  string `json:"username,omitempty"` // 操作用户
}

type AppBackUPResp struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    AppBackUPData `json:"data"`
}

type AppBackUPData struct {
	BackupName string `json:"backup_name"`
	CreateTime int64  `json:"create_time"`
}

type AppBackUPStatusResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    AppBackUPStatus `json:"data"`
}

type AppBackUPStatus struct {
	Reason string `json:"reason,omitempty"`
	Status string `json:"status"`
}

type AppBackUPListResp struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    []AppBackUPListData `json:"data"`
}

type AppBackUPListData struct {
	Name          string           `json:"name"`
	Status        string           `json:"status"`
	ChartName     string           `json:"chart_name"`
	ChartVersion  string           `json:"chart_version"`
	BackupDetails BackupDetails    `json:"backup_details"`
	Restores      []BackupRestores `json:"restores"`
}

type BackupDetails struct {
	DBs     []DB     `json:"db"`
	Volumes []Volume `json:"volume"`
}

type DB struct {
	Type   string `json:"db_type"`
	Name   string `json:"db_name"`
	Status string `json:"status"`
	Cost   int64  `json:"cost"`
	Size   int64  `json:"size"`
}

type Volume struct {
	Name      string `json:"pvc_name"`
	Volume    string `json:"volume"`
	Status    string `json:"status"`
	Cost      int64  `json:"cost"`
	TotalSize int64  `json:"total_bytes"`
	DoneSize  int64  `json:"doneBytes"`
}

type BackupRestores struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
