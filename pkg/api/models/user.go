package models

import (
	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/types"
)

type Auth struct {
	gorm.Model
	Type     types.AuthType
	UserName string
	Password string
}
