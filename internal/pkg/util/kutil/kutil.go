// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package kutil

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ergoapi/util/exstr"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/version"
	"github.com/ergoapi/util/ztime"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
)

const (
	NodeToken = "/var/lib/rancher/k3s/server/node-token"
)

func GetNodeToken() string {
	b, err := os.ReadFile(NodeToken)
	if err != nil {
		return ""
	}
	return string(b)
}

// NeedCacheHelmFile helm repo update
func NeedCacheHelmFile() bool {
	cachefile := fmt.Sprintf("%s/.566964cd0285e57cd088caa251ae863a.lock", common.DefaultCacheDir)
	if file.CheckFileExists(cachefile) {
		data, _ := file.ReadAll(cachefile)
		old := time.Unix(exstr.Str2Int64(string(data)), 0)
		now := time.Now()
		if now.Sub(old) > 2*time.Minute {
			os.Remove(cachefile)
			file.WriteFile(cachefile, ztime.NowUnixString(), true)
			return true
		}
		return false
	}
	file.WriteFile(cachefile, ztime.NowUnixString(), true)
	return true
}

// IsLegalDomain check domain legal
func IsLegalDomain(host string) bool {
	for _, d := range common.ValidDomainSuffix {
		if strings.HasSuffix(host, d) {
			return true
		}
	}
	return false
}

func SplitDomain(domain string) (string, string) {
	for _, d := range common.ValidDomainSuffix {
		if strings.HasSuffix(domain, d) {
			return strings.ReplaceAll(domain, "."+d, ""), d
		}
	}
	return domain, common.ValidDomainSuffix[0]
}

func GetConsoleURL(cfg *config.Config) string {
	domain := cfg.Domain
	if len(domain) > 0 {
		prekey := "console."
		if cfg.Quickon.DevOps {
			prekey = ""
			if len(cfg.Install.Version) == 0 || version.LTv2(cfg.Install.Version, common.Version) {
				prekey = "zentao."
			}
		}
		if !IsLegalDomain(cfg.Domain) || cfg.Quickon.Domain.Type != "api" {
			domain = fmt.Sprintf("http://%s%s", prekey, cfg.Domain)
		} else {
			domain = fmt.Sprintf("https://%s", cfg.Domain)
		}
	} else {
		domain = fmt.Sprintf("http://%s:32379", cfg.Cluster.InitNode)
	}
	return domain
}

// ValidNamespace valid namespace
func ValidNamespace(ns string) bool {
	switch ns {
	case common.DefaultAppNamespace, common.DefaultSystemNamespace:
		return true
	default:
		return false
	}
}
