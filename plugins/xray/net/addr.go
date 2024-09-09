package net

import (
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/infra/conf"
)

type Address net.Address

func A2Cfg(addr Address) *conf.Address {
	return &conf.Address{
		Address: addr,
	}
}

func ParseAddress(addr string) Address {
	return net.ParseAddress(addr)
}
