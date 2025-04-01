// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package crd

import (
	"context"
	"net"

	"github.com/ergoapi/util/color"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/k8s"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"

	quchengv1beta1 "github.com/easysoft/quickon-api/qucheng/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// cmdExternalDbSvc 外部数据库主命令
func cmdExternalDbSvc(f factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "external",
		Aliases: []string{"etdb"},
		Short:   "manage external database for platform",
	}

	// 添加子命令
	cmd.AddCommand(cmdCreateExternalDbSvc(f))
	cmd.AddCommand(cmdDeleteExternalDbSvc(f))
	return cmd
}

// cmdCreateExternalDbSvc 创建外部数据库子命令
func cmdCreateExternalDbSvc(f factory.Factory) *cobra.Command {
	var host string
	var port int32
	var namespace string
	var name string
	var superUser string
	var superPassword string

	log := f.GetLog()
	cmd := &cobra.Command{
		Use:     "new",
		Short:   "create new external database for platform",
		Version: "4.0.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			if host == "" || port == 0 {
				log.Fatalf("host and port are required")
				return nil
			}

			// 创建 k8s client
			client, err := k8s.NewSimpleQClient()
			if err != nil {
				log.Fatalf("failed to connect to k8s cluster: %v", err)
				return nil
			}

			ctx := context.Background()
			svc := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				},
			}

			// 判断是IP还是域名
			if net.ParseIP(host) != nil {
				// IP方式: 使用 Endpoints + ClusterIP Service
				svc.Spec = corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Port:       port,
							Protocol:   corev1.ProtocolTCP,
							Name:       "db",
							TargetPort: intstr.FromInt(int(port)),
						},
					},
				}

				// 创建 Endpoints
				eps := &corev1.Endpoints{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: namespace,
					},
					Subsets: []corev1.EndpointSubset{
						{
							Addresses: []corev1.EndpointAddress{
								{
									IP: host,
								},
							},
							Ports: []corev1.EndpointPort{
								{
									Port:     port,
									Protocol: corev1.ProtocolTCP,
									Name:     "db",
								},
							},
						},
					},
				}

				// 创建或更新 Endpoints
				_, err = client.GetEndpoint(ctx, namespace, name, metav1.GetOptions{})
				if err != nil {
					_, err = client.CreateEndpoint(ctx, namespace, eps, metav1.CreateOptions{})
				} else {
					_, err = client.UpdateEndpoint(ctx, namespace, eps, metav1.UpdateOptions{})
				}
				if err != nil {
					log.Fatalf("failed to create/update endpoints: %v", err)
					return nil
				}
			} else {
				// 域名方式: 使用 ExternalName Service
				svc.Spec = corev1.ServiceSpec{
					Type:         corev1.ServiceTypeExternalName,
					ExternalName: host,
					Ports: []corev1.ServicePort{
						{
							Port:       port,
							TargetPort: intstr.FromInt32(port),
							Protocol:   corev1.ProtocolTCP,
							Name:       "db",
						},
					},
				}
			}

			_, err = client.GetService(ctx, namespace, name, metav1.GetOptions{})
			if err != nil {
				_, err = client.CreateService(ctx, namespace, svc, metav1.CreateOptions{})
			} else {
				_, err = client.UpdateService(ctx, namespace, svc, metav1.UpdateOptions{})
			}
			if err != nil {
				log.Fatalf("failed to create/update service: %v", err)
				return nil
			}
			if len(superPassword) > 0 && len(superUser) > 0 {
				log.Debug("detch super user & password for external database, will create dbsvc")
				dbsvc := &quchengv1beta1.DbService{
					ObjectMeta: metav1.ObjectMeta{
						Name:      svc.Name,
						Namespace: svc.Namespace,
						Annotations: map[string]string{
							"easycorp.io/resource_alias": dbsvcResourceAlias(svc.Name),
						},
						Labels: map[string]string{
							"easycorp.io/global_database": "true",
						},
					},
					Spec: quchengv1beta1.DbServiceSpec{
						Type: quchengv1beta1.DbTypeMysql,
						Service: quchengv1beta1.Service{
							Name:      svc.Name,
							Namespace: svc.Namespace,
							Port:      intstr.FromInt32(port),
						},
						Account: quchengv1beta1.Account{
							User: quchengv1beta1.AccountUser{
								Value: superUser,
							},
							Password: quchengv1beta1.AccountPassword{
								Value: superPassword,
							},
						},
					},
				}
				_, err = client.GetDBSvc(ctx, namespace, svc.Name, metav1.GetOptions{})
				if err != nil {
					_, err = client.CreateDBSvc(ctx, namespace, dbsvc, metav1.CreateOptions{})
				} else {
					_, err = client.UpdateDBSvc(ctx, namespace, svc.Name, dbsvc, metav1.UpdateOptions{})
				}
				if err != nil {
					log.Fatalf("failed to create/update dbsvc: %v", err)
					return nil
				}
			} else {
				log.Warn("ignore create dbsvc for external database")
				return nil
			}
			log.Donef("created external database service %s in namespace %s", name, namespace)
			log.Infof("you can access the database in cluster using: %s", color.SGreen("%s.%s:%d", name, namespace, port))
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "External database host (required)")
	cmd.Flags().Int32Var(&port, "port", common.DefaultExternalDBPort, "External database port")
	cmd.Flags().StringVar(&namespace, "namespace", common.GetDefaultSystemNamespace(true), "Kubernetes namespace")
	cmd.Flags().StringVar(&name, "name", common.DefaultExternalDBName, "Service name")
	cmd.Flags().StringVar(&superUser, "username", "root", "Super username for the database")
	cmd.Flags().StringVar(&superPassword, "password", "", "Super user password for the database")
	return cmd
}

