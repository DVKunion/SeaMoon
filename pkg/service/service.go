package service

import (
	"context"
	"net"

	"github.com/DVKunion/SeaMoon/pkg/transfer"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type Service interface {
	Conn(ctx context.Context, t transfer.Type, sOpts ...Option) (net.Conn, error)
	Serve(ln net.Listener, srvOpt ...Option) error
}

var Factory = map[tunnel.Type]Service{}

func register(t tunnel.Type, s Service) {
	Factory[t] = s
}
