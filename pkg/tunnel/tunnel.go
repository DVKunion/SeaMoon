package tunnel

import (
	"net"
)

// Tunnel is a bridge implementation for data transmission
// tunnel 不应该解决如何处理流量信息，仅仅是一个桥梁，用于信道传输。
type Tunnel interface {
	net.Conn
}

type Status int8

const (
	INITIALIZING Status = iota + 1 // 初始化
	ACTIVE                         // 可用
	INACTIVE                       // 停用
	ERROR                          // 不可用
	WAITING                        // 异常
)

type Type string

const (
	NULL = "unknown"
	WST  = "websocket-tunnel"
	GRT  = "grpc-tunnel"
)

var tpMaps = map[Type]string{
	NULL: "",
	WST:  "ws",
	GRT:  "grpc",
}

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

func (t Type) TranAddr(tls bool) string {
	proto := tpMaps[t]
	if tls {
		return proto + "s://"
	}
	return proto + "://"
}
