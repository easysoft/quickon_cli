// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package common

import "testing"

func TestGetVersion(t *testing.T) {
	type args struct {
		devops  bool
		typ     string
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				devops:  true,
				typ:     "oss",
				version: "v1.0.0",
			},
			want: "v1.0.0",
		},
		{
			name: "test2",
			args: args{
				devops:  true,
				typ:     "oss",
				version: "",
			},
			want: DefaultZentaoDevOPSOSSVersion,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersion(tt.args.devops, tt.args.typ, tt.args.version); got != tt.want {
				t.Logf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
