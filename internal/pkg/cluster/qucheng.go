// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package cluster

import (
	"context"
	"strings"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/binfile"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/ergoapi/util/expass"
	"github.com/ergoapi/util/file"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (p *Cluster) genQuChengToken() string {
	// TODO token 生成优化
	return expass.RandomPassword(32)
}

func (p *Cluster) InstallQuCheng() error {
	log.Flog.Info("executing init qucheng logic...")
	ctx := context.Background()
	log.Flog.Debug("waiting for storage to be ready...")
	waitsc := time.Now()
	// wait.BackoffUntil TODO
	for {
		sc, _ := p.KubeClient.GetDefaultSC(ctx)
		if sc != nil {
			log.Flog.Donef("default storage %s is ready", sc.Name)
			break
		}
		time.Sleep(time.Second * 5)
		trywaitsc := time.Now()
		if trywaitsc.Sub(waitsc) > time.Minute*3 {
			log.Flog.Warnf("wait storage %s ready, timeout: %v", sc.Name, trywaitsc.Sub(waitsc).Seconds())
			break
		}
	}

	_, err := p.KubeClient.CreateNamespace(ctx, common.DefaultSystem, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}
	log.Flog.Debug("start init qucheng")
	getbin := binfile.Meta{}
	helmbin, err := getbin.LoadLocalBin(common.HelmBinName)
	if err != nil {
		return err
	}
	output, err := qcexec.Command(helmbin, "repo", "add", "install", common.GetChartRepo(p.QuchengVersion)).CombinedOutput()
	if err != nil {
		errmsg := string(output)
		if !strings.Contains(errmsg, "exists") {
			log.Flog.Errorf("init qucheng install repo failed: %s", string(output))
			return err
		}
		log.Flog.Warn("qucheng install repo  already exists")
	} else {
		log.Flog.Done("init qucheng install repo done")
	}

	output, err = qcexec.Command(helmbin, "repo", "update").CombinedOutput()
	if err != nil {
		log.Flog.Errorf("update qucheng install repo failed: %s", string(output))
		return err
	}
	log.Flog.Done("update qucheng install repo done")
	token := p.genQuChengToken()
	// helm upgrade -i nginx-ingress-controller bitnami/nginx-ingress-controller -n kube-system
	output, err = qcexec.Command(helmbin, "upgrade", "-i", common.DefaultQuchengName, common.DefaultChartName, "-n", common.DefaultSystem, "--set", "env.CNE_API_TOKEN="+token, "--set", "cloud.defaultChannel="+common.GetChannel(p.QuchengVersion)).CombinedOutput()
	if err != nil {
		log.Flog.Errorf("upgrade install qucheng web failed: %s", string(output))
		return err
	}
	// Deprecated CNE_API_TOKEN
	output, err = qcexec.Command(helmbin, "upgrade", "-i", common.DefaultCneAPIName, common.DefaultAPIChartName, "-n", common.DefaultSystem, "--set", "env.CNE_TOKEN="+token, "--set", "env.CNE_API_TOKEN="+token, "--set", "cloud.defaultChannel="+common.GetChannel(p.QuchengVersion)).CombinedOutput()
	if err != nil {
		log.Flog.Errorf("upgrade install qucheng api failed: %s", string(output))
		return err
	}
	log.Flog.Done("install qucheng done")
	p.Ready()
	initfile := common.GetCustomConfig(common.InitFileName)
	if err := file.Writefile(initfile, "init done"); err != nil {
		log.Flog.Warnf("write init done file failed, reason: %v.\n\t please run: touch %s", err, initfile)
	}
	return nil
}
