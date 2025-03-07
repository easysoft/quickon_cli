/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package preflight

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/zos"
	"github.com/shirou/gopsutil/v3/disk"
	"k8s.io/apimachinery/pkg/util/validation"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/types"
	"github.com/easysoft/qcadmin/internal/pkg/util/initsystem"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"

	qnetutil "github.com/easysoft/qcadmin/internal/pkg/util/netutil"
	netutil "k8s.io/apimachinery/pkg/util/net"
	system "k8s.io/system-validators/validators"
	utilsexec "k8s.io/utils/exec"
	netutils "k8s.io/utils/net"
)

const (
	bridgenf = "/proc/sys/net/bridge/bridge-nf-call-iptables"
	// bridgenf6             = "/proc/sys/net/bridge/bridge-nf-call-ip6tables"
	ipv4Forward = "/proc/sys/net/ipv4/ip_forward"
	// ipv6DefaultForwarding = "/proc/sys/net/ipv6/conf/default/forwarding"
)

// Checker validates the state of the system to ensure kubeadm will be
// successful as often as possible.
type Checker interface {
	Check() error
	Name() string
}

// ServiceCheck verifies that the given service is enabled and active. If we do not
// detect a supported init system however, all checks are skipped and a warning is
// returned.
type ServiceCheck struct {
	Service       string
	CheckIfActive bool
	CheckIfExist  bool
}

// Name returns label for ServiceCheck. If not provided, will return based on the service parameter
func (sc ServiceCheck) Name() string {
	return fmt.Sprintf("Service-%s", strings.ToUpper(sc.Service))
}

// Check validates if the service is enabled and active.
func (sc ServiceCheck) Check() error {
	log := log.GetInstance()
	log.Debugf("validating if the %q service is existed or active", sc.Service)
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		return err
	}
	if !initSystem.ServiceExists(sc.Service) {
		return errors.Errorf("%s service does not exist", sc.Service)
	}

	if sc.CheckIfExist {
		return nil
	}

	if !initSystem.ServiceIsEnabled(sc.Service) {
		return errors.Errorf("%s service is not enabled, please run '%s'",
			sc.Service, initSystem.EnableCommand(sc.Service))
	}

	if sc.CheckIfActive && !initSystem.ServiceIsActive(sc.Service) {
		return errors.Errorf("%s service is not active, please run 'systemctl start %s.service'",
			sc.Service, sc.Service)
	}

	return nil
}

// FirewalldCheck checks if firewalld is enabled or active. If it is, warn the user that there may be problems
// if no actions are taken.
type FirewalldCheck struct {
	ports []int
}

// Name returns label for FirewalldCheck.
func (FirewalldCheck) Name() string {
	return "Firewalld"
}

// Check validates if the firewall is enabled and active.
func (fc FirewalldCheck) Check() error {
	log := log.GetInstance()
	log.Debug("validating if the firewall is enabled and active")
	initSystem, err := initsystem.GetInitSystem()
	if err != nil {
		return err
	}

	if !initSystem.ServiceExists("firewalld") {
		return nil
	}

	if initSystem.ServiceIsActive("firewalld") {
		log.Warnf("firewalld is active, please ensure ports %v are open or your cluster may not function correctly", fc.ports)
	}
	return nil
}

// PortOpenCheck ensures the given port is available for use.
type PortOpenCheck struct {
	port  int
	label string
}

// Name returns name for PortOpenCheck. If not known, will return "PortXXXX" based on port number
func (poc PortOpenCheck) Name() string {
	if poc.label != "" {
		return poc.label
	}
	return fmt.Sprintf("Port-%d", poc.port)
}

