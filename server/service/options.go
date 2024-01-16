package service

import (
	"crypto/tls"
	"time"
)

type Options struct {
	addr string

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

func WithTLSConf(t *tls.Config) Option {
	return func(o *Options) {
		o.tlsConf = t
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
