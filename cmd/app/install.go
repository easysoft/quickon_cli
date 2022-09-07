package app

import (
	"fmt"
	"strings"

	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/ergoapi/util/color"
	"github.com/ergoapi/util/exnet"
	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Account string `json:"account"`
	} `json:"data"`
}

type Body struct {
	Chart  string `json:"chart"`
	Domain string `json:"domain"`
}

func NewCmdAppInstall(f factory.Factory) *cobra.Command {
	var name, domain string
	log := f.GetLog()
	app := &cobra.Command{
		Use:     "install",
		Short:   "install app",
		Example: `q app install -n zentao-open or q app install zentao-open`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				name = args[0]
			}
			if len(domain) == 0 {
				domain = name
			}
			cfg, _ := config.LoadConfig()
			apiHost := cfg.Domain
			if apiHost == "" {
				apiHost = fmt.Sprintf("%s:32379", exnet.LocalIPs()[0])
			} else if !strings.HasSuffix(apiHost, "haogs.cn") && !strings.HasSuffix(apiHost, "corp.cc") {
				apiHost = fmt.Sprintf("console.%s", cfg.Domain)
			}
			log.Infof("install app %s, domain: %s.%s", name, domain, cfg.Domain)
			client := req.C()
			if log.GetLevel() > logrus.InfoLevel {
				client = client.DevMode().EnableDumpAll()
			}
			resp, err := client.R().
				SetHeader("accept", "application/json").
				SetHeader("TOKEN", cfg.APIToken).
				SetBody(&Body{Chart: name, Domain: domain}).
				Post(fmt.Sprintf("http://%s/instance-apiInstall.html", apiHost))
			if err != nil {
				log.Errorf("install app %s failed, reason: %v", name, err)
				return fmt.Errorf("install app %s failed, reason: %v", name, err)
			}
			if !resp.IsSuccess() {
				log.Errorf("install app %s failed, reason: bad response status %v", name, resp.Status)
				return fmt.Errorf("install app %s failed, reason: bad response status %v", name, resp.Status)
			}
			log.Donef("app %s install success.", name)
			log.Infof("please wait, the app is starting. \n\t app url: %s", color.SGreen(fmt.Sprintf("%s.%s", domain, cfg.Domain)))
			return nil
		},
	}
	app.Flags().StringVarP(&name, "name", "n", "zentao-open", "app name")
	app.Flags().StringVarP(&domain, "domain", "d", "", "app subdomain")
	return app
}
