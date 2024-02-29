package transfer

import (
	"bufio"
	"log/slog"
	"net"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/xlog"
)

func Socks5Transport(conn net.Conn) error {
	br := &network.BufferedConn{Conn: conn, Br: bufio.NewReader(conn)}
	b, err := br.Peek(1)

	if err != nil || b[0] != network.SOCKS5Version {
		slog.Error(xlog.CLIENT_PROTOCOL_UNSUPPORT_ERROR, "err", err)
		return err
	} else {
		// select method
		_, err := network.ReadMethods(br)
		if err != nil {
			slog.Error(`[socks5] read methods failed`, "err", err)
			return err
		}

		// TODO AUTH
		if err := network.WriteMethod(network.MethodNoAuth, br); err != nil {
			if err != nil {
				slog.Error(`[socks5] write method failed`, "err", err)
			} else {
				slog.Error(`[socks5] methods is not acceptable`)
			}
			return err
		}

		// read command
		request, err := network.ReadSOCKS5Request(br)
		if err != nil {
			slog.Error(`[socks5] read command failed`, "err", err)
			return err
		}
		switch request.Cmd {
		case network.SOCKS5CmdConnect:
			handleConnect(br, request)
		case network.SOCKS5CmdBind:
			slog.Error("not support cmd bind")
			//handleBind(conn, request)
		case network.SOCKS5CmdUDPOverTCP:
			//handleUDP(conn, request)
			slog.Error("not support cmd upd")
		}
	}
	return nil
}

func handleConnect(conn net.Conn, req *network.SOCKS5Request) {
	slog.Info(xlog.SOCKS5_CONNECT_SERVER, "src", conn.RemoteAddr(), "dest", req.Addr)
	// default socks timeout : 10
	dialer := net.Dialer{Timeout: 10 * time.Second}
	destConn, err := dialer.Dial("tcp", req.Addr.String())

	if err != nil {
		slog.Error(xlog.SOCKS5_CONNECT_DIAL_ERROR, "err", err)
		return
	}

	// if utils.Transport get out , then close conn of remote
	defer destConn.Close()

	if err := network.NewReply(network.SOCKS5RespSucceeded, nil).Write(conn); err != nil {
		slog.Error(xlog.SOCKS5_CONNECT_WRITE_ERROR, "err", err)
		return
	}

	slog.Info(xlog.SOCKS5_CONNECT_ESTAB, "src", conn.RemoteAddr(), "dest", req.Addr)

	if _, _, err := network.Transport(conn, destConn); err != nil {
		slog.Error(xlog.CONNECT_TRANS_ERROR, "err", err)
	}

	slog.Info(xlog.SOCKS5_CONNECT_DIS, "src", conn.RemoteAddr(), "dest", req.Addr)
}

func handleBind() {
	// TODO
}

func handleUDPOverTCP() {
	// TODO
}
