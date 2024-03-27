package service

import (
	"context"
	"net"
	"sync"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type Service interface {
	Conn(ctx context.Context, t enum.ProxyType, sOpts ...Option) (tunnel.Tunnel, error)
	Serve(ln net.Listener, srvOpt ...Option) error
}

var Factory = sync.Map{}

func register(t enum.TunnelType, s Service) {
	Factory.Store(t, s)
}
