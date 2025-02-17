// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/ergoapi/util/exhash"
	"github.com/ergoapi/util/exmap"

	"github.com/easysoft/qcadmin/internal/pkg/k8s"

	quchengv1beta1 "github.com/easysoft/quickon-api/qucheng/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func fakeDbUserInfo(qclient *k8s.Client, db *quchengv1beta1.Db) error {
	if db.Spec.Account.User.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), db.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: db.Spec.Account.User.ValueFrom.SecretKeyRef.Name,
			},
			Key: db.Spec.Account.User.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		db.Spec.Account.User.Value = string(user)
	}
	if db.Spec.Account.Password.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), db.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: db.Spec.Account.Password.ValueFrom.SecretKeyRef.Name,
			},
			Key: db.Spec.Account.Password.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		db.Spec.Account.Password.Value = string(user)
	}
	return nil
}

func vaildGlobalDatabase(l map[string]string) bool {
	if exmap.CheckLabel(l, "easycorp.io/global_database") {
		return exmap.GetLabelValue(l, "easycorp.io/global_database") == "true"
	}
	return false
}

func fakeDbSvcUserInfo(qclient *k8s.Client, dbsvc *quchengv1beta1.DbService) error {
	if dbsvc.Spec.Account.User.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), dbsvc.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: dbsvc.Spec.Account.User.ValueFrom.SecretKeyRef.Name,
			},
			Key: dbsvc.Spec.Account.User.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		dbsvc.Spec.Account.User.Value = string(user)
	}
	if dbsvc.Spec.Account.Password.Value == "" {
		user, err := qclient.GetSecretKeyBySelector(context.TODO(), dbsvc.Namespace, &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: dbsvc.Spec.Account.Password.ValueFrom.SecretKeyRef.Name,
			},
			Key: dbsvc.Spec.Account.Password.ValueFrom.SecretKeyRef.Key,
		})
		if err != nil {
			return err
		}
		dbsvc.Spec.Account.Password.Value = string(user)
	}
	return nil
}

func dbsvcResourceAlias(name string) string {
	return exhash.B64EnCode(name)
}
