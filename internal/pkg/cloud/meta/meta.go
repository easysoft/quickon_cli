// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package meta

import (
	"strings"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/idoubi/goz"
)

// Meta is the meta data of the system
func CheckCloudMetadataAPI() {
	for _, apiInstance := range common.CloudAPI {
		cli := goz.NewClient(goz.Options{
			Timeout: 1,
		})
		resp, err := cli.Get(apiInstance.API)
		if err != nil {
			continue
		}
		r, _ := resp.GetBody()
		if strings.Contains(r.String(), apiInstance.ResponseMatch) {
			log.Flog.Debug("%s Metadata API available in %s", apiInstance.CloudProvider, apiInstance.API)
		} else {
			log.Flog.Debug("failed to dial %s API.", apiInstance.CloudProvider)
		}
	}
}
