package database

import (
	"os"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/models"
)

const dbPath string = ".seamoon.db"

var (
	initFlag bool
	globalDB *gorm.DB
	once     sync.Once
)

func DB() *gorm.DB {
	return globalDB
}

func init() {
	initFlag = false

	once.Do(func() {
		_, err := os.Stat(dbPath)
		if os.IsNotExist(err) {
			initFlag = true
		}

		globalDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}

		if initFlag {
			autoMigrate(globalDB)
		}
	})
}

func autoMigrate(db *gorm.DB) {
	for _, m := range models.ModelList {
		db.AutoMigrate(m)
	}
}
