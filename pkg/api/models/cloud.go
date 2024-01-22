package models

import (
	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/types"
)

type CloudProvider struct {
	gorm.Model

	Type types.CloudType

	// 普通云厂商使用的认证
	AccessKey string
	SecretKey string

	// Sealos 认证信息
	KubeConfig string
}
