package service

import (
	"crypto/tls"
	"time"
)

type Options struct {
	addr string

	tor       bool
	pass      string
	uid       string
	crypt     string
	tlsConf   *tls.Config
	keepalive *KeepAliveOpt
	buffers   *BufferOpt
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
