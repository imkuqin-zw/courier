package gorm

import (
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/imkuqin-zw/courier/pkg/gorm/plugin/trace"
	"gorm.io/gorm"
)

const (
	defaultMaxIdleConn     = 10
	defaultMaxOpenConn     = 100
	defaultConnMaxLifetime = time.Second * 300
	defaultSlowThreshold   = time.Millisecond * 500
)

// conf options
type Config struct {
	Driver string
	// DSN: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	DSN string

	PrepareStmt            bool
	DryRun                 bool
	SkipDefaultTransaction bool
	MaxIdleConn            int
	MaxOpenConn            int
	ConnMaxLifetime        time.Duration

	// 慢日志阈值
	SlowThreshold time.Duration
	EnableTrace   bool

	// 记录错误sql时,是否打印包含参数的完整sql语句
	DetailSQL bool

	plugins []gorm.Plugin
	dbCfg   *mysql.Config
}

//WithInterceptors
func (config *Config) WithInterceptors(plugins ...gorm.Plugin) *Config {
	if config.plugins == nil {
		config.plugins = make([]gorm.Plugin, 0, len(plugins))
	}
	config.plugins = append(config.plugins, plugins...)
	return config
}

//Check
func (config *Config) Check() (err error) {
	config.dbCfg, err = mysql.ParseDSN(config.DSN)
	if err != nil {
		return
	}

	return nil
}

//Build
func (config *Config) Build() *gorm.DB {
	if err := config.Check(); err != nil {
		panic(err)
	}
	if config.MaxIdleConn == 0 {
		config.MaxIdleConn = defaultMaxIdleConn
	}
	if config.MaxOpenConn == 0 {
		config.MaxOpenConn = defaultMaxOpenConn
	}
	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = defaultConnMaxLifetime
	}
	if config.SlowThreshold == 0 {
		config.SlowThreshold = defaultSlowThreshold
	}
	if config.EnableTrace {
		config.WithInterceptors(trace.New(
			config.dbCfg.User,
			config.dbCfg.DBName,
			config.dbCfg.Addr,
			config.DetailSQL,
		))
	}
	return open(config)
}
