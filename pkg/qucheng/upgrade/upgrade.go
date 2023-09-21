// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

import (
	"fmt"
	"os"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/ergoapi/util/color"

	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/version"
)

type Option struct {
	client *helm.Client
	log    log.Logger
}

func (opt *Option) Fetch(ns, name string) (ComponentVersion, error) {
	cmv := ComponentVersion{
		Name: name,
	}
	helmClient, err := helm.NewClient(&helm.Config{Namespace: ns})
	if err != nil {
		return cmv, errors.Errorf("create helm client err: %v", err)
	}
	opt.client = helmClient
	// update helm repo cache
	if err := opt.client.UpdateRepo(); err != nil {
		opt.log.Warn("update helm repo failed")
	}
	localcv, localav, err := opt.fetchDeploy(ns, name)
	if err != nil {
		opt.log.Debugf("fecth local %s failed, reason: %v", name, err)
	}
	cmv.Deploy.AppVersion = localav
	cmv.Deploy.ChartVersion = localcv
	// remote version
	remotecv, remoteav, err := opt.fetchCR(ns, name)
	if err != nil {
		opt.log.Debugf("fecth remote %s failed, reason: %v", name, err)
	}
	cmv.Remote.AppVersion = remoteav
	cmv.Remote.ChartVersion = remotecv
	// can upgrade
	cmv.CanUpgrade = version.LT(cmv.Remote.ChartVersion, cmv.Deploy.ChartVersion)
	if cmv.CanUpgrade {
		cmv.UpgradeMessage = fmt.Sprintf("Now you can use %s to upgrade component %s to the latest version", color.SGreen("%s upgrade %s", os.Args[0], name), name)
	}
	opt.log.Debugf("local: %s(%s), remote : %s(%s), upgrade: %v", localcv, localav, remotecv, remoteav, cmv.CanUpgrade)
	return cmv, err
}

// fetchDeploy get helm deployed version from k8s
func (opt *Option) fetchDeploy(ns, name string) (string, string, error) {
	result, _, err := opt.client.List(0, 0, name)
	if err != nil || len(result) < 1 {
		return "", "", errors.Errorf("not found chart: %s", name)
	}
	if len(result) > 1 {
		return "", "", errors.Errorf("chart more than 1, now count: %d", len(result))
	}
	last := result[0]
	return last.Chart.Metadata.Version, last.Chart.Metadata.AppVersion, nil
}

// fetchDeploy get helm remote version from cr
func (opt *Option) fetchCR(ns, name string) (string, string, error) {
	result, err := opt.client.GetLastCharts(common.DefaultHelmRepoName, name)
	if err != nil || len(result) < 1 {
		return "", "", errors.Errorf("not found chart: %s", name)
	}
	if len(result) > 1 {
		return "", "", errors.Errorf("chart more than 1, now count: %d", len(result))
	}
	last := result[0]
	return last.Chart.Version, last.Chart.AppVersion, nil
}

func Upgrade(flagVersion string, log log.Logger) error {
	helmClient, _ := helm.NewClient(&helm.Config{Namespace: common.GetDefaultSystemNamespace(true)})
	if err := helmClient.UpdateRepo(); err != nil {
		log.Errorf("update repo failed, reason: %v", err)
		return err
	}
	cfg, _ := config.LoadConfig()
	devops := false
	if cfg != nil && cfg.Quickon.DevOps {
		devops = true
	}
	qv, err := QuchengVersion(devops)
	if err != nil {
		return err
	}
	count := 0
	for _, cv := range qv.Components {
		if cv.CanUpgrade {
			defaultValue, _ := helmClient.GetValues(cv.Name)
			if _, err := helmClient.Upgrade(cv.Name, common.DefaultHelmRepoName, cv.Name, "", defaultValue); err != nil {
				log.Warnf("upgrade %s failed, reason: %v", cv.Name, err)
			} else {
				log.Donef("upgrade %s success", cv.Name)
				count++
			}
		} else {
			log.Infof("%s current version is the latest", cv.Name)
		}
	}
	if count == 0 {
		log.Done("the current version is the latest")
	}
	return nil
}

func QuchengVersion(devops bool) (Version, error) {
	v := Version{}
	opt := Option{
		log: log.GetInstance(),
	}
	name := common.DefaultQuchengName
	if devops {
		name = common.DefaultZentaoPaasName
	}
	if apiVersion, err := opt.Fetch(common.GetDefaultSystemNamespace(true), name); err == nil {
		v.Components = append(v.Components, apiVersion)
	}
	if apiVersion, err := opt.Fetch(common.GetDefaultSystemNamespace(true), common.DefaultCneOperatorName); err == nil {
		v.Components = append(v.Components, apiVersion)
	}
	return v, nil
}
