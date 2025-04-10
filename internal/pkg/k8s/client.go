// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/exmap"
	"golang.org/x/term"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/easysoft/qcadmin/common"

	quchengclientset "github.com/easysoft/quickon-api/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	kubeerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	//EvictionKind EvictionKind
	EvictionKind = "Eviction"
	//EvictionSubresource EvictionSubresource
	EvictionSubresource = "pods/eviction"
)

type Client struct {
	Clientset        kubernetes.Interface
	DynamicClientset dynamic.Interface
	MetricsClientset *metrics.Clientset
	Config           *rest.Config
	RawConfig        clientcmdapi.Config
	restClientGetter genericclioptions.RESTClientGetter
	contextName      string
	QClient          *quchengclientset.Clientset
}

func NewSimpleQClient() (*Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = common.GetKubeConfig()
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	qclient, err := quchengclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		Clientset: client,
		Config:    config,
		QClient:   qclient,
	}, nil
}

func NewSimpleClient(kubecfg ...string) (*Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		if len(kubecfg) > 0 {
			kubeconfig = kubecfg[0]
		} else {
			kubeconfig = common.GetKubeConfig()
		}
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		Clientset: client,
		Config:    config,
	}, nil
}

func NewClient(contextName, kubeconfig string) (*Client, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = common.GetKubeConfig()
		}
	}
	restClientGetter := genericclioptions.ConfigFlags{
		Context:    &contextName,
		KubeConfig: &kubeconfig,
	}
	rawKubeConfigLoader := restClientGetter.ToRawKubeConfigLoader()

	config, err := rawKubeConfigLoader.ClientConfig()
	if err != nil {
		return nil, err
	}

	rawConfig, err := rawKubeConfigLoader.RawConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	if contextName == "" {
		contextName = rawConfig.CurrentContext
	}

	return &Client{
		Clientset:        clientset,
		Config:           config,
		DynamicClientset: dynamicClientset,
		MetricsClientset: metricsClient,
		RawConfig:        rawConfig,
		restClientGetter: &restClientGetter,
		contextName:      contextName,
	}, nil
}

func (c *Client) ContextName() (name string) {
	return c.contextName
}

func (c *Client) GetVersion(_ context.Context) (string, error) {
	v, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return "", errors.Errorf("failed to get Kubernetes version: %w", err)
	}
	return fmt.Sprintf("%#v", *v), nil
}

func (c *Client) GetGitVersion(_ context.Context) (string, error) {
	v, err := c.Clientset.Discovery().ServerVersion()
	if err != nil {
		return "", errors.Errorf("failed to get Kubernetes version: %w", err)
	}
	return v.GitVersion, nil
}

// SupportEviction uses Discovery API to find out if the server support eviction subresource
// If support, it will return its groupVersion; Otherwise, it will return ""
func (c *Client) SupportEviction() (string, error) {
	discoveryClient := c.Clientset.Discovery()
	groupList, err := discoveryClient.ServerGroups()
	if err != nil {
		return "", err
	}
	foundPolicyGroup := false
	var policyGroupVersion string
	for _, group := range groupList.Groups {
		if group.Name == "policy" {
			foundPolicyGroup = true
			policyGroupVersion = group.PreferredVersion.GroupVersion
			break
		}
	}
	if !foundPolicyGroup {
		return "", nil
	}
	resourceList, err := discoveryClient.ServerResourcesForGroupVersion("v1")
	if err != nil {
		return "", err
	}
	for _, resource := range resourceList.APIResources {
		if resource.Name == EvictionSubresource && resource.Kind == EvictionKind {
			return policyGroupVersion, nil
		}
	}
	return "", nil
}

func (c *Client) ListSC(ctx context.Context, opts metav1.ListOptions) (*storagev1.StorageClassList, error) {
	return c.Clientset.StorageV1().StorageClasses().List(ctx, opts)
}