// Check validates if the particular port is available.
func (poc PortOpenCheck) Check() error {
	log := log.GetInstance()
	log.Debugf("validating availability of port %d", poc.port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", poc.port))
	if err != nil {
		return errors.Errorf("Port %d is in use", poc.port)
	}
	if ln != nil {
		if err = ln.Close(); err != nil {
			return errors.Errorf("when closing port %d, encountered %v", poc.port, err)
		}
	}
	return nil
}

// IsPrivilegedUserCheck verifies user is privileged (linux - root, windows - Administrator)
type IsPrivilegedUserCheck struct{}

// Name returns name for IsPrivilegedUserCheck
func (IsPrivilegedUserCheck) Name() string {
	return "IsPrivilegedUser"
}

type NetworkCheck struct {
	offline bool
}

// Name returns name for NetworkCheck
func (NetworkCheck) Name() string {
	return "Network"
}

// Check validates if the user is privileged.
func (nc NetworkCheck) Check() error {
	log := log.GetInstance()
	if gw, err := qnetutil.CheckDefaultRoute(); err == nil {
		log.Donef("default route %s reachable via icmp", color.SGreen(gw.String()))
	} else {
		log.Warnf("default route %s unreachable via icmp", color.SRed(gw.String()))
	}
	if ns, err := qnetutil.GetDefaultNameserver(); err == nil {
		a := "reachable"
		if !qnetutil.CheckReachabilityWithICMP(ns) {
			a = "unreachable"
		}
		if err := qnetutil.CheckNameserverAvailability(ns + ":53"); err != nil {
			log.Warnf("nameserver %s (ICMP %s, DNS unreachable: %s)", color.SRed(ns), a, err)
		} else {
			log.Donef("nameserver %s (ICMP %s, DNS reachable)", color.SGreen(ns), a)
		}
	} else {
		log.Warnf("failed to reading default nameserver from system: %s", err)
	}
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		if nc.offline {
			log.Infof("skipping captive portal %s detection", color.SGreen(common.DefaultGenerate204URL))
		} else {
			if err := qnetutil.CheckCaptivePortal(); err == nil {
				log.Donef("captive portal %s detected success", color.SGreen(common.DefaultGenerate204URL))
			} else {
				log.Warnf("captive portal %s detected failed: %s", color.SRed(common.DefaultGenerate204URL), err)
			}
		}
	}()
	go func() {
		defer wg.Done()
		if nc.offline {
			log.Infof("skipping captive portal %s detection", color.SGreen(common.MiuiGenerate204URL))
		} else {
			if err := qnetutil.CheckCaptivePortal(common.MiuiGenerate204URL); err == nil {
				log.Donef("captive portal %s detected success", color.SGreen(common.MiuiGenerate204URL))
			} else {
				log.Warnf("captive portal %s detected failed: %s", color.SRed(common.MiuiGenerate204URL), err)
			}
		}
	}()
	go func() {
		defer wg.Done()
		if nc.offline {
			log.Infof("skipping detection remote dns %s", color.SGreen("119.29.29.29"))
		} else {
			if err := qnetutil.CheckNameserverAvailability("119.29.29.29:53"); err != nil {
				log.Warnf("remote dns %s is unavailable: %s", color.SRed("119.29.29.29"), err)
			} else {
				log.Donef("remote dns %s is available", color.SGreen("119.29.29.29"))
			}
		}
	}()
	go func() {
		defer wg.Done()
		if nc.offline {
			log.Infof("skipping detection remote dns %s", color.SGreen("1.2.4.8"))
		} else {
			if err := qnetutil.CheckNameserverAvailability("1.2.4.8:53"); err != nil {
				log.Warnf("remote dns %s is unavailable: %s", color.SRed("1.2.4.8"), err)
			} else {
				log.Donef("remote dns %s is available", color.SGreen("1.2.4.8"))
			}
		}
	}()
	go func() {
		defer wg.Done()
		if nc.offline {
			log.Infof("skipping match cdn edge trace")
		} else {
			if loc, err := qnetutil.GetCloudflareEdgeTrace(); err == nil {
				log.Donef("match cdn edge trace loc %s", color.SGreen(loc))
			} else {
				log.Warnf("miss cdn edge trace")
			}
		}
	}()
	wg.Wait()
	return nil
}

