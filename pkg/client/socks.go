package client

import (
	"bufio"
	"context"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/server"
	"github.com/DVKunion/SeaMoon/pkg/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
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
				log.Errorf(consts.SOCKS5_LISTEN_ERROR, err)
				return
			}
			var proxyAddr string
			for _, p := range Config().ProxyAddr {
				if strings.HasPrefix(p, "socks-proxy") {
					proxyAddr = "ws://" + p
				}
			}
			if proxyAddr == "" {
				log.Error(consts.PROXY_ADDR_ERROR)
				break
			}
			sg.wg.Add(1)
			go func() {
				NewSocks5Client(c, server, proxyAddr)
				sg.wg.Done()
			}()
		case <-sg.SocksStopChannel:
			Config().Socks5.Status = "inactive"
			log.Info(consts.SOCKS5_LISTEN_STOP)
			cancel()
			return
		}
	}
}

func NewSocks5Client(ctx context.Context, server net.Listener, proxyAddr string) {
	var closeFlag = false
	log.Infof(consts.SOCKS5_LISTEN_START, server.Addr())
	log.Debugf(consts.PROXY_ADDR, proxyAddr)
	defer func() {
		closeFlag = true
		server.Close()
	}()
	go func() {
		for {
			conn, err := server.Accept()
			if err == nil {
				log.Debugf(consts.SOCKS5_ACCEPT_START, conn.RemoteAddr())
				br := bufio.NewReader(conn)
				b, err := br.Peek(1)
				if err != nil || b[0] != utils.Version {
					log.Errorf(consts.CLIENT_PROTOCOL_UNSUPPORT_ERROR, err)
					return
				}
				Socks5Handler(&bufferedConn{conn, br}, proxyAddr)
			} else {
				if closeFlag {
					// except close
					return
				} else {
					log.Errorf(consts.SOCKS5_ACCEPT_ERROR, err)
				}
			}
		}
	}()
	<-ctx.Done()
}

func Socks5Handler(conn net.Conn, raddr string) {
	// select method
	_, err := utils.ReadMethods(conn)
	if err != nil {
		log.Errorf(`[socks5] read methods failed: %s`, err)
		return
	}

	// TODO AUTH
	if err := utils.WriteMethod(utils.MethodNoAuth, conn); err != nil {
		if err != nil {
			log.Errorf(`[socks5] write method failed: %s`, err)
		} else {
			log.Errorf(`[socks5] methods is not acceptable`)
		}
		return
	}

	// read command
	request, err := utils.ReadRequest(conn)
	if err != nil {
		log.Errorf(`[socks5] read command failed: %s`, err)
		return
	}
	switch request.Cmd {
	case utils.CmdConnect:
		handleConnect(conn, request, raddr)
		break
	case utils.CmdBind:
		log.Error("not support cmd bind")
		//handleBind(conn, request)
		break
	case utils.CmdUDP:
		//handleUDP(conn, request)
		log.Error("not support cmd upd")
		break
	}
}

func handleConnect(conn net.Conn, req *utils.Request, rAddr string) {

	log.Infof(consts.SOCKS5_CONNECT_SERVER, req.Addr, conn.RemoteAddr())

	dialer := &websocket.Dialer{}
	s := http.Header{}
	s.Set("SM-CMD", "CONNECT")
	s.Set("SM-TARGET", req.Addr.String())
	wsConn, _, err := dialer.Dial(rAddr, s)

	if err != nil {
		log.Errorf(consts.SOCKS_UPGRADE_ERROR, err)
		return
	}

	newConn := server.NewWebsocketServer(wsConn)

	defer newConn.Close()

	if err := utils.NewReply(utils.Succeeded, nil).Write(conn); err != nil {
		log.Errorf(consts.SOCKS5_CONNECT_WRITE_ERROR, err)
		return
	}
	log.Infof(consts.SOCKS5_CONNECT_ESTAB, conn.RemoteAddr(), req.Addr)

	if err := utils.Transport(conn, newConn); err != nil {
		log.Errorf(consts.SOCKS5_CONNECT_TRANS_ERROR, err)
	}

	log.Infof(consts.SOCKS5_CONNECT_DIS, conn.RemoteAddr(), req.Addr)

}
