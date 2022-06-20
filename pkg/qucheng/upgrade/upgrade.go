// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package upgrade

import (
	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/helm"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/version"
	"github.com/pkg/errors"
)

type Option struct {
	client *helm.Client
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
		log.Flog.Warn("update helm repo failed")
	}
	// TODO fetch local version
	localcv, localav, err := opt.fetchDeploy(ns, name)
	if err != nil {
		log.Flog.Debugf("fecth local %s failed, reason: %v", name, err)
	}
	cmv.Deploy.AppVersion = localav
	cmv.Deploy.ChartVersion = localcv
	// remote version
	remotecv, remoteav, err := opt.fetchCR(ns, name)
	if err != nil {
		log.Flog.Debugf("fecth remote %s failed, reason: %v", name, err)
	}
	cmv.Remote.AppVersion = remoteav
	cmv.Remote.ChartVersion = remotecv
	// can upgrade
	if cmv.Deploy.ChartVersion != cmv.Remote.ChartVersion && version.LT(cmv.Deploy.ChartVersion, cmv.Remote.ChartVersion) {
		cmv.CanUpgrade = true
	}
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

func Upgrade(flagVersion string) error {
	helmClient, _ := helm.NewClient(&helm.Config{Namespace: common.DefaultSystem})
	if err := helmClient.UpdateRepo(); err != nil {
		log.Flog.Errorf("update repo failed, reason: %v", err)
		return err
	}

	qv, err := QuchengVersion()
	if err != nil {
		return err
	}
	count := 0
	for _, cv := range qv.Components {
		if cv.CanUpgrade {
			defaultValue, _ := helmClient.GetValues(cv.Name)
			if _, err := helmClient.Upgrade(cv.Name, common.DefaultHelmRepoName, cv.Name, "", defaultValue); err != nil {
				log.Flog.Warnf("upgrade %s failed, reason: %v", cv.Name, err)
			} else {
				log.Flog.Donef("upgrade %s success", cv.Name)
				count++
			}
		}
	}
	if count == 0 {
		log.Flog.Done("the current version is the latest")
	}
	return nil
}

func QuchengVersion() (Version, error) {
	v := Version{}
	opt := Option{}
	if uiVersion, err := opt.Fetch(common.DefaultSystem, "cne-api"); err == nil {
		v.Components = append(v.Components, uiVersion)
	}
	if apiVersion, err := opt.Fetch(common.DefaultSystem, "qucheng"); err == nil {
		v.Components = append(v.Components, apiVersion)
	}
	return v, nil
}
