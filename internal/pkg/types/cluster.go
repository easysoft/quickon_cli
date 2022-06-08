// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Metadata struct {
	Name     string `json:"name" yaml:"name"`
	Provider string `json:"provider" yaml:"provider"`
	Master   string `json:"master" yaml:"master"`
	Worker   string `json:"worker" yaml:"worker"`
	// Token           string      `json:"token,omitempty" yaml:"token,omitempty"`
	IP              string      `json:"ip,omitempty" yaml:"ip,omitempty"`
	EIP             string      `json:"eip,omitempty" yaml:"eip,omitempty"`
	TLSSans         StringArray `json:"tls-sans,omitempty" yaml:"tls-sans,omitempty"`
	ClusterCidr     string      `json:"cluster-cidr,omitempty" yaml:"cluster-cidr,omitempty"`
	ServiceCidr     string      `json:"service-cidr,omitempty" yaml:"service-cidr,omitempty"`
	MasterExtraArgs string      `json:"master-extra-args,omitempty" yaml:"master-extra-args,omitempty"`
	WorkerExtraArgs string      `json:"worker-extra-args,omitempty" yaml:"worker-extra-args,omitempty"`
	DataStore       string      `json:"datastore,omitempty" yaml:"datastore,omitempty"`
	Network         string      `json:"network,omitempty" yaml:"network,omitempty"`
	// Plugins         StringArray `json:"plugins,omitempty" yaml:"plugins,omitempty"`
	Mode           string `json:"mode,omitempty" yaml:"mode,omitempty"`
	QuchengVersion string `json:"qucheng-version,omitempty" yaml:"qucheng-version,omitempty"`
	DisableIngress bool   `json:"disable-ingress,omitempty" yaml:"disable-ingress,omitempty"`
	CNEAPI         string `json:"cne-api,omitempty" yaml:"cne-api,omitempty"`
	CNEToken       string `json:"cne-token,omitempty" yaml:"cne-token,omitempty"`
}

type Status struct {
	Status string `json:"status,omitempty"`
}

// Flag struct for flag.
type Flag struct {
	Name      string
	P         interface{}
	V         interface{}
	ShortHand string
	Usage     string
	Required  bool
	EnvVar    string
}

type StringArray []string

// Scan gorm Scan implement.
func (a *StringArray) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		if v != "" {
			*a = strings.Split(v, ",")
		}
	default:
		return fmt.Errorf("failed to scan array value %v", value)
	}
	return nil
}

// Value gorm Value implement.
func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return strings.Join(a, ","), nil
}
