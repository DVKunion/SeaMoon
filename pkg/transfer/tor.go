package transfer

import (
	"log/slog"
	"net"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/xlog"
)

const defaultTorAddr = "127.0.0.1:9050"

func TorTransport(conn net.Conn) error {
	// tor 转发非常简单，但是要求入口流量必须是一个 s5，然后直接把 s5 的口子转发给 tor 服务即可。
	dialer := net.Dialer{Timeout: 10 * time.Second}
	destConn, err := dialer.Dial("tcp", defaultTorAddr)

	if err != nil {
		return err
	}

	defer destConn.Close()

	slog.Info(xlog.SOCKS5_CONNECT_ESTAB, "src", conn.RemoteAddr(), "dest", defaultTorAddr)

	if _, _, err := network.Transport(conn, destConn); err != nil {
		slog.Error(xlog.CONNECT_TRANS_ERROR, "err", err)
	}

	slog.Info(xlog.SOCKS5_CONNECT_DIS, "src", conn.RemoteAddr(), "dest", defaultTorAddr)

	return nil
}
