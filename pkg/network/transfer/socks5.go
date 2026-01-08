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
	// 检查是否启用级联代理，如果启用，直接把流量转发给下一跳处理
	// 必须在 Socks5Check 之前检查，因为级联代理模式下应该让下一跳处理协议
	if IsCascadeEnabled() {
		return CascadeTransport(conn, "socks5")
	}

	var err error
	if !check {
		if conn, err = Socks5Check(conn); err != nil {
			return err
		}
	}

	// 原有逻辑：本节点处理 socks5 协议
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
