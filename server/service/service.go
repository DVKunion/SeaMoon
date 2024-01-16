package service

import (
	"net"

	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type Service interface {
	Serve(ln net.Listener, srvOpt ...Option) error
}

var Factory = map[tunnel.Type]Service{}

func register(t tunnel.Type, s Service) {
	Factory[t] = s
}