func (c *Client) GetDefaultSC(ctx context.Context) (*storagev1.StorageClass, error) {
	scs, err := c.ListSC(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, sc := range scs.Items {
		defaultclass := exmap.GetLabelValue(sc.Annotations, "storageclass.kubernetes.io/is-default-class")
		if defaultclass == "true" {
			return &sc, nil
		}
	}
	return nil, errors.Errorf("no default storage class found")
}

func (c *Client) PatchDefaultSC(ctx context.Context, sc *storagev1.StorageClass, isDefaultSC bool) error {
	scan := sc.Annotations
	scan = exmap.MergeLabels(scan, map[string]string{
		"storageclass.kubernetes.io/is-default-class": fmt.Sprintf("%v", isDefaultSC),
	})
	sc.Annotations = scan
	_, err := c.Clientset.StorageV1().StorageClasses().Update(ctx, sc, metav1.UpdateOptions{})
	return err
}

func (c *Client) ListNodes(ctx context.Context, opts metav1.ListOptions) (*corev1.NodeList, error) {
	return c.Clientset.CoreV1().Nodes().List(ctx, opts)
}

func (c *Client) GetNodeByName(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Node, error) {
	return c.Clientset.CoreV1().Nodes().Get(ctx, name, opts)
}

func (c *Client) GetNodeByIP(ctx context.Context, ip string) (*corev1.Node, error) {
	nodes, err := c.Clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if len(nodes.Items) == 0 {
		return nil, kubeerr.NewNotFound(schema.GroupResource{Resource: "nodes"}, ip)
	}
	for _, node := range nodes.Items {
		for _, nodeip := range node.Status.Addresses {
			if nodeip.Address == ip {
				return &node, nil
			}
		}
	}
	return nil, kubeerr.NewNotFound(schema.GroupResource{Resource: "nodes"}, ip)
}

func (c *Client) DownNode(ctx context.Context, ip string) error {
	nodeinfo, err := c.GetNodeByIP(ctx, ip)
	if err != nil {
		if kubeerr.IsNotFound(err) {
			return nil
		}
		return err
	}
	// TODO 需要处理已存在
	c.CordonOrUnCordonNode(ctx, nodeinfo.Name, true, metav1.PatchOptions{})
	if err := c.DeleteOrEvictPodsSimple(ctx, nodeinfo.Name); err != nil {
		return err
	}
	return c.Clientset.CoreV1().Nodes().Delete(ctx, nodeinfo.Name, metav1.DeleteOptions{})
}

func (c *Client) DeleteNode(ctx context.Context, name string) error {
	return c.Clientset.CoreV1().Nodes().Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *Client) DrainNode(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Nodes().Delete(ctx, name, opts)
}

func (c *Client) CordonOrUnCordonNode(ctx context.Context, name string, drain bool, opts metav1.PatchOptions) (*corev1.Node, error) {
	data := fmt.Sprintf(`{"spec":{"unschedulable":%t}}`, drain)
	node, err := c.Clientset.CoreV1().Nodes().Patch(ctx, name, types.StrategicMergePatchType, []byte(data), opts)
	if err != nil {
		return node, err
	}
	return node, nil
}

func (c *Client) DeleteOrEvictPodsSimple(ctx context.Context, name string) error {
	pods, err := c.GetPodsByNodes(name)
	if err != nil {
		return err
	}
	policyGroupVersion, err := c.SupportEviction()
	if err != nil {
		return err
	}
	if policyGroupVersion == "" {
		return errors.New("kube api not support eviction subresource")
	}
	for _, v := range pods {
		c.EvictPod(ctx, v, policyGroupVersion)
	}
	return nil
}

// evictPod 驱离POD
func (c *Client) EvictPod(ctx context.Context, pod corev1.Pod, policyGroupVersion string) error {
	deleteOptions := &metav1.DeleteOptions{}
	//if o.GracePeriodSeconds >= 0 {
	//	gracePeriodSeconds := int64(o.GracePeriodSeconds)
	//	deleteOptions.GracePeriodSeconds = &gracePeriodSeconds
	//}
	eviction := &policyv1.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: policyGroupVersion,
			Kind:       EvictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: deleteOptions,
	}
	// Remember to change change the URL manipulation func when Evction's version change
	return c.Clientset.PolicyV1().Evictions(eviction.Namespace).Evict(ctx, eviction)
}

func (c *Client) CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Create(ctx, secret, opts)
}

func (c *Client) UpdateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.UpdateOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Update(ctx, secret, opts)
}

