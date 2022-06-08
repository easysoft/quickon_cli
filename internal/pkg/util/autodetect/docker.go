// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package autodetect

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/pkg/parsers"
	"github.com/ergoapi/util/file"
)

const (
	// constant for cgroup drivers
	cgroupFsDriver = "cgroupfs"
	// cgroupSystemdDriver = "systemd"
	// cgroupNoneDriver    = "none"
)

type Config struct {
	ExecOptions []string `json:"exec-opts,omitempty"`
}

// getCD gets the raw value of the native.cgroupdriver option, if set.
func getCD(config *Config) string {
	for _, option := range config.ExecOptions {
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil || !strings.EqualFold(key, "native.cgroupdriver") {
			continue
		}
		return val
	}
	return ""
}

// verifyCgroupDriver validates native.cgroupdriver
func verifyCgroupDriver(config *Config) error {
	cd := getCD(config)
	// if cd == "" || cd == cgroupFsDriver || cd == cgroupSystemdDriver {
	// 	return nil
	// }
	if cd != cgroupFsDriver {
		return fmt.Errorf("native.cgroupdriver option %s is internally used and cannot be specified manually\n you can run ~/.qc/data/hack/manifests/scripts/docker-daemon-patch.sh", cd)
	}
	return nil
}

// verifyDockerDaemonSettings performs validation of daemon config struct
func verifyDockerDaemonSettings(conf *Config) error {
	if err := verifyCgroupDriver(conf); err != nil {
		return err
	}
	return nil
}

func VerifyDockerDaemon() error {
	config := &Config{}
	if file.CheckFileExists("/etc/docker/daemon.json") {
		daemonDatam, _ := file.ReadAll("/etc/docker/daemon.json")
		json.Unmarshal(daemonDatam, config)
	}
	return verifyDockerDaemonSettings(config)
}