// cmdDeleteExternalDbSvc 删除外部数据库子命令
func cmdDeleteExternalDbSvc(f factory.Factory) *cobra.Command {
	var namespace string
	var name string

	log := f.GetLog()
	cmd := &cobra.Command{
		Use:     "clean",
		Aliases: []string{"delete"},
		Version: "4.0.0",
		Short:   "delete external database from platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 创建 k8s client
			client, err := k8s.NewSimpleQClient()
			if err != nil {
				log.Fatalf("failed to connect to k8s cluster: %v", err)
				return nil
			}

			ctx := context.Background()

			// 检查服务是否存在
			_, err = client.GetService(ctx, namespace, name, metav1.GetOptions{})
			if err != nil {
				log.Fatalf("service %s not found in namespace %s: %v", name, namespace, err)
				return nil
			}

			// 删除 DbService (如果存在)
			_, err = client.GetDBSvc(ctx, namespace, name, metav1.GetOptions{})
			if err == nil {
				err = client.DeleteDBSvc(ctx, namespace, name, metav1.DeleteOptions{})
				if err != nil {
					log.Warnf("failed to delete dbsvc %s: %v", name, err)
				} else {
					log.Infof("deleted dbsvc %s in namespace %s", name, namespace)
				}
			}

			// 删除 Service
			err = client.DeleteService(ctx, namespace, name, metav1.DeleteOptions{})
			if err != nil {
				log.Fatalf("failed to delete service %s: %v", name, err)
				return nil
			}

			// 尝试删除 Endpoints (如果存在)
			_, err = client.GetEndpoint(ctx, namespace, name, metav1.GetOptions{})
			if err == nil {
				err = client.DeleteEndpoint(ctx, namespace, name, metav1.DeleteOptions{})
				if err != nil {
					log.Warnf("failed to delete endpoints %s: %v", name, err)
				}
			}

			log.Donef("deleted external database service %s in namespace %s", name, namespace)
			return nil
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", common.GetDefaultSystemNamespace(true), "Kubernetes namespace")
	cmd.Flags().StringVar(&name, "name", common.DefaultExternalDBName, "Service name")
	return cmd
}