func (c *Client) PatchSecret(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) DeleteSecret(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Secrets(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetSecret(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Get(ctx, name, opts)
}

func (c *Client) CreateServiceAccount(ctx context.Context, namespace string, account *corev1.ServiceAccount, opts metav1.CreateOptions) (*corev1.ServiceAccount, error) {
	return c.Clientset.CoreV1().ServiceAccounts(namespace).Create(ctx, account, opts)
}

func (c *Client) DeleteServiceAccount(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateClusterRole(ctx context.Context, role *rbacv1.ClusterRole, opts metav1.CreateOptions) (*rbacv1.ClusterRole, error) {
	return c.Clientset.RbacV1().ClusterRoles().Create(ctx, role, opts)
}

func (c *Client) DeleteClusterRole(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().ClusterRoles().Delete(ctx, name, opts)
}

func (c *Client) CreateClusterRoleBinding(ctx context.Context, role *rbacv1.ClusterRoleBinding, opts metav1.CreateOptions) (*rbacv1.ClusterRoleBinding, error) {
	return c.Clientset.RbacV1().ClusterRoleBindings().Create(ctx, role, opts)
}

func (c *Client) DeleteClusterRoleBinding(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().ClusterRoleBindings().Delete(ctx, name, opts)
}

func (c *Client) GetConfigMap(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
}

func (c *Client) PatchConfigMap(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) UpdateConfigMap(ctx context.Context, configMap *corev1.ConfigMap, opts metav1.UpdateOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMap, opts)
}

func (c *Client) CreateConfigMap(ctx context.Context, namespace string, config *corev1.ConfigMap, opts metav1.CreateOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Create(ctx, config, opts)
}

func (c *Client) DeleteConfigMap(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateService(ctx context.Context, namespace string, service *corev1.Service, opts metav1.CreateOptions) (*corev1.Service, error) {
	return c.Clientset.CoreV1().Services(namespace).Create(ctx, service, opts)
}

func (c *Client) UpdateService(ctx context.Context, namespace string, service *corev1.Service, opts metav1.UpdateOptions) (*corev1.Service, error) {
	return c.Clientset.CoreV1().Services(namespace).Update(ctx, service, opts)
}

func (c *Client) DeleteService(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Services(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetService(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Service, error) {
	return c.Clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
}

func (c *Client) ListServices(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.ServiceList, error) {
	return c.Clientset.CoreV1().Services(namespace).List(ctx, options)
}

func (c *Client) CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, opts)
}

func (c *Client) GetDeployment(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

func (c *Client) DeleteDeployment(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}

func (c *Client) PatchDeployment(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) CheckDeploymentStatus(ctx context.Context, namespace, deployment string) error {
	d, err := c.GetDeployment(ctx, namespace, deployment, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if d == nil {
		return errors.Errorf("deployment is not available")
	}

	if d.Status.ObservedGeneration != d.Generation {
		return errors.Errorf("observed generation (%d) is older than generation of the desired state (%d)",
			d.Status.ObservedGeneration, d.Generation)
	}

	if d.Status.Replicas == 0 {
		return errors.Errorf("replicas count is zero")
	}

	if d.Status.AvailableReplicas != d.Status.Replicas {
		return errors.Errorf("only %d of %d replicas are available", d.Status.AvailableReplicas, d.Status.Replicas)
	}

	if d.Status.ReadyReplicas != d.Status.Replicas {
		return errors.Errorf("only %d of %d replicas are ready", d.Status.ReadyReplicas, d.Status.Replicas)
	}

	if d.Status.UpdatedReplicas != d.Status.Replicas {
		return errors.Errorf("only %d of %d replicas are up-to-date", d.Status.UpdatedReplicas, d.Status.Replicas)
	}

	return nil
}

func (c *Client) CreateDaemonSet(ctx context.Context, namespace string, ds *appsv1.DaemonSet, opts metav1.CreateOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Create(ctx, ds, opts)
}

func (c *Client) PatchDaemonSet(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) GetDaemonSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
}

func (c *Client) ListDaemonSet(ctx context.Context, namespace string, o metav1.ListOptions) (*appsv1.DaemonSetList, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).List(ctx, o)
}

func (c *Client) DeleteDaemonSet(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.AppsV1().DaemonSets(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateNamespace(ctx context.Context, namespace string, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	return c.Clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}, opts)
}

func (c *Client) GetNamespace(ctx context.Context, namespace string, options metav1.GetOptions) (*corev1.Namespace, error) {
	return c.Clientset.CoreV1().Namespaces().Get(ctx, namespace, options)
}

func (c *Client) DeleteNamespace(ctx context.Context, namespace string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Namespaces().Delete(ctx, namespace, opts)
}

func (c *Client) CheckNamespace(ctx context.Context, namespace string) error {
	_, err := c.GetNamespace(ctx, namespace, metav1.GetOptions{})
	if kubeerr.IsNotFound(err) {
		_, err = c.CreateNamespace(ctx, namespace, metav1.CreateOptions{})
		return err
	}
	return err
}

func (c *Client) ListNamespaces(ctx context.Context, o metav1.ListOptions) (*corev1.NamespaceList, error) {
	return c.Clientset.CoreV1().Namespaces().List(ctx, o)
}

func (c *Client) ListEvents(ctx context.Context, o metav1.ListOptions) (*corev1.EventList, error) {
	return c.Clientset.CoreV1().Events(corev1.NamespaceAll).List(ctx, o)
}

func (c *Client) DeletePod(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (c *Client) DeletePodCollection(ctx context.Context, namespace string, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return c.Clientset.CoreV1().Pods(namespace).DeleteCollection(ctx, opts, listOpts)
}

func (c *Client) ListPods(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.PodList, error) {
	return c.Clientset.CoreV1().Pods(namespace).List(ctx, options)
}

func (c *Client) PodLogs(namespace, name string, opts *corev1.PodLogOptions) *rest.Request {
	return c.Clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
}

func (c *Client) GetLogs(ctx context.Context, namespace, name, container string, sinceTime time.Time, limitBytes int64, previous bool) (string, error) {
	t := metav1.NewTime(sinceTime)
	o := corev1.PodLogOptions{
		Container:  container,
		Follow:     false,
		LimitBytes: &limitBytes,
		Previous:   previous,
		SinceTime:  &t,
		Timestamps: true,
	}
	r := c.Clientset.CoreV1().Pods(namespace).GetLogs(name, &o)
	s, err := r.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer s.Close()
	var b bytes.Buffer
	if _, err = io.Copy(&b, s); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (c *Client) GetFollowLogs(ctx context.Context, namespace, name, container string, previous bool) error {
	// t := metav1.NewTime(sinceTime)
	o := corev1.PodLogOptions{
		Container: container,
		Follow:    true,
		// LimitBytes: &limitBytes,
		Previous: previous,
		// SinceTime:  &t,
		Timestamps: true,
	}
	r := c.Clientset.CoreV1().Pods(namespace).GetLogs(name, &o)
	s, err := r.Stream(ctx)
	if err != nil {
		return err
	}
	defer s.Close()
	scanlogs := bufio.NewReader(s)
	for {
		bytes, err := scanlogs.ReadBytes('\n')
		fmt.Println(string(bytes))
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}

func (c *Client) ExecInPodWithStderr(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
	})
	return result.Stdout, result.Stderr, err
}

func (c *Client) ExecInPodWithTTY(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
		TTY:       true,
	})
	// Using TTY (for context cancellation support) fuses stderr into stdout
	return result.Stdout, err
}

func (c *Client) ExecInPod(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
	})
	if err != nil {
		return bytes.Buffer{}, err
	}

	if errString := result.Stderr.String(); errString != "" {
		return bytes.Buffer{}, errors.Errorf("command failed: %s", errString)
	}

	return result.Stdout, nil
}

