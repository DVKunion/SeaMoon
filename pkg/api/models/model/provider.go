package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

type Provider struct {
	gorm.Model

	// 元信息
	Name    *string            `gorm:"not null"`
	Desc    *string            `gorm:"not null"`
	Regions Regions            `gorm:"not null"`
	Type    *enum.ProviderType `gorm:"not null"`

	Status        *enum.ProviderStatus `gorm:"not null"`
	StatusMessage *string              `gorm:"not null"`
	MaxLimit      *int                 `gorm:"not null"`

	Info *ProviderInfo `gorm:"embedded"`

	// 连表
	Tunnels []Function `gorm:"foreignKey:ProviderId;references:ID"`
}

type ProviderInfo struct {
	// 账户详细信息
	Amount *float64 `json:"amount"`
	Cost   *float64 `json:"cost"`
}

type ProviderList []*Provider

// ProviderApi api 不展示敏感的账户数据
type ProviderApi struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name    *string            `json:"name"`
	Desc    *string            `json:"desc"`
	Type    *enum.ProviderType `json:"type"`
	Regions Regions            `json:"regions"`

	// 账户信息
	Info          *ProviderInfo        `json:"info"`
	Status        *enum.ProviderStatus `json:"status"`
	StatusMessage *string              `json:"status_message"`
	Count         *int                 `json:"count"`
	MaxLimit      *int                 `json:"max_limit"`
}

// ProviderCreateApi api 用于创建时接受数据的模型
type ProviderCreateApi struct {
	ID      uint                 `json:"id"`
	Name    *string              `json:"name"`
	Regions Regions              `json:"regions"`
	Desc    *string              `json:"desc"`
	Status  *enum.ProviderStatus `json:"status"`
	Type    *enum.ProviderType   `json:"type"`

	// 认证信息
	CloudAuth *Account `json:"cloud_auth"`
}

func (pl ProviderList) ToApi() []*ProviderApi {
	res := make([]*ProviderApi, 0)
	for _, d := range pl {
		api := toApi(d, &ProviderApi{}, d.extra())
		res = append(res, api.(*ProviderApi))
	}
	return res
}

func (p Provider) ToApi() *ProviderApi {
	return toApi(p, &ProviderApi{}, p.extra()).(*ProviderApi)
}

func (pa ProviderCreateApi) ToModel(full bool) *Provider {
	return toModel(pa, &Provider{}, full).(*Provider)
}

func (p Provider) extra() func(api interface{}) {
	return func(api interface{}) {
		ref := reflect.ValueOf(api).Elem()
		field := ref.FieldByName("Count")
		if field.CanSet() {
			a := len(p.Tunnels)
			field.Set(reflect.ValueOf(&a))
		}
	}
}

// Regions
// gorm 不支持 []string 这种类型，因此需要自定义一个字段类型实现 text <-> jsonString 的转换
// 这里两个方法一个指针一个非指针是不能动的。。。不知道为什么，动了就炸。
type Regions []string

func (sl *Regions) Scan(value any) error {
	if b, ok := value.(string); ok {
		return json.Unmarshal([]byte(b), &sl)
	}
	return errors.New("can not transfer")
}
func (sl Regions) Value() (driver.Value, error) {
	data, err := json.Marshal(sl)
	if data == nil {
		return "", err
	}
	return string(data), err
}
