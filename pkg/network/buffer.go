package network

// fork from go-gost/core

import (
	"bufio"
	"net"
	"sync"
)

type BufferedConn struct {
	net.Conn
	Br *bufio.Reader
}

func (c *BufferedConn) Read(b []byte) (int, error) {
	return c.Br.Read(b)
}

func (c *BufferedConn) Peek(n int) ([]byte, error) {
	return c.Br.Peek(n)
}

var (
	pools = []struct {
		size int
		pool sync.Pool
	}{
		{
			size: 128,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 128)
					return b
				},
			},
		},
		{
			size: 512,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 512)
					return b
				},
			},
		},
		{
			size: 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 1024)
					return b
				},
			},
		},
		{
			size: 2048,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 2048)
					return b
				},
			},
		},
		{
			size: 4096,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 4096)
					return b
				},
			},
		},
		{
			size: 8192,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 8192)
					return b
				},
			},
		},
		{
			size: 16 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 16*1024)
					return b
				},
			},
		},
		{
			size: 32 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 32*1024)
					return b
				},
			},
		},
		{
			size: 64 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 64*1024)
					return b
				},
			},
		},
		{
			size: 65 * 1024,
			pool: sync.Pool{
				New: func() any {
					b := make([]byte, 65*1024)
					return b
				},
			},
		},
	}
)

// GetBuffer returns a buffer of specified size.
func GetBuffer(size int) []byte {
	for i := range pools {
		if size <= pools[i].size {
			b := pools[i].pool.Get().([]byte)
			return b[:size]
		}
	}
	b := make([]byte, size)
	return b
}

func PutBuffer(b []byte) {
	for i := range pools {
		if cap(b) == pools[i].size {
			pools[i].pool.Put(b)
		}
	}
}
