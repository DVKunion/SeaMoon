package service

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/transfer"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

const (
	defaultTimeout        = 10 * time.Second
	defaultReadBufferSize = 32 * 1024
)

type WSService struct {
	upGrader *websocket.Upgrader
}

func init() {
	register(tunnel.WST, &WSService{})
}

func (s *WSService) Handle(m *http.ServeMux) {
	s.upGrader = &websocket.Upgrader{
		HandshakeTimeout:  defaultTimeout,
		ReadBufferSize:    defaultReadBufferSize,
		WriteBufferSize:   defaultReadBufferSize,
		EnableCompression: true,
		CheckOrigin:       func(r *http.Request) bool { return true },
	}
	// websocket http proxy handler
	m.HandleFunc("/http", s.http)

	// websocket socks5 proxy handler
	m.HandleFunc("/socks5", s.socks5)
}

func (s *WSService) http(w http.ResponseWriter, r *http.Request) {
	// means use http to connector
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	t := tunnel.WsWrapConn(conn)
	go func() {
		if err := transfer.HttpTransport(t); err != nil {
			slog.Error("connection error", "msg", err)
		}
	}()
}

func (s *WSService) socks5(w http.ResponseWriter, r *http.Request) {

	cmd := r.Header.Get("SM-CMD")
	target := r.Header.Get("SM-TARGET")

	addr, _ := network.NewAddr(target)

	request := &network.SOCKS5Request{
		Addr: addr,
	}

	transCommand := map[string]uint8{
		"CONNECT": network.SOCKS5CmdConnect,
		"BIND":    network.SOCKS5CmdBind,
		"UDP":     network.SOCKS5CmdUDPOverTCP,
	}

	if command, ok := transCommand[cmd]; ok {
		request.Cmd = command
	}
	// means use socks5 to connector
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	t := tunnel.WsWrapConn(conn)
	go func() {
		if err := transfer.Socks5Transport(t, request); err != nil {
			slog.Error("connection error", "msg", err)
		}
	}()
}
