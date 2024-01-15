package server

import "github.com/DVKunion/SeaMoon/pkg/tunnel"

type Options struct {
	host  string      // 监听地址
	port  string      // 监听端口
	proto tunnel.Type // 监听协议
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
		tt := tunnel.TranType(t)
		if tt == tunnel.NULL {
			// todo
		}
		o.proto = tt
		return nil
	}
}
