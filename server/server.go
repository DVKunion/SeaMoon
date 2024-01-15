package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/server/service"
)

type Server struct {
	srv  service.Service
	opts Options

	startAt time.Time
}

func New(opts ...Option) (*Server, error) {
	s := &Server{}
	for _, o := range opts {
		err := o(&s.opts)
		if err != nil {
			return s, err
		}
	}
	// check
	if srv, ok := service.Factory[s.opts.proto]; ok {
		s.srv = srv
	}

	if s.srv == nil {
		return s, errors.New("xxxx")
	}

	return s, nil
}

func (s *Server) Serve(ctx context.Context) error {
	// http server
	serverAddr := strings.Join(append([]string{s.opts.host, s.opts.port}), ":")

	mux := http.NewServeMux()

	lc := net.ListenConfig{}
	lc.SetMultipathTCP(true)

	ln, err := lc.Listen(ctx, "tcp", serverAddr)

	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	s.srv.Handle(mux)
	s.startAt = time.Now()
	// inject
	mux.HandleFunc("/_health", s.health)

	slog.Info("seamoon server start with", "options", s.opts)
	// http服务
	if err := server.Serve(ln); err != nil {
		slog.Error("server error", err)
		return err
	}
	return nil
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK\n" + s.startAt.Format("2006-01-02 15:04:05") + "\n" + consts.Version + "-" + consts.Commit))
	if err != nil {
		slog.Error("server status error", "msg", err)
		return
	}
}
