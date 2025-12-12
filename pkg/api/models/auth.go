package models

import (
	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

var DefaultAuth = []Auth{
	{
		Type:     enum.AuthAdmin,
		Username: "seamoon",
		Password: "2575a6f37310dd27e884a0305a2dd210",
	},
}

type Auth struct {
	gorm.Model

	Type     enum.AuthType // 认证类型，用于判断该认证信息适用于啥的
	Username string        `json:"username"`
	Password string        `json:"password"`

	LastLogin string `json:"last_login"`
	LastAddr  string `json:"last_addr"`
}

type AuthApi struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CloudAuth struct {
	// 普通云厂商使用的认证
	AccessKey    string `json:"access_key" gorm:"not null"`
	AccessSecret string `json:"access_secret" gorm:"not null"`

	// 接口类型认证信息
	Token string `json:"token" gorm:"not null"`

	// Sealos 认证信息
	KubeConfig string `json:"kube_config" gorm:"not null"`
}
