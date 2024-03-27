package server

import (
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

type Options struct {
	host  string          // 监听地址
	port  string          // 监听端口
	proto enum.TunnelType // 监听协议

	mtcp bool //
}

type Option func(o *Options) (err error)

func WithHost(host string) Option {
	return func(o *Options) (err error) {
		o.host = host
		return
	}
}

func WithPort(port string) Option {
	return func(o *Options) (err error) {
		o.port = port
		return
	}
}

func WithProto(t string) Option {
	return func(o *Options) (err error) {
		tt := enum.TransTunnelType(t)
		if tt == enum.TunnelTypeNULL {
			// todo
		}
		o.proto = tt
		return nil
	}
}

func WithMTcp(flag bool) Option {
	return func(o *Options) (err error) {
		o.mtcp = flag
		return nil
	}
}
