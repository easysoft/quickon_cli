// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import "time"

const (
	FileMode0777 = 0o777
	FileMode0755 = 0o755
	FileMode0644 = 0o644
	FileMode0600 = 0o600
)

const (
	DefaultDaemonPort     = 60080
	DefaultLogDir         = ".qc/log"
	DefaultDataDir        = ".qc/data"
	DefaultBinDir         = ".qc/bin"
	DefaultCfgDir         = ".qc/config"
	DefaultCacheDir       = ".qc/cache"
	DefaultQuickonDataDir = "/opt/quickon"
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
	ControlPlaneMem = 3700
	// ControlPlaneNumDisk is the number of Disk required on control-plane
	ControlPlaneNumDisk = 40
	// ControlPlaneLowDisk is the number of Disk required on node
	ControlPlaneLowDisk = 20
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
	DefaultQuickonOSSVersion      = "3.0.2" // 开源版本
	DefaultQuickonEEVersion       = "1.0"
	DefaultZentaoDevOPSOSSVersion = "20.4.0"
	DefaultZentaoDevOPSMaxVersion = "5.3.0"
	DefaultZentaoDevOPSBizVersion = "10.3.0"
	DefaultZentaoDevOPSIPDVersion = "2.3.0"
	K3sBinName                    = "k3s"
	K3sBinPath                    = "/usr/local/bin/k3s"
	HelmBinName                   = "helm"
	HelmBinPath                   = "/usr/local/bin/helm"
	K3sBinVersion                 = "v1.28.12+k3s1"
	K3sBinURL                     = "https://github.com/k3s-io/k3s/releases/download"
	K3sAgentEnv                   = "/etc/systemd/system/k3s-agent.service.env"
	K3sKubeConfig                 = "/etc/rancher/k3s/k3s.yaml"
	K3sDefaultDir                 = "/var/lib/rancher/k3s"
	KubeQPS                       = 5.0
	KubeBurst                     = 10
	KubectlBinPath                = "/usr/local/bin/kubectl"
	CRICrictl                     = "/usr/local/bin/crictl"
	CRICtr                        = "/usr/local/bin/ctr"
	QcAdminBinPath                = "/usr/bin/qcadmin"
	StatusWaitDuration            = 5 * time.Minute
	WaitRetryInterval             = 5 * time.Second
	DefaultKubeSystem             = "kube-system"
	DefaultHelmRepoName           = "install"
	DefaultSystem                 = "cne-system"
	DefaultSystemNamespace        = "quickon-system"
	DefaultStorageNamespace       = "quickon-storage"
	DefaultAppNamespace           = "quickon-app"
	DefaultCINamespace            = "quickon-ci"
	DefaultQuchengName            = "qucheng"
	DefaultZentaoPaasName         = "zentaopaas"
	DefaultCneOperatorName        = "cne-operator"
	DefaultIngressName            = "nginx-ingress-controller"
	DefaultKubeAPIDomain          = "kubeapi.quickon.local"
	DefaultStorageType            = "nfs"
	InitFileName                  = ".initdone"
	InitLockFileName              = ".qlock"
	InitModeCluster               = ".incluster"
	DefaultOSUserRoot             = "root"
	DefaultSSHPort                = 22
	PrivateKeyFilename            = "id_rsa"
	PublicKeyFilename             = "id_rsa.pub"
	DefaultGenerate204URL         = "https://connectivitycheck.gstatic.com/generate_204"
	MiuiGenerate204URL            = "https://connect.rom.miui.com/generate_204"
	V2exGenerate204URL            = "https://captive.v2ex.co/generate_204"
	CloudflareEdgeTraceURL        = "https://www.cloudflare.com/cdn-cgi/trace"
	CneAPITokenHeader             = "X-Auth-Token"
	CDNVersionURL                 = "https://pkg.zentao.net/cli/devops/version.json"
)

const (
	DownloadAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0"
	QuchengDocs        = "https://www.qucheng.com"
	ZentaoDocs         = "https://www.zentao.net/book/devops.html"
	QuchengDefaultUser = "qadmin"
	QuchengDefaultPass = "pass4Quickon"
)

const K3SServiceTpl = `
[Unit]
Description=Lightweight Kubernetes
Documentation=https://k3s.io
Wants=network-online.target
After=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
{{if .TypeMaster -}}
Type=notify
{{else}}
Type=exec
{{ end -}}
EnvironmentFile=-/etc/default/%N
EnvironmentFile=-/etc/sysconfig/%N
EnvironmentFile=-/etc/systemd/system/k3s.service.env
KillMode=process
Delegate=yes
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
ExecStartPre=/bin/sh -xc '! /usr/bin/systemctl is-enabled --quiet nm-cloud-setup.service'
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/k3s \
    {{if .TypeMaster -}}
      server \
      --tls-san kubeapi.corp.cc \
      --tls-san apiserver.cluster.local \
      --tls-san {{ .KubeAPI }} \
      {{if .PodCIDR -}}
        --cluster-cidr {{ .PodCIDR }} \
      {{end -}}
      {{if .ServiceCIDR -}}
        --service-cidr {{ .ServiceCIDR }} \
      {{end -}}
      {{if .DataStore -}}
        --datastore-endpoint {{ .DataStore }} \
      {{else -}}
        --cluster-init \
        --etcd-expose-metrics \
        --etcd-snapshot-name auto-snapshot \
        --etcd-snapshot-compress \
      {{end -}}
      {{if .LocalStorage -}}
        --disable servicelb,traefik \
      {{else -}}
        --disable servicelb,traefik,local-storage \
      {{end -}}
      --disable-cloud-controller \
      --system-default-registry {{ .Registry }} \
      --disable-network-policy \
      --disable-helm-controller \
      {{if .CNI -}}
        --flannel-backend {{ .CNI }} \
      {{end -}}
    {{else -}}
      --server https://{{ .KubeAPI }}:6443 \
      agent \
    {{end -}}
      --token {{ .KubeToken }} \
      --data-dir {{.DataDir}} \
      --docker \
      --pause-image {{ .Registry }}/rancher/mirrored-pause:3.6 \
      --prefer-bundled-bin \
      --kube-proxy-arg "proxy-mode=ipvs" "masquerade-all=true" \
      --kube-proxy-arg "metrics-bind-address=0.0.0.0"
`
