package drivers

import (
	"gorm.io/gorm"
)

var migrateFunc = make([]func(), 0)
var driversMap = map[string]Driver{
	"sqlite3": &sqlite3{},
}

// Driver 用于后期支持其他 db 格式，此处目前仅支持 sqlite
type Driver interface {
	Init(migrateFunc []func())
	GetConn() *gorm.DB
	QueryPage(page, size int) *gorm.DB
	Generate() error
}

func Init() {
	driversMap["sqlite3"].Init(migrateFunc)
}

func Drive() Driver {
	return driversMap["sqlite3"]
}

func RegisterMigrate(f func()) {
	migrateFunc = append(migrateFunc, f)
}
