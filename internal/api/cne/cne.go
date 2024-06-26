// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cne

import (
	"fmt"
	"reflect"

	"github.com/imroc/req/v3"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
)

type CneAPI struct {
	Req      *req.Request
	Endpoint string
}

func NewCneAPI() *CneAPI {
	config, _ := config.LoadConfig()
	c := req.C().SetUserAgent(common.GetUG()).R().SetHeader(common.CneAPITokenHeader, config.APIToken)
	return &CneAPI{
		Req:      c,
		Endpoint: fmt.Sprintf("http://%s:32380", config.Cluster.InitNode),
	}
}

func (c *CneAPI) CreateAppBackUP(ns, chartName string) (string, error) {
	var result AppBackUPResp
	appbackup := &AppBackUPReq{Namespace: ns, Name: chartName}
	resp, err := c.Req.SetBody(appbackup).
		SetSuccessResult(&result).
		Post(c.Endpoint + "/api/cne/app/backup")
	if resp.IsSuccessState() {
		return result.Data.BackupName, nil
	}
	return "", err
}

func (c *CneAPI) AppBackUPStatus(ns, chartName, backupName string) (*AppBackUPStatus, error) {
	var result AppBackUPStatusResp
	resp, err := c.Req.SetQueryParams(map[string]string{
		"namespace":   ns,
		"name":        chartName,
		"backup_name": backupName,
	}).
		SetSuccessResult(&result).
		Get(c.Endpoint + "/api/cne/app/backup/status")
	if resp.IsSuccessState() {
		return &result.Data, nil
	}
	return nil, err
}

func (c *CneAPI) ListAppBackUP(ns, chartName string) ([]AppBackUPListData, error) {
	var result AppBackUPListResp
	resp, err := c.Req.SetQueryParams(map[string]string{
		"namespace": ns,
		"name":      chartName,
	}).
		SetSuccessResult(&result).
		Get(c.Endpoint + "/api/cne/app/backups")
	if resp.IsSuccessState() {
		return result.Data, nil
	}
	return nil, err
}

func struct2map(obj any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		if field.Kind() == reflect.Struct {
			result[fieldName] = struct2map(field.Interface())
		} else {
			result[fieldName] = field.Interface()
		}
	}
	return result
}
