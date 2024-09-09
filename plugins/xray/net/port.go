package net

import (
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/infra/conf"
)

type Port *net.Port

func PortFromInt(val uint32) Port {
	p, err := net.PortFromInt(val)
	if err != nil {
		p = net.PortFromBytes([]byte("0"))
	}
	return &p
}
func PortFromString(val string) Port {
	p, err := net.PortFromString(val)
	if err != nil {
		p = net.PortFromBytes([]byte("0"))
	}
	return &p
}

func P2Cfg(p Port) *conf.PortList {
	pr := net.SinglePortRange(*p)
	return &conf.PortList{Range: []conf.PortRange{
		{
			From: pr.From,
			To:   pr.To,
		},
	}}
}
