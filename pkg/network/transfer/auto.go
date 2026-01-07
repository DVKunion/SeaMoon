package transfer

import (
	"bufio"
	"net"

	"github.com/DVKunion/SeaMoon/pkg/network/basic"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

// AutoTransport 自适应解析 http / socks
func AutoTransport(conn net.Conn) error {
	br := &basic.BufferedConn{Conn: conn, Br: bufio.NewReader(conn)}
	b, err := br.Peek(1)

	// 如果启用了级联代理，记录日志
	if IsCascadeEnabled() {
		xlog.Debug("Cascade proxy enabled for auto transport")
	}

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
