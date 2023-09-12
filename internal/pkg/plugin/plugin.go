// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	gv "github.com/Masterminds/semver/v3"
	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/common"
	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	corev1 "k8s.io/api/core/v1"
	kubeerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetAll() ([]Meta, error) {
	log := log.GetInstance()
	var plugins []Meta
	pf := fmt.Sprintf("%s/hack/manifests/plugins/plugins.json", common.GetDefaultDataDir())
	log.Debug("load local plugin config from", pf)
	content, err := os.ReadFile(pf)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &plugins)
	if err != nil {
		log.Errorf("unmarshal plugin meta failed: %v", err)
		return nil, err
	}
	return plugins, nil
}

func GetMaps() (map[string]Meta, error) {
	plugins, err := GetAll()
	if err != nil {
		return nil, err
	}
	maps := make(map[string]Meta)
	for _, p := range plugins {
		maps[p.Type] = p
	}
	return maps, nil
}

func GetMeta(args ...string) (Item, error) {
	log := log.GetInstance()
	ps, err := GetMaps()
	if err != nil {
		return Item{}, err
	}
	t := args[0]
	name := ""
	if len(args) == 2 {
		name = args[1]
	} else if strings.Contains(t, "/") {
		ts := strings.Split(t, "/")
		t = ts[0]
		if len(ts) > 1 {
			name = ts[1]
		}
	}
	var plugin Item
	if v, ok := ps[t]; ok {
		if name == "" {
			name = v.Default
		}
		exist := false
		for _, item := range v.Item {
			if item.Name == name {
				exist = true
				plugin = item
				plugin.Type = v.Type
				break
			}
		}
		if !exist {
			log.Warnf("%s not found %s, will use default: %s", t, name, v.Default)
			return GetMeta(t, v.Default)
		}
		log.Infof("install %s plugin: %s", t, name)
		plugin.log = log
		return plugin, nil
	}
	return Item{}, errors.Errorf("plugin %s not found", t)
}

func (p *Item) UnInstall() error {
	if p.BuiltIn {
		p.log.Warnf("builtin plugin %s cannot be uninstalled", p.Type)
		return nil
	}
	pluginName := fmt.Sprintf("qc-plugin-%s", p.Type)
	_, err := p.Client.GetSecret(context.TODO(), common.GetDefaultSystemNamespace(true), pluginName, metav1.GetOptions{})
	if err != nil {
		if kubeerr.IsNotFound(err) {
			p.log.Warnf("plugin %s is already uninstalled", p.Type)
			return nil
		}
		p.log.Fatalf("get plugin secret failed: %v", err)
		return nil
	}
	// #nosec
	if p.Tool == "helm" {
		applycmd := qcexec.Command(os.Args[0], "experimental", "helm", "delete", p.Type, "-n", common.GetDefaultSystemNamespace(true))
		if output, err := applycmd.CombinedOutput(); err != nil {
			p.log.Errorf("helm uninstall %s plugin %s failed: %s", p.Type, p.Name, string(output))
			return err
		}
	} else {
		// #nosec
		applycmd := qcexec.Command(os.Args[0], "experimental", "kubectl", "delete", "-f", fmt.Sprintf("%s/%s", common.GetDefaultDataDir(), p.Path), "-n", common.GetDefaultSystemNamespace(true), "--kubeconfig", common.GetKubeConfig())
		if output, err := applycmd.CombinedOutput(); err != nil {
			p.log.Errorf("kubectl uninstall %s plugin %s failed: %s", p.Type, p.Name, string(output))
			return err
		}
	}
	p.log.Donef("uninstall %s plugin success", p.Type)
	p.Client.DeleteSecret(context.TODO(), common.GetDefaultSystemNamespace(true), pluginName, metav1.DeleteOptions{})
	return nil
}

