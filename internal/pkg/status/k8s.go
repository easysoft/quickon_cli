// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package status

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/file"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/plugin"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"

	corev1 "k8s.io/api/core/v1"
	kubeerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sStatusOption struct {
	Namespace      string
	KubeConfig     string
	Wait           bool
	WaitDuration   time.Duration
	IgnoreWarnings bool
	ListOutput     string
	Log            log.Logger
}

type K8sStatusCollector struct {
	client *k8s.Client
	option K8sStatusOption
	cfg    *config.Config
}

func NewK8sStatusCollector(option K8sStatusOption) (*K8sStatusCollector, error) {
	client, err := k8s.NewClient("", option.KubeConfig)
	cfg, _ := config.LoadConfig()
	return &K8sStatusCollector{
		client: client,
		option: option,
		cfg:    cfg,
	}, err
}

func (s K8sStatusOption) waitTimeout() time.Duration {
	if s.WaitDuration != time.Duration(0) {
		return s.WaitDuration
	}

	return common.StatusWaitDuration
}

func (k *K8sStatusCollector) Status(ctx context.Context) (*Status, error) {
	var mostRecentStatus *Status
	ctx, cancel := context.WithTimeout(ctx, k.option.waitTimeout())
	defer cancel()

retry:
	select {
	case <-ctx.Done():
		return mostRecentStatus, errors.Errorf("timeout while waiting for status to become successful: %w", ctx.Err())
	default:
	}
	s := k.status(ctx)
	if s != nil {
		mostRecentStatus = s
	}
	if !k.statusIsReady(s) && k.option.Wait {
		time.Sleep(common.WaitRetryInterval)
		goto retry
	}

	return mostRecentStatus, nil
}

func (k *K8sStatusCollector) statusIsReady(s *Status) bool {
	// TODO Check status
	return true
}

func (k *K8sStatusCollector) status(ctx context.Context) *Status {
	status := newStatus(k.option.ListOutput)
	if err := k.nodesStatus(ctx, status); err != nil {
		k.option.Log.Errorf("failed to get nodes status: %v", err)
	}
	if err := k.serviceStatus(ctx, status); err != nil {
		k.option.Log.Errorf("failed to get service status: %v", err)
	}
	return status
}

func (k *K8sStatusCollector) deploymentStatus(ctx context.Context, ns, name, aliasname, t string, status *Status) (bool, error) {
	k.option.Log.Debugf("check cm %s status", aliasname)
	stateCount := PodStateCount{Type: "Deployment"}
	d, err := k.client.GetDeployment(ctx, ns, name, metav1.GetOptions{})
	if kubeerr.IsNotFound(err) {
		stateCount.Disabled = true
		if t == "k8s" {
			status.KubeStatus.PodState[aliasname] = stateCount
		} else {
			status.QStatus.PodState[aliasname] = stateCount
		}
		return true, nil
	}

	if err != nil {
		stateCount.Disabled = false
		if t == "k8s" {
			status.KubeStatus.PodState[aliasname] = stateCount
		} else {
			status.QStatus.PodState[aliasname] = stateCount
		}
		return false, err
	}

	if d == nil {
		stateCount.Disabled = false
		if t == "k8s" {
			status.KubeStatus.PodState[aliasname] = stateCount
		} else {
			status.QStatus.PodState[aliasname] = stateCount
		}
		return false, errors.Errorf("component %s is not available", aliasname)
	}

	stateCount.Desired = int(d.Status.Replicas)
	stateCount.Ready = int(d.Status.ReadyReplicas)
	stateCount.Available = int(d.Status.AvailableReplicas)
	stateCount.Unavailable = int(d.Status.UnavailableReplicas)
	stateCount.Disabled = false
	if t == "k8s" {
		status.KubeStatus.PodState[aliasname] = stateCount
	} else {
		status.QStatus.PodState[aliasname] = stateCount
	}
	notReady := stateCount.Desired - stateCount.Ready
	if notReady > 0 {
		k.option.Log.Warnf("%d pods of Deployment %s are not ready", notReady, name)
	}
	if unavailable := stateCount.Unavailable - notReady; unavailable > 0 {
		k.option.Log.Warnf("%d pods of Deployment %s are not available", unavailable, name)
	}
	return false, nil
}

