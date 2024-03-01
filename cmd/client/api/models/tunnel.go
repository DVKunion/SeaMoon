package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

// Tunnel 表示着实际部署的一个函数节点
type Tunnel struct {
	gorm.Model

	CloudProviderId uint

	Name   *string        // 隧道名称，建议英文
	Addr   *string        // 服务地址
	Port   *string        // 服务端口
	Type   *tunnel.Type   // 隧道协议类型
	Status *tunnel.Status // 隧道状态

	TunnelConfig *TunnelConfig `gorm:"embedded"`
	// 连表查询
	Proxies []Proxy `gorm:"foreignKey:TunnelID;references:ID"`
}

func (t *Tunnel) GetAddr() string {
	switch *t.Type {
	case tunnel.WST:
		if t.TunnelConfig.TLS {
			return "wss://" + *t.Addr
		}
		return "ws://" + *t.Addr
	case tunnel.GRT:
		if t.TunnelConfig.TLS {
			return "grpcs://" + *t.Addr
		}
		return "grpc://" + *t.Addr
	}
	return ""
}

type TunnelConfig struct {
	// 函数配置
	CPU            float32        `json:"cpu"`              // CPU 资源
	Memory         int32          `json:"memory"`           // 内存资源
	Instance       int32          `json:"instance"`         // 最大实例处理数
	TunnelAuthType types.AuthType `json:"tunnel_auth_type"` // 函数认证方式

	TLS bool `json:"tls"` // 是否开启 TLS 传输, 开启后自动使用 wss  协议
	Tor bool `json:"tor"` // 是否开启 Tor 转发
}

type TunnelApi struct {
	CloudProviderId     uint             `json:"cloud_provider_id"`
	CloudProviderRegion *string          `json:"cloud_provider_region"`
	CloudProviderType   *types.CloudType `json:"cloud_provider_type"`

	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   *string        `json:"name"`
	Addr   *string        `json:"address"`
	Port   *string        `json:"port"`
	Type   *tunnel.Type   `json:"type"`
	Status *tunnel.Status `json:"status"`

	TunnelConfig *TunnelConfig `json:"tunnel_config"`
}

type TunnelCreateApi struct {
	CloudProviderId uint           `json:"cloud_provider_id"`
	Name            *string        `json:"name"`
	Port            *string        `json:"port"`
	Type            *tunnel.Type   `json:"type"`
	Status          *tunnel.Status `json:"status"`

	TunnelConfig *TunnelConfig `json:"tunnel_config"`
}