// DirAvailableCheck checks if the given directory either does not exist, or is empty.
type DirAvailableCheck struct {
	Path  string
	Label string
}

// Name returns label for individual DirAvailableChecks. If not known, will return based on path.
func (dac DirAvailableCheck) Name() string {
	if dac.Label != "" {
		return dac.Label
	}
	return fmt.Sprintf("DirAvailable-%s", strings.Replace(dac.Path, "/", "-", -1))
}

// Check validates if a directory does not exist or empty.
func (dac DirAvailableCheck) Check() error {
	log := log.GetInstance()
	log.Debugf("validating the existence and emptiness of directory %s", dac.Path)

	// If it doesn't exist we are good:
	if _, err := os.Stat(dac.Path); os.IsNotExist(err) {
		return nil
	}

	f, err := os.Open(dac.Path)
	if err != nil {
		return errors.Wrapf(err, "unable to check if %s is empty", dac.Path)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != io.EOF {
		return errors.Errorf("%s is not empty", dac.Path)
	}
	return nil
}

// FileAvailableCheck checks that the given file does not already exist.
type FileAvailableCheck struct {
	Path  string
	Label string
}

// Name returns label for individual FileAvailableChecks. If not known, will return based on path.
func (fac FileAvailableCheck) Name() string {
	if fac.Label != "" {
		return fac.Label
	}
	return fmt.Sprintf("FileAvailable-%s", strings.Replace(fac.Path, "/", "-", -1))
}

// Check validates if the given file does not already exist.
func (fac FileAvailableCheck) Check() error {
	log := log.GetInstance()
	log.Debugf("validating the existence of file %s", fac.Path)

	if _, err := os.Stat(fac.Path); err == nil {
		return errors.Errorf("%s already exists", fac.Path)
	}
	return nil
}

// FileExistingCheck checks that the given file does not already exist.
type FileExistingCheck struct {
	Path  string
	Label string
}

// Name returns label for individual FileExistingChecks. If not known, will return based on path.
func (fac FileExistingCheck) Name() string {
	if fac.Label != "" {
		return fac.Label
	}
	return fmt.Sprintf("FileExisting-%s", strings.Replace(fac.Path, "/", "-", -1))
}

// Check validates if the given file already exists.
func (fac FileExistingCheck) Check() error {
	log := log.GetInstance()
	log.Debugf("validating the existence of file %s", fac.Path)

	if _, err := os.Stat(fac.Path); err != nil {
		return errors.Errorf("%s doesn't exist", fac.Path)
	}
	return nil
}

// FileContentCheck checks that the given file contains the string Content.
type FileContentCheck struct {
	Path    string
	Content []byte
	Label   string
}

// Name returns label for individual FileContentChecks. If not known, will return based on path.
func (fcc FileContentCheck) Name() string {
	if fcc.Label != "" {
		return fcc.Label
	}
	return fmt.Sprintf("FileContent-%s", strings.Replace(fcc.Path, "/", "-", -1))
}

// Check validates if the given file contains the given content.
func (fcc FileContentCheck) Check() error {
	log := log.GetInstance()
	log.Debugf("validating the contents of file %s", fcc.Path)
	f, err := os.Open(fcc.Path)
	if err != nil {
		log.Warnf("%s does not exist", fcc.Path)
		return nil
	}

	lr := io.LimitReader(f, int64(len(fcc.Content)))
	defer f.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, lr)
	if err != nil {
		return errors.Errorf("%s could not be read", fcc.Path)
	}

	if !bytes.Equal(buf.Bytes(), fcc.Content) {
		log.Warnf("%s contents are not set to %s", fcc.Path, fcc.Content)
	}
	return nil
}

// InPathCheck checks if the given executable is present in $PATH
type InPathCheck struct {
	executable string
	mandatory  bool
	exec       utilsexec.Interface
	label      string
	suggestion string
}

// Name returns label for individual InPathCheck. If not known, will return based on path.
func (ipc InPathCheck) Name() string {
	if ipc.label != "" {
		return ipc.label
	}
	return fmt.Sprintf("FileExisting-%s", strings.Replace(ipc.executable, "/", "-", -1))
}

