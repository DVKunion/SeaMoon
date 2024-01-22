package models

import "gorm.io/gorm"

// SystemConfig 系统标准配置表
type SystemConfig struct {
	gorm.Model

	Key   string
	Value string
}
