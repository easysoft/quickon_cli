// Copyright (c) 2021-2024 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) CountPodNumByNamespace(ctx context.Context, ns string) (int, error) {
	if len(ns) == 0 {
		ns = metav1.NamespaceAll
	}
	podList, err := c.Clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(podList.Items), nil
}