package mysql

import (
	"github.com/imkuqin-zw/courier/pkg/gorm/driver"
	"gorm.io/driver/mysql"
)

func init() {
	driver.InstallerCtor("mysql", mysql.Open)
}
