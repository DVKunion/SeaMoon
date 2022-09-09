package server

import (
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type SocksServer struct {
	CloudServer
	request *utils.Request
}

func (s *Server) SocksServerTransfer(r *http.Request) CloudServer {
	return newSocksServer(r)
}

func (s *SocksServer) Verification(w http.ResponseWriter) bool {
	return true
}

func (s *SocksServer) Serve(w http.ResponseWriter, r *http.Request) {
	// socks upgrade websocket
	var upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var conn, _ = upGrader.Upgrade(w, r, nil)
	wsConn := NewWebsocketServer(conn)

	switch s.request.Cmd {
	case utils.CmdConnect:
		s.handleConnect(wsConn)
	case utils.CmdBind:
		s.handleBind()
	case utils.CmdUDPOverTCP:
		s.handleUDPOverTCP()
	default:
		conn.WriteMessage(websocket.BinaryMessage, []byte("UnSupport Command"))
	}
}

func newSocksServer(r *http.Request) *SocksServer {

	cmd := r.Header.Get("SM-CMD")
	target := r.Header.Get("SM-TARGET")
	// TODO some panic
	host := strings.Split(target, ":")[0]
	port := strings.Split(target, ":")[1]

	transCommand := map[string]uint8{
		"CONNECT": utils.CmdConnect,
		"BIND":    utils.CmdBind,
		"UDP":     utils.CmdUDPOverTCP,
	}

	command, ok := transCommand[cmd]

	if !ok {

	}

	reqPort, err := strconv.ParseUint(port, 10, 16)

	if err != nil {

	}

	return &SocksServer{
		request: &utils.Request{
			Cmd: command,
			Addr: &utils.Addr{
				Host: host,
				Port: uint16(reqPort),
				Type: utils.AddrDomain,
			},
		},
	}
}

func (s *SocksServer) handleConnect(conn net.Conn) {
	log.Infof(consts.SOCKS5_CONNECT_SERVER, s.request.Addr, conn.RemoteAddr())
	// 默认socks设置十秒超时
	dialer := net.Dialer{Timeout: 10}
	newConn, err := dialer.Dial("tcp", s.request.Addr.String())

	if err != nil {
		log.Errorf(consts.SOCKS5_CONNECT_DIAL_ERROR, err)
		conn.Close()
		return
	}

	defer newConn.Close()
	defer conn.Close()

	log.Infof(consts.SOCKS5_CONNECT_ESTAB, conn.RemoteAddr(), s.request.Addr)

	if err := utils.Transport(conn, newConn); err != nil {
		log.Errorf(consts.SOCKS5_CONNECT_TRANS_ERROR, err)
	}

	log.Infof(consts.SOCKS5_CONNECT_DIS, conn.RemoteAddr(), s.request.Addr)
}

func (s *SocksServer) handleBind() {

}

func (s *SocksServer) handleUDPOverTCP() {

}
