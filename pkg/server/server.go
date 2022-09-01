package server

import (
	"net/http"
	"strings"
)

type CloudServer interface {
	Verification(http.ResponseWriter) bool
	Serve(http.ResponseWriter, *http.Request)
}

type Server struct {
	proto string
	host  string
	port  string
}

func NewServer(proto string, host string, port string) *Server {
	return &Server{
		proto: proto,
		host:  host,
		port:  port,
	}
}

var protocol2Server = map[string]func(*Server, *http.Request) CloudServer{
	"http":   (*Server).HttpServerTransfer,
	"socks5": (*Server).SocksServerTransfer,
}

func (s *Server) Serve() error {
	// http server
	serverAddr := strings.Join(append([]string{s.host, s.port}), ":")
	// http handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		transfer := protocol2Server[s.proto]

		server := transfer(s, r)
		// do check
		if server.Verification(w) {
			// do handle
			server.Serve(w, r)
		}
	})
	// websocket upgrader

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		return err
	}
	return nil
}
