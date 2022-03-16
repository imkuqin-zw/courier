package driver

import "gorm.io/gorm"

type Constructor func(dsn string) gorm.Dialector

var driver map[string]Constructor

func InstallerCtor(name string, ctor Constructor) {
	driver[name] = ctor
}

func GetDriverCtor(name string) (Constructor, bool) {
	ctor, ok := driver[name]
	if ok {
		return ctor, true
	}
	return nil, false
}