func (c *Client) ExecPodWithTTY(ctx context.Context, namespace, podName, container string, command []string) error {
	req := c.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   command,
			Container: container,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(c.Config, "POST", req.URL())
	if err != nil {
		return err
	}
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		return errors.Errorf("stdin/stdout must be a terminal")
	}
	oldstate, _ := term.MakeRaw(0)
	defer term.Restore(0, oldstate)
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  screen,
		Stdout: screen,
		Stderr: screen,
		Tty:    true,
	})
	return err
}

func (c *Client) CreateIngressClass(ctx context.Context, ingressClass *networkingv1.IngressClass, opts metav1.CreateOptions) (*networkingv1.IngressClass, error) {
	return c.Clientset.NetworkingV1().IngressClasses().Create(ctx, ingressClass, opts)
}

func (c *Client) ListIngressClass(ctx context.Context, opts metav1.ListOptions) (*networkingv1.IngressClassList, error) {
	return c.Clientset.NetworkingV1().IngressClasses().List(ctx, opts)
}

func (c *Client) ListDefaultIngressClass(ctx context.Context, opts metav1.ListOptions) (*networkingv1.IngressClass, error) {
	lists, err := c.ListIngressClass(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, l := range lists.Items {
		if exmap.GetLabelValue(l.Annotations, "ingressclass.kubernetes.io/is-default-class") == "true" {
			return &l, nil
		}
	}
	return nil, errors.New("not found default ingress class")
}

func (c *Client) DeleteIngressClass(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.NetworkingV1().IngressClasses().Delete(ctx, name, opts)
}

func (c *Client) GetSecretKeyBySelector(ctx context.Context, namespace string, secretSelector *corev1.SecretKeySelector) (string, error) {
	secret, err := c.GetSecret(ctx, namespace, secretSelector.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	if data, ok := secret.Data[secretSelector.Key]; ok {
		return string(data), nil
	}
	return "", errors.Errorf("key %s not found in secret %s", secretSelector.Key, secretSelector.Name)
}

func (c *Client) ListIngress(ctx context.Context, namespace string, opts metav1.ListOptions) (*networkingv1.IngressList, error) {
	return c.Clientset.NetworkingV1().Ingresses(namespace).List(ctx, opts)
}
