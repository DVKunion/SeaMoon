package transfer

import (
	"bufio"
	"net"

	"github.com/DVKunion/SeaMoon/pkg/network/basic"
)

// AutoTransport 自适应解析 http / socks
func AutoTransport(conn net.Conn) error {
	// 检查是否启用级联代理，如果启用，直接把流量转发给下一跳处理
	if IsCascadeEnabled() {
		return CascadeTransport(conn, "auto")
	}

	br := &basic.BufferedConn{Conn: conn, Br: bufio.NewReader(conn)}
	b, err := br.Peek(1)

	if err != nil || b[0] != basic.SOCKS5Version {
		return HttpTransport(br)
	}

	return Socks5Transport(br, true, conn.LocalAddr().String())
}

func AutoTransportV2ray(conn net.Conn) error {
	br := &basic.BufferedConn{Conn: conn, Br: bufio.NewReader(conn)}
	b, err := br.Peek(1)

	if err != nil || b[0] != basic.SOCKS5Version {
		return V2rayTransport(br, "http")
	}

	return V2rayTransport(br, "socks5")
}
