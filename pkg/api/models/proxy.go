package models

import (
	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/types"
)

type Proxy struct {
	gorm.Model
	ListenAddr string            `json:"listen_addr"` // 代理监听地址
	ListenPort string            `json:"listen_port"` // 代理监听端口
	Speed      string            `json:"speed"`       // 代理速度
	Type       types.ProxyType   `json:"type"`        // 代理类型
	Status     types.ProxyStatus `json:"status"`      // 代理状态
}

// Tunnel 表示着实际部署的一个函数节点
type Tunnel struct {
	gorm.Model
	Addr string           // 服务地址
	Type types.TunnelType // 隧道协议类型
}