// Check validates if the given executable is present in the path.
func (ipc InPathCheck) Check() error {
	log := log.GetInstance()
	log.Debugf("validating the presence of executable %s", ipc.executable)
	_, err := ipc.exec.LookPath(ipc.executable)
	if err != nil {
		if ipc.mandatory {
			// Return as an error:
			return errors.Errorf("%s not found in system path", ipc.executable)
		}
		// Return as a warning:
		warningMessage := fmt.Sprintf("%s not found in system path", ipc.executable)
		if ipc.suggestion != "" {
			warningMessage += fmt.Sprintf("\nSuggestion: %s", ipc.suggestion)
		}
		log.Warn(warningMessage)
		return nil
	}
	return nil
}

// HostnameCheck checks if hostname match dns sub domain regex.
// If hostname doesn't match this regex, kubelet will not launch static pods like kube-apiserver/kube-controller-manager and so on.
type HostnameCheck struct {
	nodeName string
}

// Name will return Hostname as name for HostnameCheck
func (HostnameCheck) Name() string {
	return "Hostname"
}

// Check validates if hostname match dns sub domain regex.
// Check hostname length and format
func (hc HostnameCheck) Check() error {
	log := log.GetInstance()
	log.Debug("checking whether the given node name is valid and reachable using net.LookupHost")
	for _, msg := range validation.IsQualifiedName(hc.nodeName) {
		log.Warnf("invalid node name format %q: %s", hc.nodeName, msg)
	}
	addr, err := net.LookupHost(hc.nodeName)
	if addr == nil {
		log.Warnf("hostname \"%s\" could not be reached", hc.nodeName)
	}
	if err != nil {
		log.Warnf("hostname \"%s\", err: %v", hc.nodeName, err)
	}
	return nil
}

// HTTPProxyCheck checks if https connection to specific host is going
// to be done directly or over proxy. If proxy detected, it will return warning.
type HTTPProxyCheck struct {
	Proto string
	Host  string
}

// Name returns HTTPProxy as name for HTTPProxyCheck
func (hst HTTPProxyCheck) Name() string {
	return "HTTPProxy"
}

