package gorm

import (
	"context"
	"strings"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	gormLg "gorm.io/gorm/logger"
)

type zapLogger struct {
	slowThreshold time.Duration
}

func (zl *zapLogger) LogMode(level gormLg.LogLevel) gormLg.Interface {
	return zl
}

func (zl *zapLogger) Info(ctx context.Context, s string, i ...interface{}) {
	logger.Infof(strings.TrimRight(s, "\n"), i...)
}

func (zl *zapLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	logger.Warnf(s, i...)
}

func (zl *zapLogger) Error(ctx context.Context, s string, i ...interface{}) {
	logger.Errorf(s, i...)
}

func (zl *zapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	cost := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		logger.Errorf("sql[%s], cost[%s], rows[%d], error: %s", sql, cost, rows, err.Error())
		return
	}

	if zl.slowThreshold < cost {
		logger.Warnf("sql[%s], cost[%d], rows[%d], slow", sql, cost, rows)
		return
	}

	logger.Infof("sql[%s], cost[%d], rows[%d]", sql, cost, rows)
	return

}
