package server

import (
	"errors"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SocksServer struct {
	CloudServer
	request *utils.Request
}

func (s *Server) SocksServerTransfer(r *http.Request) CloudServer {
	return newSocksServer(r)
}

func (s *SocksServer) Verification(w http.ResponseWriter) (bool, error) {
	// check target && port
	if s.request.Addr == nil || s.request.Cmd == 0 {
		var errMsg = "socks remote error"
		return false, errors.New(errMsg)
	}
	return true, nil
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
	defer wsConn.Close()

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

	var ss = &SocksServer{request: &utils.Request{}}

	cmd := r.Header.Get("SM-CMD")
	target := r.Header.Get("SM-TARGET")
	if target != "" && len(strings.Split(target, ":")) > 1 {
		host := strings.Split(target, ":")[0]
		port := strings.Split(target, ":")[1]
		if reqPort, err := strconv.ParseUint(port, 10, 16); err == nil {
			ss.request.Addr = &utils.Addr{
				Host: host,
				Port: uint16(reqPort),
				Type: utils.AddrDomain,
			}
		} else {
			log.Error(err)
		}
	}

	transCommand := map[string]uint8{
		"CONNECT": utils.CmdConnect,
		"BIND":    utils.CmdBind,
		"UDP":     utils.CmdUDPOverTCP,
	}

	if command, ok := transCommand[cmd]; ok {
		ss.request.Cmd = command
	}

	return ss
}

func (s *SocksServer) handleConnect(conn net.Conn) {
	log.Infof(consts.SOCKS5_CONNECT_SERVER, s.request.Addr, conn.RemoteAddr())
	// default socks timeout : 10
	dialer := net.Dialer{Timeout: 10 * time.Second}
	newConn, err := dialer.Dial("tcp", s.request.Addr.String())

	if err != nil {
		log.Errorf(consts.SOCKS5_CONNECT_DIAL_ERROR, err)
		return
	}

	// if utils.Transport get out , then close conn of remote
	defer newConn.Close()

	log.Infof(consts.SOCKS5_CONNECT_ESTAB, conn.RemoteAddr(), s.request.Addr)

	if err := utils.Transport(conn, newConn); err != nil {
		log.Errorf(consts.SOCKS5_CONNECT_TRANS_ERROR, err)
	}

	log.Infof(consts.SOCKS5_CONNECT_DIS, conn.RemoteAddr(), s.request.Addr)
}

func (s *SocksServer) handleBind() {
	// TODO
}

func (s *SocksServer) handleUDPOverTCP() {
	// TODO
}