// Check validates http connectivity type, direct or via proxy.
func (hst HTTPProxyCheck) Check() error {
	log := log.GetInstance()
	log.Debug("validating if the connectivity type is via proxy or direct")
	u := &url.URL{Scheme: hst.Proto, Host: hst.Host}
	if netutils.IsIPv6String(hst.Host) {
		u.Host = net.JoinHostPort(hst.Host, "1234")
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	proxy, err := netutil.SetOldTransportDefaults(&http.Transport{}).Proxy(req)
	if err != nil {
		return err
	}
	if proxy != nil {
		return errors.Errorf("Connection to %q uses proxy %q. If that is not intended, adjust your proxy settings", u, proxy)
	}
	return nil
}

// HTTPProxyCIDRCheck checks if https connection to specific subnet is going
// to be done directly or over proxy. If proxy detected, it will return warning.
// Similar to HTTPProxyCheck above, but operates with subnets and uses API
// machinery transport defaults to simulate kube-apiserver accessing cluster
// services and pods.
type HTTPProxyCIDRCheck struct {
	Proto string
	CIDR  string
}

// Name will return HTTPProxyCIDR as name for HTTPProxyCIDRCheck
func (HTTPProxyCIDRCheck) Name() string {
	return "HTTPProxyCIDR"
}

// Check validates http connectivity to first IP address in the CIDR.
// If it is not directly connected and goes via proxy it will produce warning.
func (subnet HTTPProxyCIDRCheck) Check() error {
	log := log.GetInstance()
	log.Debug("validating http connectivity to first IP address in the CIDR")
	if len(subnet.CIDR) == 0 {
		return nil
	}

	_, cidr, err := netutils.ParseCIDRSloppy(subnet.CIDR)
	if err != nil {
		return errors.Wrapf(err, "error parsing CIDR %q", subnet.CIDR)
	}

	testIP, err := netutils.GetIndexedIP(cidr, 1)
	if err != nil {
		return errors.Wrapf(err, "unable to get first IP address from the given CIDR (%s)", cidr.String())
	}

	testIPstring := testIP.String()
	if len(testIP) == net.IPv6len {
		testIPstring = fmt.Sprintf("[%s]:1234", testIP)
	}
	url := fmt.Sprintf("%s://%s/", subnet.Proto, testIPstring)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Utilize same transport defaults as it will be used by API server
	proxy, err := netutil.SetOldTransportDefaults(&http.Transport{}).Proxy(req)
	if err != nil {
		return err
	}
	if proxy != nil {
		log.Warnf("connection to %q uses proxy %q. This may lead to malfunctional cluster setup. Make sure that Pod and Services IP ranges specified correctly as exceptions in proxy configuration", subnet.CIDR, proxy)
	}
	return nil
}

// SystemVerificationCheck defines struct used for running the system verification node check in test/e2e_node/system
type SystemVerificationCheck struct{}

// Name will return SystemVerification as name for SystemVerificationCheck
func (SystemVerificationCheck) Name() string {
	return "SystemVerification"
}

// Check runs all individual checks
func (sysver SystemVerificationCheck) Check() error {
	log := log.GetInstance()
	log.Debug("running all checks")
	// Create a buffered writer and choose a quite large value (1M) and suppose the output from the system verification test won't exceed the limit
	// Run the system verification check, but write to out buffered writer instead of stdout
	bufw := bufio.NewWriterSize(os.Stdout, 1*1024*1024)
	reporter := &system.StreamReporter{WriteStream: bufw}

	var errs []error
	// All the common validators we'd like to run:
	var validators = []system.Validator{
		&system.KernelValidator{Reporter: reporter}}

	if runtime.GOOS == "linux" {
		//add linux validators
		validators = append(validators,
			&system.OSValidator{Reporter: reporter},
			&system.CgroupsValidator{Reporter: reporter})
	}

	// Run all validators
	for _, v := range validators {
		warn, err := v.Validate(system.DefaultSysSpec)
		if err != nil {
			errs = append(errs, err...)
			log.Error(err)
		}
		if warn != nil {
			log.Warn(warn)
		}
	}

	if len(errs) != 0 {
		// Only print the output from the system verification check if the check failed
		log.Error("[preflight] The system verification failed. Printing the output from the verification:")
		bufw.Flush()
		return errors.Errorf("system verification failed")
	}
	return nil
}

// SwapCheck warns if swap is enabled
type SwapCheck struct{}

// Name will return Swap as name for SwapCheck
func (SwapCheck) Name() string {
	return "Swap"
}

// Check validates whether swap is enabled or not
func (swc SwapCheck) Check() error {
	log := log.GetInstance()
	log.Debug("validating whether swap is enabled or not")
	f, err := os.Open("/proc/swaps")
	if err != nil {
		// /proc/swaps not available, thus no reasons to warn
		return nil
	}
	defer f.Close()
	var buf []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Warnf("error reading /proc/swaps: %v", err)
		return nil
	}

	if len(buf) > 1 {
		log.Warnf("swap is enabled; production deployments should disable swap unless testing the NodeSwap feature gate of the kubelet")
		return nil
	}
	return nil
}

// NumCPUCheck checks if current number of CPUs is not less than required
type NumCPUCheck struct {
	NumCPU int
	Devops bool
}

// Name returns the label for NumCPUCheck
func (NumCPUCheck) Name() string {
	return "NumCPU"
}

