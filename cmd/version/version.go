// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package version

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/ergoapi/util/file"
	"github.com/ergoapi/util/github"

	gv "github.com/Masterminds/semver/v3"
	"github.com/cockroachdb/errors"
	"github.com/easysoft/qcadmin/common"
	logpkg "github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/easysoft/qcadmin/pkg/qucheng/upgrade"
	"github.com/ergoapi/util/color"
	"github.com/imroc/req/v3"
)

var versionTpl = `{{with .Client -}}
Client:
 Version:           {{ .Version }}
 Go version:        {{ .GoVersion }}
 Git commit:        {{ .GitCommit }}
 Built:             {{ .BuildTime }}
 OS/Arch:           {{.Os}}/{{.Arch}}
 Experimental:      {{.Experimental}}
{{- if .CanUpgrade }}
 Note:              {{ .UpgradeMessage }}
 URL:               https://github.com/easysoft/quickon_cli/releases
{{- end }}
{{- end}}
{{- if .ServerDeployed }}

Server:
 Type:              {{ .Server.ServerType -}}
{{with .Server.Components }}
 {{- range $component := .Components}}
 {{$component.Name}}:
{{- if $component.CanUpgrade }}
{{- if eq $component.Deploy.AppVersion $component.Remote.AppVersion }}
  AppVersion:       {{$component.Deploy.AppVersion}}
{{- else }}
  AppVersion:       {{$component.Deploy.AppVersion}} --> {{$component.Remote.AppVersion}}
{{- end }}
{{- if eq $component.Deploy.ChartVersion $component.Remote.ChartVersion }}
  ChartVersion:     {{$component.Deploy.ChartVersion}}
{{- else }}
  ChartVersion:     {{$component.Deploy.ChartVersion}} --> {{$component.Remote.ChartVersion}}
{{- end }}
  Note:             {{ $component.UpgradeMessage }}
{{- else }}
  AppVersion:       {{$component.Deploy.AppVersion}}
  ChartVersion:     {{$component.Deploy.ChartVersion}}
{{- end }}
 {{- end}}
{{- end}}
{{- end}}
`

const (
	defaultVersion       = "0.0.0"
	defaultGitCommitHash = "a1b2c3d4"
	defaultBuildDate     = "Mon Aug  3 15:06:50 2020"
)

type versionGet struct {
	Code int `json:"code"`
	Data struct {
		Name    string    `json:"name"`
		Version string    `json:"version"`
		Sync    time.Time `json:"sync"`
	} `json:"data"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

type versionInfo struct {
	Client clientVersion
	Server serverVersion
}

type clientVersion struct {
	Version        string
	LastVersion    string
	GoVersion      string
	GitCommit      string
	Os             string
	Arch           string
	BuildTime      string `json:",omitempty"`
	Experimental   bool
	CanUpgrade     bool
	UpgradeMessage string
}

type serverVersion struct {
	ServerType common.QuickonType
	Components *upgrade.Version
}

// ServerDeployed returns true when the client could connect to the qucheng
func (v versionInfo) ServerDeployed() bool {
	return v.Server.Components != nil
}

// PreCheckLatestVersion 检查最新版本
func PreCheckLatestVersion(log logpkg.Logger) (version, t string, err error) {
	version, _ = checkLastVersionFromGithub()
	if version != "" {
		log.Debugf("fetch version from github: %s", version)
		return version, "github", nil
	}
	version, err = checkLatestVersionFromAPI()
	if err != nil {
		return version, "api", err
	}
	log.Debugf("fetch version from api: %s", version)
	return version, "api", nil
}

func checkLastVersionFromGithub() (string, error) {
	pkg := github.Pkg{
		Owner: "easysoft",
		Repo:  "quickon_cli",
	}
	tag, err := pkg.LastTag()
	if err != nil {
		return "", err
	}
	return tag.Name, nil
}

func checkLatestVersionFromAPI() (string, error) {
	lastVersion := &versionGet{}
	client := req.C().SetLogger(nil).SetUserAgent(common.GetUG()).SetTimeout(time.Second * 5)
	_, err := client.R().SetSuccessResult(lastVersion).Get(common.GetAPI("/api/release/last/qcadmin"))
	if err != nil {
		return "", err
	}
	return lastVersion.Data.Version, nil
}

func ShowVersion(log logpkg.Logger) {
	// logo.PrintLogo()
	if common.Version == "" {
		common.Version = defaultVersion
	}
	if common.BuildDate == "" {
		common.BuildDate = defaultBuildDate
	}
	if common.GitCommitHash == "" {
		common.GitCommitHash = defaultGitCommitHash
	}
	tmpl, err := newVersionTemplate()
	if err != nil {
		log.Fatalf("gen version failed, reason: %v", err)
		return
	}
	vd := versionInfo{
		Client: clientVersion{
			Version:      common.Version,
			GoVersion:    runtime.Version(),
			GitCommit:    common.GitCommitHash,
			BuildTime:    common.BuildDate,
			Os:           runtime.GOOS,
			Arch:         runtime.GOARCH,
			Experimental: true,
		},
	}

	log.StartWait("check update...")
	lastVersion, lastType, err := PreCheckLatestVersion(log)
	log.StopWait()
	if err != nil {
		log.Debugf("get update message err: %v", err)
		return
	}
	if lastVersion != "" && !strings.Contains(common.Version, lastVersion) {
		nowVersion := gv.MustParse(strings.TrimPrefix(common.Version, "v"))
		needUpgrade := nowVersion.LessThan(gv.MustParse(lastVersion))
		if needUpgrade {
			vd.Client.CanUpgrade = true
			vd.Client.LastVersion = lastVersion
			vd.Client.Version = color.SGreen(vd.Client.Version)
			vd.Client.UpgradeMessage = fmt.Sprintf("Now you can use %s to upgrade cli to the latest version %s by %s mode", color.SGreen("%s upgrade cli", os.Args[0]), color.SGreen(lastVersion), color.SGreen(lastType))
		}
	}

	if file.CheckFileExists(common.GetCustomConfig(common.InitFileName)) {
		cfg, _ := config.LoadConfig()
		if cfg != nil {
			vd.Server.ServerType = cfg.Quickon.Type
			if cfg.Quickon.DevOps {
				vd.Server.ServerType = common.QuickonType(fmt.Sprintf("devops.%s", vd.Server.ServerType))
			}

			qv, err := upgrade.QuchengVersion(cfg.Quickon.DevOps)
			if err == nil {
				vd.Server.Components = &qv
			}
		}
	}
	if err := prettyPrintVersion(vd, tmpl); err != nil {
		panic(err)
	}
}

func prettyPrintVersion(vd versionInfo, tmpl *template.Template) error {
	t := tabwriter.NewWriter(os.Stdout, 20, 1, 1, ' ', 0)
	err := tmpl.Execute(t, vd)
	t.Write([]byte("\n"))
	t.Flush()
	return err
}

func newVersionTemplate() (*template.Template, error) {
	tmpl, err := template.New("version").Parse(versionTpl)
	return tmpl, errors.Wrap(err, "template parsing error")
}
