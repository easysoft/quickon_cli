package backup

import (
	"fmt"

	"github.com/easysoft/qcadmin/internal/app/config"
	"github.com/easysoft/qcadmin/internal/pkg/util/exec"
	"github.com/easysoft/qcadmin/internal/pkg/util/factory"
	"github.com/spf13/cobra"
)

func NewCmdBackupCluster(f factory.Factory) *cobra.Command {
	log := f.GetLog()
	bc := &cobra.Command{
		Use:     "backup",
		Long:    "backup cluster",
		Aliases: []string{"snapshot"},
		RunE: func(cmd *cobra.Command, args []string) error {
			log.StartWait("backup cluster")
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("found config err, reason: %v", err)
			}
			if cfg.DB != "" && cfg.DB != "etcd" {
				return fmt.Errorf("not support datastore %s", cfg.DB)
			}
			exec.CommandRun("bash", "-c", "")
			return nil
		},
	}
	return bc
}
