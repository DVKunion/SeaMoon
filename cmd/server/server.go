package server

import (
	"context"
	"os"
	"strings"

	net "github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/service"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

type Server struct {
	srv  service.Service
	opts Options
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
	if srv, ok := service.Factory.Load(s.opts.proto); ok {
		s.srv = srv.(service.Service)
	}

	if s.srv == nil {
		return s, errors.New("xxxx")
	}

	return s, nil
}

// Serve do common serve
func (s *Server) Serve(ctx context.Context) error {
	network := "tcp"

	if net.IsIPv4(s.opts.host) {
		network = "tcp4"
	}

	serverAddr := strings.Join(append([]string{s.opts.host, s.opts.port}), ":")

	lc := net.ListenConfig{}
	if s.opts.mtcp {
		lc.SetMultipathTCP(true)
	}

	ln, err := lc.Listen(ctx, network, serverAddr)

	if err != nil {
		return err
	}

	xlog.Info(xlog.ApiServerStart, "options", s.opts)

	var srvOpt []service.Option

	srvOpt = append(srvOpt, service.WithAddr(serverAddr))
	srvOpt = append(srvOpt, service.WithUid(os.Getenv("SM_UID")))
	srvOpt = append(srvOpt, service.WithPassword(os.Getenv("SM_SS_PASS")))
	srvOpt = append(srvOpt, service.WithCrypt(os.Getenv("SM_SS_CRYPT")))

	if err := s.srv.Serve(ln, srvOpt...); err != nil {
		xlog.Error(xlog.ServiceError, err)
		return err
	}
	return nil
}
