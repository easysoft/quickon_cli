// Copyright (c) 2021-2023 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package status

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/easysoft/qcadmin/common"
	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/output"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
)

type MapCount map[string]int

type MapMapCount map[string]MapCount

type PodStateMap map[string]PodStateCount

type Status struct {
	output     string     `json:"-" yaml:"-"`
	KubeStatus KubeStatus `json:"k8s" yaml:"k8s"`
	QStatus    QStatus    `json:"platform" yaml:"platform"`
}

type KubeStatus struct {
	Version   string      `json:"version" yaml:"version"`
	Type      string      `json:"type" yaml:"type"`
	NodeCount MapCount    `json:"nodes" yaml:"nodes"`
	PodState  PodStateMap `json:"service,omitempty" yaml:"service,omitempty"`
}

type PodStateCount struct {
	// Type is the type of deployment ("Deployment", "DaemonSet", ...)
	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// Desired is the number of desired pods to be scheduled
	Desired int `json:"desired,omitempty" yaml:"desired,omitempty"`

	// Ready is the number of ready pods
	Ready int `json:"ready,omitempty" yaml:"ready,omitempty"`

	// Available is the number of available pods
	Available int `json:"available,omitempty" yaml:"available,omitempty"`

	// Unavailable is the number of unavailable pods
	Unavailable int `json:"unavailable,omitempty" yaml:"unavailable,omitempty"`

	Disabled bool `json:"disabled" yaml:"disabled"`
}

type QStatus struct {
	PodState    PodStateMap `json:"service,omitempty" yaml:"service,omitempty"`
	PluginState PodStateMap `json:"plugin,omitempty" yaml:"plugin,omitempty"`
}

func newStatus(output string) *Status {
	return &Status{
		output: output,
		KubeStatus: KubeStatus{
			NodeCount: MapCount{},
			Version:   "unknow",
			PodState:  PodStateMap{},
		},
		QStatus: QStatus{
			PodState:    PodStateMap{},
			PluginState: PodStateMap{},
		},
	}
}

func (s *Status) Format() error {
	switch strings.ToLower(s.output) {
	case "json":
		return output.EncodeJSON(os.Stdout, s)
	case "yaml":
		return output.EncodeYAML(os.Stdout, s)
	default:
		var buf bytes.Buffer
		w := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', 0)
		// k8s
		fmt.Fprintf(w, "Cluster Status: \n")
		fmt.Fprintf(w, "  %s\t%s\n", "version", s.KubeStatus.Version)
		fmt.Fprintf(w, "  %s\t%s\n", "mode", s.KubeStatus.Type)
		if s.KubeStatus.NodeCount["ready"] > 0 {
			fmt.Fprintf(w, "  %s\t%s\n", "status", color.SGreen("health"))
		} else {
			fmt.Fprintf(w, "  %s\t%s\n", "status", color.SRed("unhealth"))
		}
		for name, state := range s.KubeStatus.PodState {
			if state.Disabled {
				fmt.Fprintf(w, "  %s\t%s\n", name, color.SBlue("disabled"))
			} else {
				if state.Available > 0 {
					fmt.Fprintf(w, "  %s\t%s\n", name, color.SGreen("ok"))
				} else {
					fmt.Fprintf(w, "  %s\t%s\n", name, color.SRed("warn"))
				}
			}
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "Platform Status: \n")
		if s.QStatus.PodState["platform"].Disabled {
			fmt.Fprintf(w, "  %s\t%s\n", "status", color.SBlue("disabled"))
			w.Flush()
			return output.EncodeText(os.Stdout, buf.Bytes())
		}
		cfg, _ := config.LoadConfig()
		domain := ""
		loginIP := exnet.LocalIPs()[0]
		if cfg != nil {
			domain = cfg.Domain
		}
		consoleURL := ""
		if len(domain) > 0 {
			// TODO 优化域名显示结果
			if !kutil.IsLegalDomain(domain) {
				if cfg.Quickon.DevOps {
					consoleURL = fmt.Sprintf("http://zentao.%s", domain)
				} else {
					consoleURL = fmt.Sprintf("http://console.%s", domain)
				}
			} else {
				if cfg.Quickon.Domain.Type == "custom" {
					consoleURL = fmt.Sprintf("http://%s", domain)
				} else {
					consoleURL = fmt.Sprintf("https://%s", domain)
				}
			}
		} else {
			consoleURL = fmt.Sprintf("http://%s:32379", loginIP)
		}
		fmt.Fprintf(w, "  namespace:\t%s\n", color.SBlue(common.GetDefaultSystemNamespace(true)))
		if cfg.Quickon.DevOps {
			fmt.Fprintf(w, "  console:      %s\n", color.SGreen(consoleURL))
		} else {
			fmt.Fprintf(w, "  console:      %s(%s/%s)\n", color.SGreen(consoleURL), color.SGreen(common.QuchengDefaultUser), color.SGreen(cfg.ConsolePassword))
		}
		ptOK := true
		fmt.Fprintf(w, "  component status: \n")
		for name, state := range s.QStatus.PodState {
			if state.Disabled {
				fmt.Fprintf(w, "    %s\t%s\n", name, color.SBlue("disabled"))
			} else {
				if state.Available > 0 {
					fmt.Fprintf(w, "    %s\t%s\n", name, color.SGreen("ok"))
				} else {
					fmt.Fprintf(w, "    %s\t%s\n", name, color.SRed("warn"))
					ptOK = false
				}
			}
		}
		fmt.Fprintf(w, "  plugin status: \n")
		for name, state := range s.QStatus.PluginState {
			if state.Disabled {
				fmt.Fprintf(w, "    %s\t%s\n", name, color.SBlue("disabled"))
			} else {
				fmt.Fprintf(w, "    %s\t%s\n", name, color.SGreen("enabled"))
			}
		}
		if ptOK {
			fmt.Fprintf(w, "  %s\t%s\n", "status", color.SGreen("health"))
		} else {
			fmt.Fprintf(w, "  %s\t%s\n", "status", color.SRed("unhealth"))
		}
		w.Flush()
		return output.EncodeText(os.Stdout, buf.Bytes())
	}
}
