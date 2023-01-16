// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package config

import (
	"os"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/ergoapi/util/file"
	"sigs.k8s.io/yaml"
)

var Cfg *Config

type SSH struct {
	User     string `json:"user,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
	PkName   string `json:"pkName,omitempty"`
	PkData   string `json:"pkData,omitempty"`
	Pk       string `json:"pk,omitempty"`
	PkPasswd string `json:"pkPasswd,omitempty"`
}

// Node node meta
type Node struct {
	Host string `yaml:"host" json:"host"`
	Name string `yaml:"name" json:"name"`
	Init bool   `yaml:"-" json:"-"`
}

// Config
type Config struct {
	Generated       time.Time `json:"-" yaml:"-"`
	ConsolePassword string    `yaml:"console-password,omitempty" json:"console-password,omitempty"`
	ClusterID       string    `yaml:"cluster_id" json:"cluster_id"`
	DB              string    `yaml:"db" json:"db"`
	Domain          string    `yaml:"domain" json:"domain"`
	APIToken        string    `yaml:"api_token" json:"api_token"`
	InitNode        string    `yaml:"init_node" json:"init_node"`
	Token           string    `yaml:"token" json:"token"`
	Master          []Node    `yaml:"master" json:"master"`
	Worker          []Node    `yaml:"worker" json:"worker"`
	S3              S3Config  `yaml:"s3" json:"s3"`
	DataDir         string    `yaml:"datadir" json:"datadir"`
}

type S3Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewConfig() *Config {
	return &Config{
		Generated: time.Now(),
	}
}

func LoadConfig() (*Config, error) {
	path := common.GetDefaultConfig()
	r := new(Config)
	if file.CheckFileExists(path) {
		b, _ := os.ReadFile(path)
		_ = yaml.Unmarshal(b, r)
	}
	return r, nil
}

func (r *Config) SaveConfig() error {
	path := common.GetDefaultConfig()
	b, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, b, common.FileMode0644)
	if err != nil {
		return err
	}
	return nil
}
