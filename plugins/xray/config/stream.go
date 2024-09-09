package config

import (
	"github.com/xtls/xray-core/infra/conf"

	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

// WithWSSettings websocket stream 配置
func WithWSSettings(tls bool, path string) *conf.StreamConfig {
	return &conf.StreamConfig{
		Network: (*conf.TransportProtocol)(tools.StringPtr("websocket")),
		Security: func() string {
			if tls {
				return "tls"
			}
			return ""
		}(),
		WSSettings: &conf.WebSocketConfig{
			Path: "/" + path,
		},
	}
}

// WithGrpcSettings grpc stream 配置
func WithGrpcSettings(name string) *conf.StreamConfig {
	return &conf.StreamConfig{
		Network:  (*conf.TransportProtocol)(tools.StringPtr("grpc")),
		Security: "tls",
		GRPCConfig: &conf.GRPCConfig{
			ServiceName: name,
		},
	}
}
