// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package k8s

import (
	"context"

	quchengv1beta1 "github.com/easysoft/quickon-api/qucheng/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) ListDBSvc(ctx context.Context, namespace string, opts metav1.ListOptions) (*quchengv1beta1.DbServiceList, error) {
	return c.QClient.QuchengV1beta1().DbServices(namespace).List(ctx, opts)
}

func (c *Client) GetDBSvc(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*quchengv1beta1.DbService, error) {
	return c.QClient.QuchengV1beta1().DbServices(namespace).Get(ctx, name, opts)
}

func (c *Client) CreateDBSvc(ctx context.Context, namespace string, dbSvc *quchengv1beta1.DbService, opts metav1.CreateOptions) (*quchengv1beta1.DbService, error) {
	return c.QClient.QuchengV1beta1().DbServices(namespace).Create(ctx, dbSvc, opts)
}

func (c *Client) UpdateDBSvc(ctx context.Context, namespace, name string, dbSvc *quchengv1beta1.DbService, opts metav1.UpdateOptions) (*quchengv1beta1.DbService, error) {
	return c.QClient.QuchengV1beta1().DbServices(namespace).Update(ctx, dbSvc, opts)
}

func (c *Client) DeleteDBSvc(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.QClient.QuchengV1beta1().DbServices(namespace).Delete(ctx, name, opts)
}

func (c *Client) ListDB(ctx context.Context, namespace string, opts metav1.ListOptions) (*quchengv1beta1.DbList, error) {
	return c.QClient.QuchengV1beta1().Dbs(namespace).List(ctx, opts)
}

func (c *Client) GetDB(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*quchengv1beta1.Db, error) {
	return c.QClient.QuchengV1beta1().Dbs(namespace).Get(ctx, name, opts)
}

func (c *Client) CreateDB(ctx context.Context, namespace string, db *quchengv1beta1.Db, opts metav1.CreateOptions) (*quchengv1beta1.Db, error) {
	return c.QClient.QuchengV1beta1().Dbs(namespace).Create(ctx, db, opts)
}

func (c *Client) UpdateDB(ctx context.Context, namespace, name string, db *quchengv1beta1.Db, opts metav1.UpdateOptions) (*quchengv1beta1.Db, error) {
	return c.QClient.QuchengV1beta1().Dbs(namespace).Update(ctx, db, opts)
}

func (c *Client) DeleteDB(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.QClient.QuchengV1beta1().Dbs(namespace).Delete(ctx, name, opts)
}
