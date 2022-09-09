package client

import (
	"bufio"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/server"
	"github.com/DVKunion/SeaMoon/pkg/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type bufferedConn struct {
	net.Conn
	br *bufio.Reader
}

func (c *bufferedConn) Read(b []byte) (int, error) {
	return c.br.Read(b)
}

func NewSocks5Client(listenAddr string, proxyAddr string, verbose bool) {
	server, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Error(consts.SOCKS5_LISTEN_ERROR, err)
	}
	log.Infof(consts.SOCKS5_LISTEN_START, listenAddr)
	log.Debugf(consts.PROXY_ADDR, proxyAddr)
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Error(consts.SOCKS5_ACCEPT_ERROR)
			// 连接的error可能是非预期不可控制的，不干扰client进程。
			continue
		}
		log.Debugf(consts.SOCKS5_ACCEPT_START, conn.RemoteAddr())

		go func() {
			br := bufio.NewReader(conn)
			b, err := br.Peek(1)
			if err != nil || b[0] != utils.Version {
				conn.Close()
				log.Errorf(consts.CLIENT_PROTOCOL_UNSUPPORT_ERROR, err)
				return
			}
			Socks5Handler(&bufferedConn{conn, br}, proxyAddr)
		}()
	}
}

func Socks5Handler(conn net.Conn, raddr string) {
	defer conn.Close()

	// select method
	_, err := utils.ReadMethods(conn)
	if err != nil {
		log.Printf(`[socks5] read methods failed: %s`, err)
		return
	}

	// TODO AUTH
	if err := utils.WriteMethod(utils.MethodNoAuth, conn); err != nil {
		if err != nil {
			log.Printf(`[socks5] write method failed: %s`, err)
		} else {
			log.Printf(`[socks5] methods is not acceptable`)
		}
		return
	}

	// read command
	request, err := utils.ReadRequest(conn)
	if err != nil {
		log.Printf(`[socks5] read command failed: %s`, err)
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
	// TODO fix with target format check
	s.Set("SM-TARGET", req.Addr.String())
	wsConn, _, err := dialer.Dial(rAddr, s)

	if err != nil {
		log.Errorf("websockect connect error: %s ", err)
		return
	}

	newConn := server.NewWebsocketServer(wsConn)

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
