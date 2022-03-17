// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorm

import (
	"context"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"github.com/imkuqin-zw/courier/pkg/config"
	"github.com/imkuqin-zw/courier/pkg/gorm/driver"
	"gorm.io/gorm"
)

func open(config *Config) *gorm.DB {
	cfg := &gorm.Config{
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		Logger: &zapLogger{
			slowThreshold: config.SlowThreshold,
		},
		DryRun:      config.DryRun,
		PrepareStmt: config.PrepareStmt,
	}

	ctor, ok := driver.GetDriverCtor(config.Driver)
	if !ok {
		logger.Fatalf("unknown driver[%s]", config.Driver)
		return nil
	}
	db, err := gorm.Open(ctor(config.DSN), cfg)
	if err != nil {
		logger.Fatalf("fault to open mysql, err: %s", err.Error())
		return nil
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil
	}
	sqlDb.SetMaxOpenConns(config.MaxOpenConn)
	sqlDb.SetMaxIdleConns(config.MaxIdleConn)
	sqlDb.SetMaxIdleConns(config.MaxIdleConn)
	sqlDb.SetConnMaxLifetime(config.ConnMaxLifetime)

	for _, plugin := range config.plugins {
		if err := db.Use(plugin); err != nil {
			logger.Fatalf("fault to use plugin, err: %s", err.Error())
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	defer cancel()
	if err := sqlDb.PingContext(ctx); err != nil {
		logger.Fatalf("fault to ping mysql server, err: %s", err.Error())
		return nil
	}
	return db
}

func New() *gorm.DB {
	cfg := &Config{}
	if err := config.Get("gorm").Scan(&cfg); err != nil {
		logger.Fatalf("fault to scan mongo config, error: %s", err.Error())
	}
	return cfg.Build()
}
