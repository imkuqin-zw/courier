package gorm

import (
	"context"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"gorm.io/driver/mysql"
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
	db, err := gorm.Open(mysql.Open(config.DSN), cfg)
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
