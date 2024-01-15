package client

import (
	"bufio"
	"context"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type bufferedConn struct {
	net.Conn
	br *bufio.Reader
}

func (c *bufferedConn) Read(b []byte) (int, error) {
	return c.br.Read(b)
}

func Socks5Controller(ctx context.Context, sg *SigGroup) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case <-sg.SocksStartChannel:

			server, err := net.Listen("tcp", Config().Socks5.ListenAddr)
			if err != nil {
				slog.Error(consts.SOCKS5_LISTEN_ERROR, err)
				return
			}
			var proxyAddr string
			for _, p := range Config().ProxyAddr {
				if strings.HasPrefix(p, "ws://") || strings.HasPrefix(p, "wss://") {
					proxyAddr = p
				} else if strings.HasPrefix(p, "socks-proxy") {
					proxyAddr = "ws://" + p
				}
			}
			if proxyAddr == "" {
				slog.Error(consts.PROXY_ADDR_ERROR)
				break
			}
			sg.wg.Add(1)
			go func() {
				NewSocks5Client(c, server, proxyAddr)
				sg.wg.Done()
			}()
		case <-sg.SocksStopChannel:
			slog.Info(consts.SOCKS5_LISTEN_STOP)
			cancel()
		}
	}
}

func NewSocks5Client(ctx context.Context, server net.Listener, proxyAddr string) {
	var closeFlag = false
	slog.Info(consts.SOCKS5_LISTEN_START, "addr", server.Addr())
	slog.Debug(consts.PROXY_ADDR, "proxy", proxyAddr)
	defer func() {
		closeFlag = true
		server.Close()
	}()
	go func() {
		for {
			lock := &sync.Mutex{}
			conn, err := server.Accept()
			if err == nil {
				slog.Debug(consts.SOCKS5_ACCEPT_START, "addr", conn.RemoteAddr())
				br := bufio.NewReader(conn)
				b, err := br.Peek(1)
				if err != nil || b[0] != network.SOCKS5Version {
					slog.Error(consts.CLIENT_PROTOCOL_UNSUPPORT_ERROR, "err", err)
				} else {
					go Socks5Handler(&bufferedConn{conn, br}, proxyAddr, lock)
				}
			} else {
				if closeFlag {
					// except close
					return
				} else {
					slog.Error(consts.SOCKS5_ACCEPT_ERROR, "err", err)
				}
			}
		}
	}()
	<-ctx.Done()
}

func Socks5Handler(conn net.Conn, raddr string, lock *sync.Mutex) {
	// select method
	_, err := network.ReadMethods(conn)
	if err != nil {
		slog.Error(`[socks5] read methods failed`, "err", err)
		return
	}

	// TODO AUTH
	if err := network.WriteMethod(network.MethodNoAuth, conn); err != nil {
		if err != nil {
			slog.Error(`[socks5] write method failed`, "err", err)
		} else {
			slog.Error(`[socks5] methods is not acceptable`)
		}
		return
	}

	// read command
	request, err := network.ReadSOCKS5Request(conn)
	if err != nil {
		slog.Error(`[socks5] read command failed`, "err", err)
		return
	}
	switch request.Cmd {
	case network.SOCKS5CmdConnect:
		handleConnect(conn, request, raddr, lock)
		break
	case network.SOCKS5CmdBind:
		slog.Error("not support cmd bind")
		//handleBind(conn, request)
		break
	case network.SOCKS5CmdUDPOverTCP:
		//handleUDP(conn, request)
		slog.Error("not support cmd upd")
		break
	}
}

func handleConnect(conn net.Conn, req *network.SOCKS5Request, rAddr string, lock *sync.Mutex) {

	slog.Info(consts.SOCKS5_CONNECT_SERVER, "src", conn.RemoteAddr(), "dest", req.Addr)

	dialer := &websocket.Dialer{}
	s := http.Header{}
	s.Set("SM-CMD", "CONNECT")
	s.Set("SM-TARGET", req.Addr.String())
	wsConn, _, err := dialer.Dial(rAddr, s)

	if err != nil {
		slog.Error(consts.SOCKS_UPGRADE_ERROR, "err", err)
		return
	}

	newConn := tunnel.WsWrapConn(wsConn)

	defer newConn.Close()

	if err := network.NewReply(network.SOCKS5RespSucceeded, nil).Write(conn); err != nil {
		slog.Error(consts.SOCKS5_CONNECT_WRITE_ERROR, "err", err)
		return
	}
	slog.Info(consts.SOCKS5_CONNECT_ESTAB, "src", conn.RemoteAddr(), "dest", req.Addr)

	if err := network.Transport(conn, newConn); err != nil {
		slog.Error(consts.SOCKS5_CONNECT_TRANS_ERROR, "err", err)
	}

	slog.Info(consts.SOCKS5_CONNECT_DIS, "src", conn.RemoteAddr(), "dest", req.Addr)

}
