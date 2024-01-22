package models

import (
	"reflect"
	"time"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
)

type CloudProvider struct {
	gorm.Model

	// 元信息
	Name   *string          `gorm:"not null"`
	Desc   *string          `gorm:"not null"`
	Region *string          `gorm:"not null"`
	Type   *types.CloudType `gorm:"not null"`

	// 账户信息
	Amount   *float64           `gorm:"not null"`
	Cost     *float64           `gorm:"not null"`
	Status   *types.CloudStatus `gorm:"not null"`
	MaxLimit *int               `gorm:"not null"`

	CloudAuth *CloudAuth `gorm:"embedded"`

	// 连表
	Tunnels []Tunnel `gorm:"foreignKey:CloudProviderId;references:ID"`
}

// CloudProviderApi api 不展示敏感的账户数据
type CloudProviderApi struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   *string          `json:"name"`
	Desc   *string          `json:"desc"`
	Type   *types.CloudType `json:"type"`
	Region *string          `json:"region"`

	// 账户信息
	Amount   *float64           `json:"amount"`
	Cost     *float64           `json:"cost"`
	Status   *types.CloudStatus `json:"status"`
	Count    *int               `json:"count"`
	MaxLimit *int               `json:"max_limit"`
}

// CloudProviderCreateApi api 用于创建时接受数据的模型
type CloudProviderCreateApi struct {
	Name   *string            `json:"name"`
	Region *string            `json:"region"`
	Desc   *string            `json:"desc"`
	Status *types.CloudStatus `json:"status"`
	Type   *types.CloudType   `json:"type"`

	// 认证信息
	CloudAuth *CloudAuth `json:"cloud_auth"`
}

func (p CloudProvider) Extra() func(api interface{}) {
	return func(api interface{}) {
		ref := reflect.ValueOf(api).Elem()
		field := ref.FieldByName("Count")
		if field.CanSet() {
			a := len(p.Tunnels)
			field.Set(reflect.ValueOf(&a))
		}
	}
}
