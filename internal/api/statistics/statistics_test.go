// Copyright (c) 2021-2024 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package statistics

import "testing"

func TestSendStatistics(t *testing.T) {
	type args struct {
		action string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestSendStatistics",
			args:    args{action: "install"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendStatistics(tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("SendStatistics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
