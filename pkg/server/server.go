package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type CloudServer interface {
	Verification(http.ResponseWriter) (bool, error)
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

func (s *Server) Serve() {
	// http server
	serverAddr := strings.Join(append([]string{s.host, s.port}), ":")
	// http handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		transfer := protocol2Server[s.proto]

		service := transfer(s, r)
		// do check
		if ok, err := service.Verification(w); ok {
			service.Serve(w, r)
		} else {
			log.Error(err)
		}
	})
	// http服务
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		os.Exit(0)
	}
}
