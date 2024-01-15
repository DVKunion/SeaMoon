package service

import (
	"net/http"

	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type Service interface {
	Handle(m *http.ServeMux)
}

var Factory = map[tunnel.Type]Service{}

func register(t tunnel.Type, s Service) {
	Factory[t] = s
}
