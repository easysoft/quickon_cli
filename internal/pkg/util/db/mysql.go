// Copyright (c) 2021-2025 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"

	"github.com/easysoft/qcadmin/internal/pkg/util/log"

	_ "github.com/go-sql-driver/mysql"
)

func CheckMySQL(cfg *Config) bool {
	log := log.GetInstance()
	// Test MySQL connection
	log.Debugf("test mysql connection %s:%d", cfg.Host, cfg.Port)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		log.Errorf("failed to connect to mysql(%s:%d), err: %v", cfg.Host, cfg.Port, err)
		return false
	}
	defer db.Close()
	// Test database creation
	log.Debugf("try create database z_test_db")
	_, err = db.Exec("CREATE DATABASE z_test_db")
	if err != nil {
		log.Errorf("failed to create database z_test_db,err: %v", err)
		return false
	}
	log.Debugf("clean test database z_test_db")
	_, err = db.Exec("DROP DATABASE z_test_db")
	if err != nil {
		log.Errorf("failed to drop test database: %v", err)
		return false
	}
	log.Done("mysql check success")
	return true
}
