package drivers

import (
	"os"

	"github.com/djherbis/times"
	"github.com/glebarez/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/keys"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

const dbPath string = ".seamoon.db"

type sqlite3 struct {
	db *gorm.DB
}

func (s *sqlite3) Init(migrateFunc []func()) {
	gormConfig := gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var err error

	if _, exist := os.Stat(dbPath); os.IsNotExist(exist) {
		defer func() {
			xlog.Info(xlog.DatabaseInit)
			for _, m := range models.ModelList {
				// 初始化表
				s.db.AutoMigrate(m)
			}
			// 写表
			for _, fn := range migrateFunc {
				fn()
			}
		}()
	}

	defer func() {
		t, err := times.Stat(dbPath)
		if err != nil || !t.HasBirthTime() {
			panic(err)
		}
		keys.SetGlobalKey(tools.GenerateKey(t.BirthTime()))
	}()

	s.db, err = gorm.Open(sqlite.Open(dbPath), &gormConfig)
	dao.SetDefault(s.db)
	if err != nil {
		panic(err)
	}
}

func (s *sqlite3) GetConn() *gorm.DB {
	return s.db
}

// QueryPage 公共查询所有数据的方法
func (s *sqlite3) QueryPage(page, size int) *gorm.DB {
	if page < 0 {
		page = 0
	}
	if size <= 0 {
		size = 10
	}
	return s.db.Offset(page * size).Limit(size)
}

func (s *sqlite3) Generate() error {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "pkg/api/database/dao", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(s.db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(models.ModelList...)

	// Execute the generator
	g.Execute()
	return nil
}
