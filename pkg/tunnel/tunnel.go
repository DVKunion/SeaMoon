package tunnel

import (
	"net"
)

// Tunnel is a bridge implementation for data transmission
// tunnel 不应该解决如何处理流量信息，仅仅是一个桥梁，用于信道传输。
type Tunnel interface {
	net.Conn
}

type Type string

const (
	NULL = "unknown"
	WST  = "websocket-tunnel"
	GRT  = "grpc-tunnel"
)

func TranType(t string) Type {
	switch t {
	case "websocket":
		return WST
	case "grpc":
		return GRT
	default:
		return NULL
	}
}
