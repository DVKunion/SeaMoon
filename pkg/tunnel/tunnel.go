package tunnel

import (
	"net"
)

// Tunnel is a bridge implementation for data transmission
// tunnel 不应该解决如何处理流量信息，仅仅是一个桥梁，用于信道传输。
type Tunnel interface {
	net.Conn
	Delay() int64 // 计算延迟
}
