package abstract

import (
	"fmt"
)

// Auth 抽象各类认证接口
type Auth interface {
	Identity() string
}

type AuthList []Auth

func (al AuthList) IsEmpty() bool {
	return al == nil || len(al) == 0
}

type AdminAuth struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	LastLogin string `json:"last_login"`
	LastAddr  string `json:"last_addr"`
}

func (a AdminAuth) Identity() string {
	return fmt.Sprintf("%s:%s", a.Username, a.Password)
}

type CloudKeyAuth struct {
	// 阿里需要 ID
	AccessId string `json:"access_id"`
	// 普通云厂商使用的认证
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
}

type CertsAuth struct {
	// Sealos 认证信息
	KubeConfig string `json:"kube_config" gorm:"not null"`
}

func (a CloudKeyAuth) Identity() string {
	return fmt.Sprintf("%s:%s:%s", a.AccessId, a.AccessKey, a.AccessSecret)
}

func (a CertsAuth) Identity() string {
	return a.KubeConfig
}

type UserPassAuth struct {
	User  string `json:"user"`
	Pass  string `json:"pass"`
	Level *int   `json:"level,omitempty"`
}

func (a UserPassAuth) Identity() string {
	return fmt.Sprintf("%s:%s", a.User, a.Pass)
}

// EmailPassAuth for shadowsocks / trojan
type EmailPassAuth struct {
	Password string `json:"password"`
	Method   string `json:"method,omitempty"`
	Level    int    `json:"level"`
	Email    string `json:"email"`
}

func (a EmailPassAuth) Identity() string {
	return fmt.Sprintf("%s:%s:%s", a.Method, a.Email, a.Password)
}

// IdEncryptAuth used for vmess / vless
type IdEncryptAuth struct {
	Id       string `json:"id"` // vmess / vless
	Level    int    `json:"level"`
	Email    string `json:"email"`
	Security string `json:"security,omitempty"` // vmess client needs
	Flow     string `json:"flow,omitempty"`
}

func (a IdEncryptAuth) Identity() string {
	return fmt.Sprintf("%s:%s:%s", a.Security, a.Email, a.Id)
}
