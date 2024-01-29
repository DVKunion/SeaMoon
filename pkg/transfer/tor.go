package transfer

import (
	"log/slog"
	"net"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/network"
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

	slog.Info(consts.SOCKS5_CONNECT_ESTAB, "src", conn.RemoteAddr(), "dest", defaultTorAddr)

	if err := network.Transport(conn, destConn); err != nil {
		slog.Error(consts.CONNECT_TRANS_ERROR, "err", err)
	}

	slog.Info(consts.SOCKS5_CONNECT_DIS, "src", conn.RemoteAddr(), "dest", defaultTorAddr)

	return nil
}