func (p *Item) Install() error {
	pluginName := fmt.Sprintf("qc-plugin-%s", p.Type)
	oldSecret, err := p.Client.GetSecret(context.TODO(), common.GetDefaultSystemNamespace(true), pluginName, metav1.GetOptions{})
	updatestatus := false
	if err == nil {
		nowversion := gv.MustParse(strings.TrimPrefix(p.Version, "v"))
		oldversion := string(oldSecret.Data["version"])
		p.log.Debugf("type: %s, old version: %s, now version: %s", p.Type, oldversion, nowversion)
		needupgrade := nowversion.GreaterThan(gv.MustParse(oldversion))
		if !needupgrade {
			p.log.Warnf("plugin %s is the latest version", p.Type)
			return nil
		}
		updatestatus = true
	} else {
		if !kubeerr.IsNotFound(err) {
			p.log.Debugf("get plugin secret failed: %v", err)
			return errors.Errorf("plugin %s install failed", p.Name)
		}
	}
	if p.Tool == "helm" {
		args := []string{"experimental", "helm", "upgrade", "--name", p.Type, "--repo", common.DefaultHelmRepoName, "--chart", p.Path, "--namespace", common.GetDefaultSystemNamespace(true)}
		if len(p.InstallVersion) > 0 {
			args = append(args, "--version", p.InstallVersion)
		}
		applycmd := qcexec.Command(os.Args[0], args...)
		if output, err := applycmd.CombinedOutput(); err != nil {
			p.log.Errorf("helm install %s plugin %s failed: %s", p.Type, p.Name, string(output))
			return err
		}
	} else {
		// #nosec
		applycmd := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", fmt.Sprintf("%s/%s", common.GetDefaultDataDir(), p.Path), "-n", common.GetDefaultSystemNamespace(true), "--kubeconfig", common.GetKubeConfig())
		if output, err := applycmd.CombinedOutput(); err != nil {
			p.log.Errorf("kubectl install %s plugin %s failed: %s", p.Type, p.Name, string(output))
			return err
		}
	}

	p.log.Donef("upgrade install %s plugin %s success", p.Type, p.Name)
	plugindata := map[string]string{
		"type":       p.Type,
		"name":       p.Name,
		"version":    p.Version,
		"cliversion": common.Version,
	}
	if updatestatus {
		_, err = p.Client.UpdateSecret(context.TODO(), common.GetDefaultSystemNamespace(true), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: pluginName,
			},
			StringData: plugindata,
		}, metav1.UpdateOptions{})
	} else {
		_, err = p.Client.CreateSecret(context.TODO(), common.GetDefaultSystemNamespace(true), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: pluginName,
			},
			StringData: plugindata,
		}, metav1.CreateOptions{})
	}

	return err
}

func (p *Item) Upgrade() (err error) {
	pluginName := fmt.Sprintf("qc-plugin-%s", p.Type)
	oldSecret, _ := p.Client.GetSecret(context.TODO(), common.GetDefaultSystemNamespace(true), pluginName, metav1.GetOptions{})
	updatestatus := true
	if oldSecret == nil {
		updatestatus = false
	}

	if p.Tool == "helm" {
		applycmd := qcexec.Command(os.Args[0], "experimental", "helm", "upgrade", "--name", p.Type, "--repo", common.DefaultHelmRepoName, "--chart", p.Path, "--namespace", common.GetDefaultSystemNamespace(true))
		if output, err := applycmd.CombinedOutput(); err != nil {
			p.log.Errorf("helm upgrade %s plugin %s failed: %s", p.Type, p.Name, string(output))
			return err
		}
	} else {
		// #nosec
		applycmd := qcexec.Command(os.Args[0], "experimental", "kubectl", "apply", "-f", fmt.Sprintf("%s/%s", common.GetDefaultDataDir(), p.Path), "-n", common.GetDefaultSystemNamespace(true), "--kubeconfig", common.GetKubeConfig())
		if output, err := applycmd.CombinedOutput(); err != nil {
			p.log.Errorf("kubectl upgrade %s plugin %s failed: %s", p.Type, p.Name, string(output))
			return err
		}
	}
	p.log.Donef("upgrade %s plugin %s success", p.Type, p.Name)
	plugindata := map[string]string{
		"type":       p.Type,
		"name":       p.Name,
		"version":    p.Version,
		"cliversion": common.Version,
	}
	if updatestatus {
		oldSecret.StringData = plugindata
		_, err = p.Client.UpdateSecret(context.TODO(), common.GetDefaultSystemNamespace(true), oldSecret, metav1.UpdateOptions{})
	} else {
		_, err = p.Client.CreateSecret(context.TODO(), common.GetDefaultSystemNamespace(true), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: pluginName,
			},
			StringData: plugindata,
		}, metav1.CreateOptions{})
	}

	return err
}
