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

	// 检查是否启用级联代理
	if IsCascadeEnabled() {
		// 启用级联代理时，通过 v2ray 进行转发
		if err := basic.NewReply(basic.SOCKS5RespSucceeded, nil).Write(conn); err != nil {
			xlog.Error(xlog.ServiceSocks5ReplyError, "err", err)
			return
		}
		xlog.Info(xlog.ServiceSocks5Establish, "src", conn.RemoteAddr(), "dest", req.Addr, "cascade", "enabled")

		// 使用 v2ray 进行级联转发
		if err := cascadeTransport(conn, req.Addr.String()); err != nil {
			xlog.Error(xlog.NetworkTransportError, "err", err, "cascade", "enabled")
		}
		xlog.Info(xlog.ServiceSocks5DisConnect, "src", conn.RemoteAddr(), "dest", req.Addr, "cascade", "enabled")
		return
	}

	// 原有逻辑：直接连接目标
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

// cascadeTransport 通过 v2ray 进行级联代理转发
func cascadeTransport(conn net.Conn, targetAddr string) error {
	// 当启用级联代理时，所有流量将通过 v2ray 的 outbound 进行转发
	// v2ray 的配置在服务启动时已经配置好了级联 vless outbound
	// 这里直接使用 v2ray 进行转发
	if v2ray == nil {
		// 如果 v2ray 未初始化，回退到直连模式
		dialer := net.Dialer{Timeout: 10 * time.Second}
		destConn, err := dialer.Dial("tcp", targetAddr)
		if err != nil {
			return err
		}
		defer destConn.Close()
		_, _, err = basic.Transport(conn, destConn)
		return err
	}

	// 使用 v2ray 进行转发（通过级联 vless 协议）
	return V2rayTransport(conn, "vless")
}
