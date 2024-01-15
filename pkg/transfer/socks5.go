package transfer

import (
	"errors"
	"log/slog"
	"net"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/network"
)

func Socks5Transport(conn net.Conn, req *network.SOCKS5Request) error {
	switch req.Cmd {
	case network.SOCKS5CmdConnect:
		handleConnect(conn, req)
	case network.SOCKS5CmdBind:
		handleBind()
	case network.SOCKS5CmdUDP:
		handleUDPOverTCP()
	case network.SOCKS5CmdUDPOverTCP:
		handleUDPOverTCP()
	default:
		return errors.New("")
	}
	return errors.New("")
}

func handleConnect(conn net.Conn, req *network.SOCKS5Request) {
	slog.Info(consts.SOCKS5_CONNECT_SERVER, "src", conn.RemoteAddr(), "dest", req.Addr)
	// default socks timeout : 10
	dialer := net.Dialer{Timeout: 10 * time.Second}
	destConn, err := dialer.Dial("tcp", req.Addr.String())

	if err != nil {
		slog.Error(consts.SOCKS5_CONNECT_DIAL_ERROR, "err", err)
		return
	}

	// if utils.Transport get out , then close conn of remote
	defer destConn.Close()

	slog.Info(consts.SOCKS5_CONNECT_ESTAB, "src", conn.RemoteAddr(), "dest", req.Addr)

	if err := network.Transport(conn, destConn); err != nil {
		slog.Error(consts.SOCKS5_CONNECT_TRANS_ERROR, err)
	}

	slog.Info(consts.SOCKS5_CONNECT_DIS, "src", conn.RemoteAddr(), "dest", req.Addr)
}

func handleBind() {
	// TODO
}

func handleUDPOverTCP() {
	// TODO
}
