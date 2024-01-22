package database

import (
	"log/slog"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
)

const dbPath string = ".seamoon.db"

var (
	globalDB    *gorm.DB
	migrateFunc []func()
)

func Init() {
	gormConfig := gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var err error

	if _, exist := os.Stat(dbPath); os.IsNotExist(exist) {
		defer func() {
			slog.Info("初始化数据库......")
			for _, m := range models.ModelList {
				// 初始化表
				globalDB.AutoMigrate(m)
			}
			// 写表
			for _, fn := range migrateFunc {
				fn()
			}
		}()
	}

	globalDB, err = gorm.Open(sqlite.Open(dbPath), &gormConfig)
	if err != nil {
		panic(err)
	}
}

func GetConn() *gorm.DB {
	return globalDB
}

// QueryPage 公共查询所有数据的方法
func QueryPage(page, size int) *gorm.DB {
	if page < 0 {
		page = 0
	}
	if size <= 0 {
		size = 10
	}
	return globalDB.Offset(page * size).Limit(size)
}

func RegisterMigrate(f func()) {
	migrateFunc = append(migrateFunc, f)
}
