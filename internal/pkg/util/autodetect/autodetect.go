// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package autodetect

import "context"

type ValidationCheck interface {
	Name() string
	Check(ctx context.Context) error
	NeedFix() bool
	Fix(ctx context.Context) error
}

// podcidrValidationCheck is a ValidationCheck for podcidr
type podcidrValidationCheck struct {
	podcidr string
}

func (v *podcidrValidationCheck) Name() string {
	return "podcidr"
}

func (v *podcidrValidationCheck) Check(ctx context.Context) error {
	return nil
}

func (v *podcidrValidationCheck) NeedFix() bool {
	return true
}

func (v *podcidrValidationCheck) Fix(ctx context.Context) error {
	return nil
}

// NewPodCIDRValidationCheck returns a new ValidationCheck for podcidr
func NewPodCIDRValidationCheck(podcidr string) ValidationCheck {
	return &podcidrValidationCheck{
		podcidr: podcidr,
	}
}

// svccidrValidationCheck is a ValidationCheck for svccidr
type svccidrValidationCheck struct {
	svccidr string
}

func (v *svccidrValidationCheck) Name() string {
	return "svccidr"
}

func (v *svccidrValidationCheck) Check(ctx context.Context) error {
	return nil
}

func (v *svccidrValidationCheck) NeedFix() bool {
	return true
}

func (v *svccidrValidationCheck) Fix(ctx context.Context) error {
	return nil
}

// NewSvcCIDRValidationCheck returns a new ValidationCheck for svccidr
func NewSvcCIDRValidationCheck(svccidr string) ValidationCheck {
	return &svccidrValidationCheck{
		svccidr: svccidr,
	}
}

// dockerValidationCheck is a ValidationCheck for docker
type dockerValidationCheck struct{}

func (v *dockerValidationCheck) Name() string {
	return "docker"
}

func (v *dockerValidationCheck) Check(ctx context.Context) error {
	return nil
}

func (v *dockerValidationCheck) NeedFix() bool {
	return true
}

func (v *dockerValidationCheck) Fix(ctx context.Context) error {
	return nil
}

// NewDockerValidationCheck returns a new ValidationCheck for docker
func NewDockerValidationCheck() ValidationCheck {
	return &dockerValidationCheck{}
}

// KubeValidationCheck is a ValidationCheck for k8s
type KubeValidationCheck struct{}

func (v *KubeValidationCheck) Name() string {
	return "k8s"
}

func (v *KubeValidationCheck) Check(ctx context.Context) error {
	return nil
}

func (v *KubeValidationCheck) NeedFix() bool {
	return true
}

func (v *KubeValidationCheck) Fix(ctx context.Context) error {
	return nil
}

func NewKubeValidationCheck() ValidationCheck {
	return &KubeValidationCheck{}
}

// K8sVersion is a ValidationCheck for k8s version
type K8sVersionValidationCheck struct{}

func (v *K8sVersionValidationCheck) Name() string {
	return "k8s-version"
}

func (v *K8sVersionValidationCheck) Check(ctx context.Context) error {
	return nil
}

func (v *K8sVersionValidationCheck) NeedFix() bool {
	return true
}

func (v *K8sVersionValidationCheck) Fix(ctx context.Context) error {
	return nil
}

func NewK8sVersionValidationCheck() ValidationCheck {
	return &K8sVersionValidationCheck{}
}

// KernelModules is a ValidationCheck for linux kernel
type KernelModulesValidationCheck struct{}

func (v *KernelModulesValidationCheck) Name() string {
	return "kernel-modules"
}

func (v *KernelModulesValidationCheck) Check(ctx context.Context) error {
	return nil
}

func (v *KernelModulesValidationCheck) NeedFix() bool {
	return true
}

func (v *KernelModulesValidationCheck) Fix(ctx context.Context) error {
	return nil
}

func NewKernelModulesValidationCheck() ValidationCheck {
	return &KernelModulesValidationCheck{}
}

// Sysctl is a ValidationCheck for sysctl
type SysctlValidationCheck struct{}

func (v *SysctlValidationCheck) Name() string {
	return "sysctl"
}

func (v *SysctlValidationCheck) Check(ctx context.Context) error {
	return nil
}

func (v *SysctlValidationCheck) NeedFix() bool {
	return true
}

func (v *SysctlValidationCheck) Fix(ctx context.Context) error {
	return nil
}

func NewSysctlValidationCheck() ValidationCheck {
	return &SysctlValidationCheck{}
}
