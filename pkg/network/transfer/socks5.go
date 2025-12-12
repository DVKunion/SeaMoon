package transfer

import (
	"bufio"
	"net"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/network/basic"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func Socks5Check(conn net.Conn) (net.Conn, error) {
	br := &basic.BufferedConn{Conn: conn, Br: bufio.NewReader(conn)}
	b, err := br.Peek(1)

	if err != nil || b[0] != basic.SOCKS5Version {
		return nil, errors.Wrap(err, xlog.ServiceProtocolNotSupportError)
	}
	return br, nil
}

func Socks5Transport(conn net.Conn, check bool, udpAddr string) error {

	var err error
	if !check {
		if conn, err = Socks5Check(conn); err != nil {
			return err
		}
	}
	// todo AUTH

	// select method
	if _, err = basic.ReadMethods(conn); err != nil {
		return errors.Wrap(err, xlog.ServiceSocks5ReadMethodError)
	}

	if err = basic.WriteMethod(basic.MethodNoAuth, conn); err != nil {
		return errors.Wrap(err, xlog.ServiceSocks5WriteMethodError)
	}

	// read command
	request, err := basic.ReadSOCKS5Request(conn)
	if err != nil {
		return errors.Wrap(err, xlog.ServiceSocks5ReadCmdError)
	}
	switch request.Cmd {
	case basic.SOCKS5CmdConnect:
		handleConnect(conn, request)
	case basic.SOCKS5CmdBind:
		// todo: support cmd bind
		xlog.Debug("unexpect not support cmd bind")
		handleBind(conn, request)
	case basic.SOCKS5CmdUDPOverTCP:
		handleUDPOverTCP(conn, request, udpAddr)
	}

	return nil
}

func handleConnect(conn net.Conn, req *basic.SOCKS5Request) {
	xlog.Info(xlog.ServiceSocks5ConnectServer, "src", conn.RemoteAddr(), "dest", req.Addr)
	// default socks timeout : 10
	dialer := net.Dialer{Timeout: 10 * time.Second}
	destConn, err := dialer.Dial("tcp", req.Addr.String())

	if err != nil {
		xlog.Error(xlog.ServiceSocks5DailError, "err", err)
		return
	}

	// if utils.Transport get out , then close conn of remote
	defer destConn.Close()

	if err := basic.NewReply(basic.SOCKS5RespSucceeded, nil).Write(conn); err != nil {
		xlog.Error(xlog.ServiceSocks5ReplyError, "err", err)
		return
	}

	xlog.Info(xlog.ServiceSocks5Establish, "src", conn.RemoteAddr(), "dest", req.Addr)

	if _, _, err := basic.Transport(conn, destConn); err != nil {
		xlog.Error(xlog.NetworkTransportError, "err", err)
	}

	xlog.Info(xlog.ServiceSocks5DisConnect, "src", conn.RemoteAddr(), "dest", req.Addr)
}

func handleBind(conn net.Conn, req *basic.SOCKS5Request) {
	// TODO
}

func handleUDPOverTCP(conn net.Conn, req *basic.SOCKS5Request, udpAddr string) {
	xlog.Info(xlog.ServiceSocks5ConnectServer, "src", conn.RemoteAddr(), "dest", req.Addr)

	localAddr, err := basic.NewAddr(udpAddr)
	if err != nil {
		xlog.Error(xlog.ServiceSocks5ReplyError, "err", err)
		return
	}

	if err := basic.NewReply(basic.SOCKS5RespSucceeded, localAddr).Write(conn); err != nil {
		xlog.Error(xlog.ServiceSocks5ReplyError, "err", err)
		return
	}
}