// Check number of CPUs required by qcadmin
func (ncc NumCPUCheck) Check() error {
	log := log.GetInstance()
	log.Debug("validating number of CPUs")
	numCPU := runtime.NumCPU()
	if numCPU < ncc.NumCPU {
		return errors.Errorf("the number of available CPUs %d is less than the required %d", numCPU, ncc.NumCPU)
	}
	if ncc.Devops {
		if numCPU < ncc.NumCPU*2 {
			log.Warnf("the number of available CPUs %d is greater than the minimum %d, less than the recommended %d", numCPU, ncc.NumCPU, ncc.NumCPU*2)
			return nil
		}
		log.Donef("the number of available CPUs %d is greater than the required %d", numCPU, ncc.NumCPU*2)
	} else {
		log.Donef("the number of available CPUs %d is greater than the required %d", numCPU, ncc.NumCPU)
	}
	return nil
}

// NumDiskCheck checks if current number of Disk is not less than required
type NumDiskCheck struct {
	NumDisk int
	LowDisk int
}

// Name returns the label for NumDiskCheck
func (NumDiskCheck) Name() string {
	return "NumDisk"
}

// Check number of Disk required by qcadmin
func (ndc NumDiskCheck) Check() error {
	log := log.GetInstance()
	log.Debug("validating number of Disk")
	mountRoot, err := disk.Usage("/")
	if err != nil || mountRoot == nil {
		return errors.Errorf("disk check failed, reason: %v", err)
	}
	numDisk := float64(mountRoot.Total) / 1024.0 / 1024.0 / 1024.0
	if numDisk < float64(ndc.LowDisk) {
		return errors.Errorf("the number of available Disk %.2f GB is less than the required %d GB", numDisk, ndc.LowDisk)
	}
	if numDisk < float64(ndc.NumDisk) {
		log.Warnf("the number of available Disk %.2f GB, suggest the number of available Disk greater than the required %d GB", numDisk, ndc.NumDisk)
	} else {
		log.Donef("the number of available Disk %.2f GB is greater than the required %d GB", numDisk, ndc.NumDisk)
	}
	return nil
}

// MemCheck checks if the number of megabytes of memory is not less than required
type MemCheck struct {
	Mem    uint64
	Devops bool
}

// Name returns the label for memory
func (MemCheck) Name() string {
	return "Mem"
}

