package hostinfo

import (
	"fmt"
	"github.com/easysoft/qcadmin/common"
	hinfo "tailscale.com/hostinfo"
	"tailscale.com/tailcfg"
)

// New returns a partially populated Hostinfo for the current host.
func New() *tailcfg.Hostinfo {
	t := hinfo.New()
	t.IPNVersion = fmt.Sprintf("%s-%s", common.Version, common.GitCommitHash)
	return t
}
