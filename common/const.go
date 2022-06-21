// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import "time"

const (
	FileMode0755 = 0o755
	FileMode0644 = 0o644
	FileMode0600 = 0o600
)

const (
	DefaultDaemonPort = 60080
	DefaultLogDir     = ".qc/log"
	DefaultDataDir    = ".qc/data"
	DefaultBinDir     = ".qc/bin"
	DefaultCfgDir     = ".qc/config"
	DefaultCacheDir   = ".qc/cache"
)

const (
	// KubernetesDir is the directory Kubernetes owns for storing various configuration files
	KubernetesDir = "/etc/kubernetes"
	// ManifestsSubDirName defines directory name to store manifests
	ManifestsSubDirName = "manifests"
	// ControlPlaneNumCPU is the number of CPUs required on control-plane
	ControlPlaneNumCPU = 2
	// ControlPlaneMem is the number of megabytes of memory required on the control-plane
	// Below that amount of RAM running a stable control plane would be difficult.
	ControlPlaneMem = 1700
)

const (
	// CRISocketContainerd is the containerd CRI endpoint
	CRISocketContainerd = "unix:///var/run/containerd/containerd.sock"
	// CRISocketCRIO is the cri-o CRI endpoint
	CRISocketCRIO = "unix:///var/run/crio/crio.sock"
	// CRISocketCRIDocker is the cri-dockerd CRI endpoint
	CRISocketCRIDocker = "unix:///var/run/cri-dockerd.sock"
	// CRISocketDocker is the cri-dockerd CRI endpoint
	CRISocketDocker = "unix:///var/run/docker.sock"
	// DefaultCRISocket defines the default CRI socket
	DefaultCRISocket = CRISocketContainerd

	// StatusRunning instance running status.
	StatusRunning = "Running"
	// StatusCreating instance creating status.
	StatusCreating = "Creating"
	// StatusFailed instance failed status.
	StatusFailed = "Failed"
)

const (
	K3sBinName          = "k3s"
	K3sBinPath          = "/usr/local/bin/k3s"
	HelmBinName         = "helm"
	HelmBinPath         = "/usr/local/bin/helm"
	K3sBinVersion       = "v1.23.5+k3s1"
	K3sBinURL           = "https://github.com/k3s-io/k3s/releases/download"
	K3sAgentEnv         = "/etc/systemd/system/k3s-agent.service.env"
	K3sKubeConfig       = "/etc/rancher/k3s/k3s.yaml"
	KubeQPS             = 5.0
	KubeBurst           = 10
	KubectlBinPath      = "/usr/local/bin/kubectl"
	QcAdminBinPath      = "/usr/local/bin/qcadmin"
	StatusWaitDuration  = 5 * time.Minute
	WaitRetryInterval   = 5 * time.Second
	DefaultHelmRepoName = "install"
	DefaultSystem       = "cne-system"
	DefaultChartName    = "qucheng"
	DefaultAPIChartName = "cne-api"
	DefaultQuchengName  = "qucheng"
	DefaultCneAPIName   = "cne-api"
	DefaultDBName       = "qucheng-mysql"
	InitFileName        = ".initdone"
	InitLockFileName    = ".qlock"
	InitModeCluster     = ".incluster"
)

const (
	DownloadAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4688.0 Safari/537.36"
	QuchengDocs   = "https://www.qucheng.com"
)