func (k *K8sStatusCollector) daemonsetStatus(ctx context.Context, ns, name, aliasname, t string, status *Status) (bool, error) {
	k.option.Log.Debugf("check cm %s status", aliasname)
	stateCount := PodStateCount{Type: "Daemonset"}
	d, err := k.client.GetDaemonSet(ctx, ns, name, metav1.GetOptions{})
	if kubeerr.IsNotFound(err) {
		stateCount.Disabled = true
		if t == "k8s" {
			status.KubeStatus.PodState[aliasname] = stateCount
		} else {
			status.QStatus.PodState[aliasname] = stateCount
		}
		return true, nil
	}

	if err != nil {
		stateCount.Disabled = false
		if t == "k8s" {
			status.KubeStatus.PodState[aliasname] = stateCount
		} else {
			status.QStatus.PodState[aliasname] = stateCount
		}
		return false, err
	}

	if d == nil {
		stateCount.Disabled = false
		if t == "k8s" {
			status.KubeStatus.PodState[aliasname] = stateCount
		} else {
			status.QStatus.PodState[aliasname] = stateCount
		}
		return false, errors.Errorf("component %s is not available", aliasname)
	}

	stateCount.Ready = int(d.Status.NumberReady)
	stateCount.Available = int(d.Status.NumberAvailable)
	stateCount.Unavailable = int(d.Status.NumberUnavailable)
	stateCount.Disabled = false
	if t == "k8s" {
		status.KubeStatus.PodState[aliasname] = stateCount
	} else {
		status.QStatus.PodState[aliasname] = stateCount
	}
	notReady := stateCount.Desired - stateCount.Ready
	if notReady > 0 {
		k.option.Log.Warnf("%d pods of Deployment %s are not ready", notReady, name)
	}
	if unavailable := stateCount.Unavailable - notReady; unavailable > 0 {
		k.option.Log.Warnf("%d pods of Deployment %s are not available", unavailable, name)
	}
	return false, nil
}

// serviceStatus 检查服务状态
func (k *K8sStatusCollector) serviceStatus(ctx context.Context, status *Status) error {
	// 集群
	k.deploymentStatus(ctx, "kube-system", "coredns", "coredns", "k8s", status)
	k.deploymentStatus(ctx, "kube-system", "metrics-server", "metrics-server", "k8s", status)
	if k.cfg.Storage.Type == "local" {
		k.deploymentStatus(ctx, "kube-system", "local-path-provisioner", "local-path-provisioner", "k8s", status)
	}
	// 业务层
	k.deploymentStatus(ctx, common.GetDefaultSystemNamespace(true), common.GetReleaseName(k.cfg.Quickon.DevOps), common.GetReleaseName(k.cfg.Quickon.DevOps), "", status)
	// 数据库
	db := fmt.Sprintf("%s-mysql", common.GetReleaseName(k.cfg.Quickon.DevOps))
	k.deploymentStatus(ctx, common.GetDefaultSystemNamespace(true), db, db, "", status)

	// 插件状态
	plugins, _ := plugin.GetMaps()
	for _, p := range plugins {
		k.platformPluginStatus(ctx, p, status)
	}
	return nil
}

func (k *K8sStatusCollector) platformPluginStatus(ctx context.Context, p plugin.Meta, status *Status) error {
	k.option.Log.Debugf("check platform plugin %s status", p.Type)
	stateCount := PodStateCount{Type: "Plugin"}
	_, err := k.client.GetSecret(ctx, common.GetDefaultSystemNamespace(true), "qc-plugin-"+p.Type, metav1.GetOptions{})
	if err != nil {
		stateCount.Disabled = true
	} else {
		stateCount.Disabled = false
		if p.Type == "ingress" {
			k.ingressStatus(ctx, common.GetDefaultSystemNamespace(true), fmt.Sprintf("ingress-%s", common.DefaultIngressName), common.DefaultIngressName, status)
		} else if p.Type == common.DefaultCneOperatorName {
			k.deploymentStatus(ctx, common.GetDefaultSystemNamespace(true), common.DefaultCneOperatorName, common.DefaultCneOperatorName, "", status)
		}
	}
	status.QStatus.PluginState[p.Type] = stateCount
	return nil
}

