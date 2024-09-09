package model

import (
	"encoding/json"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models/abstract"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
)

type Function struct {
	gorm.Model

	ProviderId    uint
	UniqID        *string                 // 唯一性ID,用于 sync 同步时识别出唯一函数与隧道关系
	Name          *string                 // 函数名称，建议英文
	Type          enum.FunctionType       // 隧道协议类型
	Status        enum.FunctionStatus     // 函数状态
	FcAuthType    enum.AuthType           // 函数认证方式
	StatusMessage *string                 // 函数状态原因，用于展示具体的异常详情
	ConfigRaw     string                  // 函数配置存储
	Config        abstract.FunctionConfig `gorm:"-"` // 实际类型
}

// BeforeSave create / update need to transfer data to encrypt
func (f *Function) BeforeSave(*gorm.DB) error {
	data, err := json.Marshal(f.Config)
	if err != nil {
		return err
	}
	f.ConfigRaw = string(data)
	return nil
}

// AfterFind select to transfer data to plaintext
func (f *Function) AfterFind(*gorm.DB) error {
	var b abstract.FunctionConfig
	switch f.Type {
	case enum.FunctionTunnel:
		b = &abstract.TunnelConfig{}
	}
	if b == nil {
		return errors.New("xxxxxx")
	}
	err := json.Unmarshal([]byte(f.ConfigRaw), b)
	if err != nil {
		return err
	}
	f.Config = b
	return nil
}
