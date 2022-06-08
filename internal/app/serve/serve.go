// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.cn) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package serve

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/easysoft/qcadmin/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func Serve(ctx context.Context) error {
	logrus.Info("Starting q daemon ...")
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	r := gin.Default()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", common.DefaultDaemonPort),
		Handler: r,
	}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Errorf("Failed to stop server, error: %s", err)
		}
		logrus.Info("server exited.")
	}()
	logrus.Infof("start application server, Listen on port: %d, pid is %v", common.DefaultDaemonPort, os.Getpid())
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Errorf("Failed to start http server, error: %s", err)
		return err
	}
	return nil
}
