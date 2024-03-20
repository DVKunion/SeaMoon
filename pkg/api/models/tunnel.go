package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

// Tunnel 表示着实际部署的一个函数节点
type Tunnel struct {
	gorm.Model

	ProviderId uint

	UniqID        *string            // 唯一性ID,用于 sync 同步时识别出唯一函数与隧道关系
	Name          *string            // 隧道名称，建议英文
	Addr          *string            // 服务地址
	Port          *int32             // 服务端口
	Type          *enum.TunnelType   // 隧道协议类型
	Status        *enum.TunnelStatus // 隧道状态
	StatusMessage *string            // 隧道状态原因，用于展示具体的异常详情

	Config *TunnelConfig `gorm:"embedded"`
	// 连表查询
	Proxies []Proxy `gorm:"foreignKey:TunnelID;references:ID"`
}

type TunnelList []*Tunnel

type TunnelConfig struct {
	// 函数配置
	Region     string        `json:"region"`       // 一个隧道只能是一个区域
	CPU        float32       `json:"cpu"`          // CPU 资源
	Memory     int32         `json:"memory"`       // 内存资源
	Instance   int32         `json:"instance"`     // 最大实例处理数
	FcAuthType enum.AuthType `json:"fc_auth_type"` // 函数认证方式

	TLS bool `json:"tls"` // 是否开启 TLS 传输, 开启后自动使用 wss  协议
	Tor bool `json:"tor"` // 是否开启 Tor 转发
}

type TunnelApi struct {
	ProviderId   uint               `json:"provider_id"`
	ProviderType *enum.ProviderType `json:"provider_type"`

	ID        uint      `json:"id"`
	UniqID    *string   `json:"uniq_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name          *string            `json:"name"`
	Addr          *string            `json:"address"`
	Port          *int32             `json:"port"`
	Type          *enum.TunnelType   `json:"type"`
	Status        *enum.TunnelStatus `json:"status"`
	StatusMessage *string            `json:"status_message"`

	Config *TunnelConfig `json:"tunnel_config"`
}

type TunnelCreateApi struct {
	ID            uint               `json:"id"`
	ProviderId    uint               `json:"provider_id"`
	UniqID        *string            `json:"uniq_id"`
	Name          *string            `json:"name"`
	Port          *int32             `json:"port"`
	Type          *enum.TunnelType   `json:"type"`
	Status        *enum.TunnelStatus `json:"status"`
	StatusMessage *string            `json:"status_message"`
	Addr          *string            `json:"address"`

	Config *TunnelConfig `json:"tunnel_config"`
}

type TunnelCreateApiList []*TunnelCreateApi

func (t Tunnel) GetAddr() string {
	switch *t.Type {
	case enum.TunnelTypeWST:
		if t.Config.TLS {
			return "wss://" + *t.Addr
		}
		return "ws://" + *t.Addr
	case enum.TunnelTypeGRT:
		if t.Config.TLS {
			return "grpcs://" + *t.Addr
		}
		return "grpc://" + *t.Addr
	}
	return ""
}

func (t Tunnel) ToApi(extra ...func(api interface{})) *TunnelApi {
	return toApi(t, &TunnelApi{}, extra...).(*TunnelApi)
}

func (tl TunnelList) ToApi(extra ...func(api interface{})) []*TunnelApi {
	res := make([]*TunnelApi, 0)
	for _, t := range tl {
		api := toApi(t, &TunnelApi{}, extra...)
		res = append(res, api.(*TunnelApi))
	}
	return res
}

func (ta TunnelCreateApi) ToModel() *Tunnel {
	return toModel(ta, &Tunnel{}, true).(*Tunnel)
}

func NewTunnelCreateApi() *TunnelCreateApi {
	res := &TunnelCreateApi{}
	autoFull(res)
	return res
}
