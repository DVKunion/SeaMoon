package service

import (
	"context"
	"net"
	"sync"

	"github.com/DVKunion/SeaMoon/pkg/transfer"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type Service interface {
	Conn(ctx context.Context, t transfer.Type, sOpts ...Option) (net.Conn, error)
	Serve(ln net.Listener, srvOpt ...Option) error
}

var Factory = sync.Map{}

func register(t tunnel.Type, s Service) {
	Factory.Store(t, s)
}
