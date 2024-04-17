package transfer

import (
	"bufio"
	"net"

	"github.com/DVKunion/SeaMoon/pkg/network/basic"
)

// AutoTransport 自适应解析 http / socks
func AutoTransport(conn net.Conn) error {
	br := &basic.BufferedConn{Conn: conn, Br: bufio.NewReader(conn)}
	b, err := br.Peek(1)

	if err != nil || b[0] != basic.SOCKS5Version {
		return HttpTransport(br)
	}

	return Socks5Transport(br, true)
}

func AutoTransportV2ray(conn net.Conn) error {
	br := &basic.BufferedConn{Conn: conn, Br: bufio.NewReader(conn)}
	b, err := br.Peek(1)

	if err != nil || b[0] != basic.SOCKS5Version {
		return V2rayTransport(br, "http")
	}

	return V2rayTransport(br, "socks5")
}
