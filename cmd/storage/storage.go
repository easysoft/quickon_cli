// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package storage

import (
	"context"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/easysoft/qcadmin/internal/pkg/util/log/survey"

	qcexec "github.com/easysoft/qcadmin/internal/pkg/util/exec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	nfsExample = templates.Examples(`
		# deploy local nfs storage
		z cluster storage nfs
		# deploy qcloud cfs v3
		z cluster storage nfs --ip cfsip --path cfspath
`)
)

// NewCmdStorage returns a cobra command for `storage` subcommands
func NewCmdStorage(f factory.Factory) *cobra.Command {
	s := &cobra.Command{
		Use:   "storage",
		Short: "storage commands",
		Long:  "manage cluster storage",
	}
	s.AddCommand(longhorn(f))
	s.AddCommand(nfs(f))
	s.AddCommand(local(f))
	s.AddCommand(defaultStorage(f))
	return s
}

func longhorn(f factory.Factory) *cobra.Command {
	logpkg := f.GetLog()
	cmd := &cobra.Command{
		Use:   "longhorn",
		Short: "deploy longhorn storage",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := qcexec.CommandRun("bash", "-c", common.GetCustomScripts("hack/manifests/storage/longhorn_environment_check.sh")); err != nil {
				return errors.Errorf("longhorn environment check failed, reason: %v", err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := qcexec.CommandRun("bash", "-c", common.GetCustomScripts("hack/manifests/storage/longhorn.sh")); err != nil {
				return errors.Errorf("longhorn install failed, reason: %v", err)
			}
			logpkg.Infof("install longhorn storage success")
			return nil
		},
	}
	return cmd
}

func local(f factory.Factory) *cobra.Command {
	logpkg := f.GetLog()
	cmd := &cobra.Command{
		Use:   "local",
		Short: "deploy local storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeargs := []string{"experimental", "kubectl", "apply", "-f", common.GetCustomScripts("hack/manifests/storage/local-storage.yaml")}
			output, err := qcexec.Command(os.Args[0], kubeargs...).CombinedOutput()
			if err != nil {
				logpkg.Errorf("upgrade install local storage failed: %s", string(output))
				return err
			}
			logpkg.Infof("install local storage class %s (%s) success", color.SGreen("q-local"), color.SGreen("/opt/quickon/storage/local"))
			return nil
		},
	}
	return cmd
}

func nfs(f factory.Factory) *cobra.Command {
	var ip, path, name string
	logpkg := f.GetLog()
	cmd := &cobra.Command{
		Use:     "nfs",
		Short:   "deploy nfs storage",
		Example: nfsExample,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(ip) == 0 && len(path) == 0 {
				an, err := logpkg.Question(&survey.QuestionOptions{
					Question:     "nfs server ip is empty, install local nfs",
					DefaultValue: "yes",
					Options:      []string{"yes", "no"},
				})
				if err != nil {
					return err
				}
				if an == "yes" {
					if err := qcexec.CommandRun("bash", "-c", common.GetCustomScripts("hack/manifests/storage/nfs-server.sh")); err != nil {
						return errors.Errorf("%s run install nfs script failed, reason: %v", ip, err)
					}
					ip = exnet.LocalIPs()[0]
					path = common.GetDefaultNFSStoragePath("")
					logpkg.Infof("install nfs server %s success", ip)
					return nil
				} else {
					return errors.Errorf("deny install local nfs, please set nfs server ip and path")
				}
			}
			if len(ip) > 0 && len(path) > 0 {
				if !exnet.CheckIP(ip) {
					return errors.Errorf("nfs server ip %s is invalid", ip)
				}
				return nil
			}
			return errors.Errorf("nfs server ip or path is empty")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO check install repo exist
			if err := qcexec.Command(os.Args[0], "experimental", "helm", "repo-init").Run(); err != nil {
				logpkg.Warnf("check helm repo failed: %v", err)
			}
			helmargs := []string{"experimental", "helm", "upgrade", "--name", name, "--repo", common.DefaultHelmRepoName, "--chart", "nfs-subdir-external-provisioner", "--namespace", common.DefaultStorageNamespace, "--set", "nfs.server=" + ip, "--set", "nfs.path=" + path, "--set", "storageClass.name=" + name}
			output, err := qcexec.Command(os.Args[0], helmargs...).CombinedOutput()
			if err != nil {
				logpkg.Errorf("upgrade install nfs failed: %s", string(output))
				return err
			}
			logpkg.Infof("install nfs storage class %s (%s:%s) success", color.SGreen(name), color.SGreen(ip), color.SGreen(path))
			return nil
		},
	}
	cmd.Flags().StringVar(&ip, "ip", "", "cloud cfs/nas ip")
	cmd.Flags().StringVar(&path, "path", "", "cloud cfs/nas path")
	cmd.Flags().StringVar(&name, "name", "q-nfs", "storage class name")
	return cmd
}

func defaultStorage(f factory.Factory) *cobra.Command {
	ds := &cobra.Command{
		Use:   "set-default",
		Short: "set default storage class",
		Long:  "set default storage class",
		RunE: func(cmd *cobra.Command, args []string) error {
			logpkg := f.GetLog()
			kubeClient, err := k8s.NewSimpleClient(common.GetKubeConfig())
			if err != nil {
				return errors.Errorf("kube client create failed, reason: %v", err)
			}
			ctx := context.Background()
			scs, err := kubeClient.ListSC(ctx, metav1.ListOptions{})
			if err != nil {
				return errors.Errorf("list storage class failed, reason: %v", err)
			}
			if len(scs.Items) == 0 {
				// TODO notice user to create storage class
				return errors.Errorf("no storage class found")
			}
			var scItems []string
			for _, i := range scs.Items {
				scItems = append(scItems, i.Name)
				if i.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
					if err := kubeClient.PatchDefaultSC(ctx, i.DeepCopy(), false); err != nil {
						return err
					}
				}
			}
			if len(scItems) == 1 {
				return kubeClient.PatchDefaultSC(ctx, scs.Items[0].DeepCopy(), true)
			}
			newDefaultSCName, err := logpkg.Question(&survey.QuestionOptions{
				Question:     "select default storage class",
				DefaultValue: scItems[0],
				Options:      scItems,
			})
			if err != nil {
				return errors.Errorf("select default storage class failed, reason: %v", err)
			}
			for _, sc := range scs.Items {
				if sc.Name == newDefaultSCName {
					return kubeClient.PatchDefaultSC(ctx, sc.DeepCopy(), true)
				}
			}
			return nil
		},
	}
	return ds
}
