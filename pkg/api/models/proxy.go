package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

type Proxy struct {
	gorm.Model

	TunnelID      uint
	Name          *string
	Type          *enum.ProxyType
	Status        *enum.ProxyStatus
	StatusMessage *string
	Conn          *int
	SpeedUp       *float64
	SpeedDown     *float64
	Lag           *int64
	InBound       *int64
	OutBound      *int64
	ListenAddr    *string
	ListenPort    *string
}

type ProxyList []*Proxy

type ProxyApi struct {
	ID            uint              `json:"id"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Name          *string           `json:"name"`
	Type          *enum.ProxyType   `json:"type"`
	Status        *enum.ProxyStatus `json:"status"`
	StatusMessage *string           `json:"status_message"`
	Conn          *int              `json:"conn"`
	SpeedUp       *float64          `json:"speed_up"`
	SpeedDown     *float64          `json:"speed_down"`
	Lag           *int64            `json:"lag"`
	InBound       *int64            `json:"in_bound"`
	OutBound      *int64            `json:"out_bound"`
	ListenAddr    *string           `json:"listen_address"`
	ListenPort    *string           `json:"listen_port"`
}

type ProxyCreateApi struct {
	ID              uint              `json:"id"`
	Name            *string           `json:"name"`
	Type            *enum.ProxyType   `json:"type"`
	ListenAddr      *string           `json:"listen_address"`
	ListenPort      *string           `json:"listen_port"`
	Status          *enum.ProxyStatus `json:"status"`
	StatusMessage   *string           `json:"status_message"`
	TunnelId        uint              `json:"tunnel_id"`
	TunnelCreateApi *TunnelCreateApi  `json:"tunnel_create_api"`
}

func (p Proxy) Addr() string {
	return strings.Join([]string{*p.ListenAddr, *p.ListenPort}, ":")
}

func (p Proxy) ProtoAddr() string {
	return fmt.Sprintf("%s://%s", p.Type, strings.Join([]string{*p.ListenAddr, *p.ListenPort}, ":"))
}

func (p Proxy) ToApi() *ProxyApi {
	return toApi(p, &ProxyApi{}).(*ProxyApi)
}

func (pl ProxyList) ToApi() []*ProxyApi {
	res := make([]*ProxyApi, 0)
	for _, d := range pl {
		api := toApi(d, &ProxyApi{})
		res = append(res, api.(*ProxyApi))
	}
	return res
}

func (pa ProxyCreateApi) ToModel() *Proxy {
	return toModel(pa, &Proxy{}, true).(*Proxy)
}
