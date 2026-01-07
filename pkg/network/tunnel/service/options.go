package service

import (
	"crypto/tls"
	"time"
)

type Options struct {
	addr string
	path string

	tor       bool
	pass      string
	uid       string
	crypt     string
	udpAddr   string
	tlsConf   *tls.Config
	keepalive *KeepAliveOpt
	buffers   *BufferOpt

	// 级联代理配置
	cascadeAddr     string
	cascadeUid      string
	cascadePassword string
}

type Option func(o *Options)

type KeepAliveOpt struct {
	MinTime           time.Duration
	MaxTime           time.Duration
	Timeout           time.Duration
	HandshakeTimeout  time.Duration
	MaxConnectionIdle time.Duration

	PermitStream bool
}

type BufferOpt struct {
	ReadBufferSize    int
	WriteBufferSize   int
	EnableCompression bool
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.addr = addr
	}
}

func WithPath(path string) Option {
	return func(o *Options) {
		o.path = path
	}
}

func WithTorFlag(tor bool) Option {
	return func(o *Options) {
		o.tor = tor
	}
}

func WithTLSConf(t *tls.Config) Option {
	return func(o *Options) {
		o.tlsConf = t
	}
}

func WithPassword(pass string) Option {
	return func(o *Options) {
		o.pass = pass
	}
}

func WithUid(uid string) Option {
	return func(o *Options) {
		o.uid = uid
	}
}

func WithCrypt(c string) Option {
	return func(o *Options) {
		o.crypt = c
	}
}

func WithUDPAddr(p string) Option {
	return func(o *Options) {
		o.udpAddr = p
	}
}

func WithKeepAlive(k *KeepAliveOpt) Option {
	return func(o *Options) {
		o.keepalive = k
	}
}

func WithBuffers(b *BufferOpt) Option {
	return func(o *Options) {
		o.buffers = b
	}
}

func WithCascadeProxy(addr, uid, password string) Option {
	return func(o *Options) {
		o.cascadeAddr = addr
		o.cascadeUid = uid
		o.cascadePassword = password
	}
}
