// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CreateEndpoint(ctx context.Context, namespace string, endpoint *corev1.Endpoints, opts metav1.CreateOptions) (*corev1.Endpoints, error) {
	return c.Clientset.CoreV1().Endpoints(namespace).Create(ctx, endpoint, opts)
}

func (c *Client) UpdateEndpoint(ctx context.Context, namespace string, endpoint *corev1.Endpoints, opts metav1.UpdateOptions) (*corev1.Endpoints, error) {
	return c.Clientset.CoreV1().Endpoints(namespace).Update(ctx, endpoint, opts)
}

func (c *Client) DeleteEndpoint(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Endpoints(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetEndpoint(ctx context.Context, namespace string, name string, opts metav1.GetOptions) (*corev1.Endpoints, error) {
	return c.Clientset.CoreV1().Endpoints(namespace).Get(ctx, name, opts)
}