func (k *K8sStatusCollector) nodesStatus(ctx context.Context, status *Status) error {
	nodes, err := k.client.ListNodes(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	nodetotal := len(nodes.Items)
	status.KubeStatus.NodeCount["total"] = nodetotal
	readynode := 0
	master := 0
	worker := 0
	for _, node := range nodes.Items {
		for _, nc := range node.Status.Conditions {
			if nc.Type == "Ready" && nc.Status == corev1.ConditionTrue {
				readynode++
			}
		}
		if node.Labels["node-role.kubernetes.io/master"] == "true" || node.Labels["node-role.kubernetes.io/control-plane"] == "true" {
			master++
		} else {
			worker++
		}
	}
	status.KubeStatus.NodeCount["ready"] = readynode
	status.KubeStatus.NodeCount["master"] = master
	status.KubeStatus.NodeCount["worker"] = worker
	if master == 0 {
		status.KubeStatus.Type = "cloud-managed"
	} else {
		if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
			status.KubeStatus.Type = "self-managed"
		} else {
			status.KubeStatus.Type = "owner-self-managed"
		}
	}
	version, _ := k.client.GetGitVersion(ctx)
	if len(version) > 0 {
		status.KubeStatus.Version = version
	}
	return nil
}

func (k *K8sStatusCollector) ingressStatus(ctx context.Context, ns, name, aliasname string, status *Status) (bool, error) {
	k.option.Log.Debugf("check cm %s status", aliasname)
	var stateCount PodStateCount
	ds, _ := k.client.GetDaemonSet(ctx, ns, name, metav1.GetOptions{})
	deploy, _ := k.client.GetDeployment(ctx, ns, name, metav1.GetOptions{})
	if ds == nil && deploy == nil {
		stateCount.Disabled = true
		status.QStatus.PodState[aliasname] = stateCount
		return true, nil
	}
	stateCount.Disabled = false
	if ds != nil && !ds.CreationTimestamp.IsZero() {
		k.option.Log.Debugf("detch %s kind DaemonSet", name)
		stateCount.Type = "DaemonSet"
		stateCount.Desired = int(ds.Status.DesiredNumberScheduled)
		stateCount.Ready = int(ds.Status.NumberReady)
		stateCount.Available = int(ds.Status.NumberAvailable)
		stateCount.Unavailable = int(ds.Status.NumberUnavailable)
		notReady := stateCount.Desired - stateCount.Ready
		if notReady > 0 {
			k.option.Log.Warnf("%d pods of DaemonSet %s are not ready", notReady, name)
		}
		if unavailable := stateCount.Unavailable - notReady; unavailable > 0 {
			k.option.Log.Warnf("%d pods of DaemonSet %s are not available", unavailable, name)
		}
	} else if deploy != nil && !deploy.CreationTimestamp.IsZero() {
		k.option.Log.Debugf("detch %s kind Deployment", name)
		stateCount.Type = "Deployment"
		if *deploy.Spec.Replicas > 0 {
			stateCount.Desired = int(deploy.Status.Replicas)
			stateCount.Ready = int(deploy.Status.ReadyReplicas)
			stateCount.Available = int(deploy.Status.AvailableReplicas)
			stateCount.Unavailable = int(deploy.Status.UnavailableReplicas)
			notReady := stateCount.Desired - stateCount.Ready
			if notReady > 0 {
				k.option.Log.Warnf("%d pods of Deployment %s are not ready", notReady, name)
			}
			if unavailable := stateCount.Unavailable - notReady; unavailable > 0 {
				k.option.Log.Warnf("%d pods of Deployment %s are not available", unavailable, name)
			}
		} else {
			k.option.Log.Warnf("Deployment %s disabled", name)
		}
	}
	status.QStatus.PodState[aliasname] = stateCount
	return false, nil
}
