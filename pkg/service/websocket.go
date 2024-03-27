package service

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
	"github.com/DVKunion/SeaMoon/pkg/transfer"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

const (
	defaultTimeout        = 10 * time.Second
	defaultReadBufferSize = 32 * 1024
)

type WSService struct {
	startAt  time.Time
	upGrader *websocket.Upgrader
}

func init() {
	register(enum.TunnelTypeWST, &WSService{})
}

func (s *WSService) Conn(ctx context.Context, t enum.ProxyType, sOpts ...Option) (tunnel.Tunnel, error) {
	// todo: useless ctx
	var srvOpts = &Options{}

	for _, o := range sOpts {
		o(srvOpts)
	}

	wsDialer := &websocket.Dialer{
		HandshakeTimeout:  defaultTimeout,
		ReadBufferSize:    defaultReadBufferSize,
		WriteBufferSize:   defaultReadBufferSize,
		EnableCompression: false,
	}

	if srvOpts.buffers != nil {
		wsDialer.ReadBufferSize = srvOpts.buffers.ReadBufferSize
		wsDialer.WriteBufferSize = srvOpts.buffers.WriteBufferSize
		wsDialer.EnableCompression = srvOpts.buffers.EnableCompression
	}

	var requestHeader = http.Header{}

	if srvOpts.tor {
		requestHeader.Add("SM-ONION", "enable")
	}

	wsConn, _, err := wsDialer.Dial(t.Path(srvOpts.addr), requestHeader)

	if err != nil {
		return nil, err
	}
	return tunnel.WsWrapConn(wsConn), nil
}

func (s *WSService) Serve(ln net.Listener, sOpts ...Option) error {
	var srvOpts = &Options{}

	for _, o := range sOpts {
		o(srvOpts)
	}

	s.upGrader = &websocket.Upgrader{
		HandshakeTimeout:  defaultTimeout,
		ReadBufferSize:    defaultReadBufferSize,
		WriteBufferSize:   defaultReadBufferSize,
		EnableCompression: false,
		CheckOrigin:       func(r *http.Request) bool { return true },
	}

	if srvOpts.buffers != nil {
		s.upGrader.ReadBufferSize = srvOpts.buffers.ReadBufferSize
		s.upGrader.WriteBufferSize = srvOpts.buffers.WriteBufferSize
		s.upGrader.EnableCompression = srvOpts.buffers.EnableCompression
	}

	if srvOpts.tlsConf != nil {
		ln = tls.NewListener(ln, srvOpts.tlsConf)
	}

	addr := strings.Split(ln.Addr().String(), ":")
	port := 443 // 默认端口
	if len(addr) > 1 {
		if p, err := strconv.Atoi(addr[1]); err == nil {
			port = p
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auto", s.auto)
	// websocket http proxy handler
	mux.HandleFunc("/http", s.http)

	// websocket socks5 proxy handler
	mux.HandleFunc("/socks5", s.socks5)

	if err := transfer.InitV2ray(uint32(port), srvOpts.uid, srvOpts.pass, srvOpts.crypt, enum.TunnelTypeWST, srvOpts.tor, srvOpts.tlsConf != nil); err == nil {
		mux.HandleFunc("/vmess", s.v2ray("vmess"))
		mux.HandleFunc("/vless", s.v2ray("vless"))
		mux.HandleFunc("/v-shadowsocks", s.v2ray("shadowsocks"))
	} else {
		xlog.Error(xlog.ServiceV2rayInitError, "err", err)
	}

	s.startAt = time.Now()
	// inject
	mux.HandleFunc("/_health", s.health)
	server := &http.Server{
		Addr:    srvOpts.addr,
		Handler: mux,
	}
	return server.Serve(ln)
}

func (s *WSService) auto(w http.ResponseWriter, r *http.Request) {
	// means use http to connector
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	t := tunnel.WsWrapConn(conn)
	go func() {
		if err := transfer.AutoTransport(t); err != nil {
			xlog.Error(xlog.ServiceTransportError, "type", "auto", "err", err)
		}
	}()
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
			xlog.Error(xlog.ServiceTransportError, "type", "http", "err", err)
		}
	}()
}

func (s *WSService) socks5(w http.ResponseWriter, r *http.Request) {
	onion := r.Header.Get("SM-ONION")
	// means use socks5 to connector
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	t := tunnel.WsWrapConn(conn)
	go func() {
		// 检测是否存在 onion 标识，代表着是否要开启 tor 转发
		if onion != "" {
			if err := transfer.TorTransport(t); err != nil {
				xlog.Error(xlog.ServiceTransportError, "type", "socks5+tor", "err", err)
			}
		} else {
			if err := transfer.Socks5Transport(t, false); err != nil {
				xlog.Error(xlog.ServiceTransportError, "type", "socks5", "err", err)
			}
		}
	}()
}

func (s *WSService) v2ray(proto string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upGrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		t := tunnel.WsWrapConn(conn)
		go func() {
			if err := transfer.V2rayTransport(t, proto); err != nil {
				xlog.Error(xlog.ServiceTransportError, "type", "v2ray", "proto", proto, "err", err)
			}
		}()
	}
}

func (s *WSService) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK\n" + s.startAt.Format("2006-01-02 15:04:05") + "\n" + xlog.Version + "-" + xlog.Commit))
	if err != nil {
		xlog.Error(xlog.ServiceStatusError, "err", err)
		return
	}
}
