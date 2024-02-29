package models

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/pkg/transfer"
)

type Proxy struct {
	gorm.Model

	TunnelID   uint
	Name       *string
	Type       *transfer.Type
	Status     *types.ProxyStatus
	Conn       *int
	Speed      *float64
	Lag        *int
	InBound    *int64
	OutBound   *int64
	ListenAddr *string
	ListenPort *string
}

type ProxyApi struct {
	ID         uint               `json:"id"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Name       *string            `json:"name"`
	Type       *transfer.Type     `json:"type"`
	Status     *types.ProxyStatus `json:"status"`
	Conn       *int               `json:"conn"`
	Speed      *float64           `json:"speed"`
	Lag        *int               `json:"lag"`
	InBound    *int64             `json:"in_bound"`
	OutBound   *int64             `json:"out_bound"`
	ListenAddr *string            `json:"listen_address"`
	ListenPort *string            `json:"listen_port"`
}

type ProxyCreateApi struct {
	Name            *string            `json:"name"`
	Type            *transfer.Type     `json:"type"`
	ListenAddr      *string            `json:"listen_address"`
	ListenPort      *string            `json:"listen_port"`
	Status          *types.ProxyStatus `json:"status"`
	TunnelId        uint               `json:"tunnel_id"`
	TunnelCreateApi *TunnelCreateApi   `json:"tunnel_create_api"`
}

func (p Proxy) Addr() string {
	return strings.Join([]string{*p.ListenAddr, *p.ListenPort}, ":")
}
