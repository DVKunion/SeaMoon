package drivers

import (
	"os"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

const dbPath string = ".seamoon.db"

type sqlite3 struct {
	db *gorm.DB
}

// migrationColumns 定义需要检查和迁移的新字段
// 格式: tableName -> []columnDefinition{name, type, default}
var migrationColumns = map[string][]struct {
	Name    string
	Type    string
	Default string
}{
	"tunnels": {
		{"version", "TEXT", ""},
		{"v2ray_version", "TEXT", ""},
		{"last_check_time", "TEXT", ""},
		{"cascade_proxy", "INTEGER", "0"},
		{"cascade_tunnel_id", "INTEGER", "0"},
		{"cascade_addr", "TEXT", ""},
		{"cascade_uid", "TEXT", ""},
		{"cascade_password", "TEXT", ""},
	},
}

func (s *sqlite3) Init(migrateFunc []func()) {
	gormConfig := gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var err error
	isNewDB := false

	if _, exist := os.Stat(dbPath); os.IsNotExist(exist) {
		isNewDB = true
	}

	s.db, err = gorm.Open(sqlite.Open(dbPath), &gormConfig)
	dao.SetDefault(s.db)
	if err != nil {
		panic(err)
	}

	// 每次启动都执行 AutoMigrate 以确保表结构是最新的
	for _, m := range models.ModelList {
		s.db.AutoMigrate(m)
	}

	// 执行手动迁移，确保新增字段存在（兼容旧版本数据库）
	s.manualMigrate()

	// 只有新数据库才执行初始化数据写入
	if isNewDB {
		xlog.Info(xlog.DatabaseInit)
		for _, fn := range migrateFunc {
			fn()
		}
	}
}

// manualMigrate 手动检查并添加可能缺失的列
// 这是为了兼容旧版本数据库，因为 SQLite 的 AutoMigrate 有时候无法正确添加新列
func (s *sqlite3) manualMigrate() {
	for tableName, columns := range migrationColumns {
		for _, col := range columns {
			if !s.columnExists(tableName, col.Name) {
				xlog.Info(xlog.DatabaseMigrate, "table", tableName, "column", col.Name)
				sql := "ALTER TABLE " + tableName + " ADD COLUMN " + col.Name + " " + col.Type
				if col.Default != "" {
					sql += " DEFAULT " + col.Default
				}
				if err := s.db.Exec(sql).Error; err != nil {
					xlog.Error(xlog.DatabaseMigrateError, "table", tableName, "column", col.Name, "err", err)
				}
			}
		}
	}
}

// columnExists 检查表中是否存在指定列
func (s *sqlite3) columnExists(tableName, columnName string) bool {
	var count int64
	s.db.Raw("SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?", tableName, columnName).Scan(&count)
	return count > 0
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

func (s *sqlite3) Generate(cmd *cobra.Command, args []string) error {
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
