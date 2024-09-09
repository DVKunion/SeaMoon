package model

import (
	"encoding/json"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models/abstract"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/keys"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

// Account 表
// 用于存放所有账户、认证相关的信息
type Account struct {
	gorm.Model

	Name     string           // 账户名
	Type     enum.AccountType // 账户类型
	AuthType enum.AuthType    // 认证类型
	Data     string           // 原始加密数据
	Auth     abstract.Auth    `gorm:"-"`
}

// BeforeSave create / update need to transfer data to encrypt
func (a *Account) BeforeSave(*gorm.DB) (err error) {
	data := tools.MarshalString(a.Auth)
	if data == "" {
		err = errors.New("auth save marshall error!")
		return
	}
	a.Data, err = tools.AESEncrypt([]byte(data), keys.GetGlobalKey())
	return
}

// AfterFind select to transfer data to plaintext
func (a *Account) AfterFind(*gorm.DB) error {
	data, err := tools.AESDecrypt(a.Data, keys.GetGlobalKey())
	if err != nil {
		return err
	}
	var b abstract.Auth
	switch a.AuthType {
	case enum.AdminAuth:
		b = &abstract.AdminAuth{}
	case enum.XrayAuthUserPass:
		b = &abstract.UserPassAuth{}
	case enum.XrayAuthEmailPass:
		b = &abstract.EmailPassAuth{}
	case enum.XrayAuthIdEncrypt:
		b = &abstract.IdEncryptAuth{}
	case enum.CloudAuthKeyMod:
		b = &abstract.CloudKeyAuth{}
	case enum.CloudAuthCert:
		b = &abstract.CertsAuth{}
	}
	if b == nil {
		return errors.New("xxxxxx")
	}
	err = json.Unmarshal(data, b)
	if err != nil {
		return err
	}
	a.Auth = b
	return nil
}