// RunInitNodeChecks executes all individual, applicable to control-plane node checks.
// The boolean flag 'isSecondaryControlPlane' controls whether we are running checks in a --join-control-plane scenario.
// The boolean flag 'downloadCerts' controls whether we should skip checks on certificates because we are downloading them.
// If the flag is set to true we should skip checks already executed by RunJoinNodeChecks.
func RunInitNodeChecks(execer utilsexec.Interface, cfg *types.Metadata, ignorePreflightErrors, offline, devops bool) error {
	log := log.GetInstance()
	if err := RunRootCheckOnly(ignorePreflightErrors); err != nil {
		return err
	}

	if err := RunKubeOnly(ignorePreflightErrors); err != nil {
		return err
	}

	// manifestsDir := filepath.Join(kubeadmconstants.KubernetesDir, kubeadmconstants.ManifestsSubDirName)
	checks := []Checker{
		NumCPUCheck{NumCPU: common.ControlPlaneNumCPU, Devops: devops},
		NumDiskCheck{NumDisk: common.ControlPlaneNumDisk, LowDisk: common.ControlPlaneLowDisk},
		// Linux only
		// TODO: support other OS, if control-plane is supported on it.
		MemCheck{Mem: common.ControlPlaneMem, Devops: devops},
		FirewalldCheck{ports: []int{80, 443, 6443, 32379, 32380}},
		PortOpenCheck{port: 80},
		PortOpenCheck{port: 443},
		PortOpenCheck{port: 6443},
		PortOpenCheck{port: 32379},
		PortOpenCheck{port: 32380},
		NetworkCheck{offline: offline},
		// FileAvailableCheck{Path: kubeadmconstants.GetStaticPodFilepath(kubeadmconstants.KubeAPIServer, manifestsDir)},
		// FileAvailableCheck{Path: kubeadmconstants.GetStaticPodFilepath(kubeadmconstants.KubeControllerManager, manifestsDir)},
		// FileAvailableCheck{Path: kubeadmconstants.GetStaticPodFilepath(kubeadmconstants.KubeScheduler, manifestsDir)},
		// FileAvailableCheck{Path: kubeadmconstants.GetStaticPodFilepath(kubeadmconstants.Etcd, manifestsDir)},
		// HTTPProxyCheck{Proto: "https", Host: cfg.LocalAPIEndpoint.AdvertiseAddress},
	}
	cidrs := strings.Split(cfg.ServiceCidr, ",")
	for _, cidr := range cidrs {
		checks = append(checks, HTTPProxyCIDRCheck{Proto: "https", CIDR: cidr})
	}
	cidrs = strings.Split(cfg.ClusterCidr, ",")
	for _, cidr := range cidrs {
		checks = append(checks, HTTPProxyCIDRCheck{Proto: "https", CIDR: cidr})
	}

	// non-windows checks
	if runtime.GOOS == "linux" {
		checks = append(checks,
			FileContentCheck{Path: bridgenf, Content: []byte{'1'}},
			FileContentCheck{Path: ipv4Forward, Content: []byte{'1'}},
			SwapCheck{},
			InPathCheck{executable: "crictl", mandatory: false, exec: execer},
			InPathCheck{executable: "conntrack", mandatory: false, exec: execer},
			InPathCheck{executable: "ip", mandatory: true, exec: execer},
			InPathCheck{executable: "iptables", mandatory: false, exec: execer},
			InPathCheck{executable: "mount", mandatory: true, exec: execer},
			InPathCheck{executable: "nsenter", mandatory: true, exec: execer},
			InPathCheck{executable: "ebtables", mandatory: false, exec: execer},
			InPathCheck{executable: "ethtool", mandatory: false, exec: execer},
			InPathCheck{executable: "socat", mandatory: false, exec: execer},
			InPathCheck{executable: "tc", mandatory: false, exec: execer},
			InPathCheck{executable: "touch", mandatory: true, exec: execer},
			InPathCheck{executable: "route", mandatory: false, exec: execer},
			InPathCheck{executable: "wget", mandatory: false, exec: execer},
			InPathCheck{executable: "curl", mandatory: false, exec: execer})
	}
	checks = append(checks,
		SystemVerificationCheck{},
		HostnameCheck{nodeName: zos.GetHostname()},
		// ServiceCheck{Service: "kubelet", CheckIfActive: false},
		// PortOpenCheck{port: kubeadmconstants.KubeletPort}
	)
	if err := RunChecks(checks, os.Stderr, ignorePreflightErrors); err != nil {
		return err
	}
	log.Done("preflight checks passed")
	return nil
}

// RunRootCheckOnly initializes checks slice of structs and call RunChecks
func RunRootCheckOnly(ignorePreflightErrors bool) error {
	checks := []Checker{
		IsPrivilegedUserCheck{},
	}

	return RunChecks(checks, os.Stderr, ignorePreflightErrors)
}

type KubeExistCheck struct{}

// Name returns name for KubeExistCheck
func (KubeExistCheck) Name() string {
	return "KubeExist"
}

// Check Exist kubernetes directory
func (kc KubeExistCheck) Check() error {
	log := log.GetInstance()
	log.Debug("validating kubernetes exist")
	return nil
}

// RunKubeOnly initializes checks slice of structs and call RunChecks
func RunKubeOnly(ignorePreflightErrors bool) error {
	checks := []Checker{
		KubeExistCheck{},
	}

	return RunChecks(checks, os.Stderr, ignorePreflightErrors)
}

// RunChecks runs each check, displays it's warnings/errors, and once all
// are processed will exit if any errors occurred.
func RunChecks(checks []Checker, ww io.Writer, ignorePreflightErrors bool) error {
	log := log.GetInstance()
	for _, c := range checks {
		// name := c.Name()
		if err := c.Check(); err != nil {
			if ignorePreflightErrors {
				log.Errorf("%s check err, reason: %v", c.Name(), err)
				continue
			}
			return err
		}
	}
	return nil
}
