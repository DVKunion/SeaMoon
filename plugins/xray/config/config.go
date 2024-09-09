package config

import (
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"

	"github.com/DVKunion/SeaMoon/plugins/xray/net"
)

const seamoonApiTag = "seamoon-xray-api"

var (
	empty = []byte("{}")
)

type Config = conf.Config

type BoundConfig struct {
	Addr net.Address
	Port net.Port
	Tag  string
}

func Render(opts ...Options) (*core.Config, error) {
	tmpl := &conf.Config{
		InboundConfigs:  make([]conf.InboundDetourConfig, 0),
		OutboundConfigs: make([]conf.OutboundDetourConfig, 0),
	}
	for _, o := range opts {
		o(tmpl)
	}
	cf, err := tmpl.Build()
	if err != nil {
		return nil, err
	}
	return cf, nil
}

type Options func(conf *conf.Config)

func WithApiConfig() Options {
	return func(cf *conf.Config) {
		cf.API = &conf.APIConfig{
			Tag:    seamoonApiTag,
			Listen: "127.0.0.1:10086",
			Services: []string{
				"HandlerService",
				//"LoggerService",
				"StatsService",
				"RoutingService",
			},
		}
	}
}
