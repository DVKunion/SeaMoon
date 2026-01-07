package transfer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon"
	"github.com/v2fly/v2ray-core/v5/infra/conf/synthetic/log"
	v4 "github.com/v2fly/v2ray-core/v5/infra/conf/v4"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

const handleTag = "seamoon-"

var (
	v2ray *core.Instance
	empty = []byte("{}")

	// 级联代理全局配置
	cascadeConfig *CascadeProxyConfig
)

// CascadeProxyConfig 级联代理配置
type CascadeProxyConfig struct {
	Enabled  bool
	Addr     string
	Uid      string
	Password string
}

// SetCascadeConfig 设置级联代理配置
func SetCascadeConfig(addr, uid, password string) {
	if addr != "" && uid != "" {
		cascadeConfig = &CascadeProxyConfig{
			Enabled:  true,
			Addr:     addr,
			Uid:      uid,
			Password: password,
		}
	}
}

// GetCascadeConfig 获取级联代理配置
func GetCascadeConfig() *CascadeProxyConfig {
	return cascadeConfig
}

// IsCascadeEnabled 检查是否启用级联代理
func IsCascadeEnabled() bool {
	return cascadeConfig != nil && cascadeConfig.Enabled
}

type v2rayConfig struct {
	mode string
	addr string // 用于出站的协议地址
	port uint32 //

	id    string
	pass  string
	crypt string

	proto string
	tp    enum.TunnelType

	tor bool
	tls bool

	// 级联代理配置
	cascadeAddr     string
	cascadeUid      string
	cascadePassword string
}

type ConfigOpt func(config *v2rayConfig)

func WithServerMod() ConfigOpt {
	return func(config *v2rayConfig) {
		config.mode = "server"
	}
}

func WithClientMod() ConfigOpt {
	return func(config *v2rayConfig) {
		config.mode = "client"
	}
}

func WithNetAddr(addr string, port uint32) ConfigOpt {
	return func(config *v2rayConfig) {
		config.addr = addr
		config.port = port
	}
}

func WithAuthInfo(id, crypt, pass string) ConfigOpt {
	return func(config *v2rayConfig) {
		config.id = id
		config.crypt = crypt
		config.pass = pass
	}
}

func WithExtra(tor, tls bool) ConfigOpt {
	return func(config *v2rayConfig) {
		config.tor = tor
		config.tls = tls
	}
}

func WithCascadeProxy(addr, uid, password string) ConfigOpt {
	return func(config *v2rayConfig) {
		config.cascadeAddr = addr
		config.cascadeUid = uid
		config.cascadePassword = password
	}
}

func WithTunnelType(proto string, tp enum.TunnelType) ConfigOpt {
	return func(config *v2rayConfig) {
		config.proto = proto
		config.tp = tp
	}
}

func NewV2rayConfig(opts ...ConfigOpt) *v2rayConfig {
	v := &v2rayConfig{}
	for _, o := range opts {
		o(v)
	}
	return v
}

func (v v2rayConfig) Build() (*core.Config, error) {
	t := &v4.Config{
		LogConfig: &log.LogConfig{
			AccessLog: "v2ray_access.log",
			ErrorLog:  "v2ray_error.log",
			LogLevel:  "ERROR",
		},
		InboundConfigs:  v.InboundConfig(),
		OutboundConfigs: v.OutboundConfig(),
	}

	return t.Build()
}

func (v v2rayConfig) InboundConfig() []v4.InboundDetourConfig {
	cs := make([]v4.InboundDetourConfig, 0)
	switch v.mode {
	case "server":
		if v.id != "" {
			cs = append(cs, v.vlessInboundConfig())
			cs = append(cs, v.vmessInboundConfig())
		}
		if v.crypt != "" && v.pass != "" {
			cs = append(cs, v.shadowsocksInboundConfig())
		}
	case "client":
		cs = append(cs, v.httpInboundConfig())
		cs = append(cs, v.socks5InboundConfig())
	}
	return cs
}

func (v v2rayConfig) OutboundConfig() []v4.OutboundDetourConfig {
	cs := make([]v4.OutboundDetourConfig, 0)
	switch v.mode {
	case "server":
		// 如果配置了级联代理，使用 vless 协议转发到下一跳
		if v.cascadeAddr != "" && v.cascadeUid != "" {
			// 解析级联代理地址
			addr, port, cascadeTLS := v.parseCascadeAddr()
			// 创建级联代理的 vless outbound 配置
			cs = append(cs, v.cascadeVlessOutboundConfig(addr, port, cascadeTLS))
		} else if v.tor {
			cs = append(cs, v.torOutboundConfig())
		} else {
			cs = append(cs, v.freedomOutboundConfig())
		}
	case "client":
		switch v.proto {
		case "vmess":
			cs = append(cs, v.vmessOutboundConfig())
		case "vless":
			cs = append(cs, v.vlessOutboundConfig())
		case "shadowsocks":
			cs = append(cs, v.shadowsocksOutboundConfig())
		}
	}
	return cs
}

func (v v2rayConfig) StreamSetting(proto string) *v4.StreamConfig {
	switch v.tp {
	case enum.TunnelTypeWST:
		return v.streamWebsocketSetting(proto)
	case enum.TunnelTypeGRT:
		return v.streamGrpcSetting(proto)
	}
	return nil
}

func (v v2rayConfig) streamGrpcSetting(proto string) *v4.StreamConfig {
	return nil
}

func (v v2rayConfig) streamWebsocketSetting(proto string) *v4.StreamConfig {
	return &v4.StreamConfig{
		Network: (*v4.TransportProtocol)(v.tp.ToPtr()),
		Security: func() string {
			if v.tls {
				return "tls"
			}
			return ""
		}(),
		WSSettings: &v4.WebSocketConfig{
			Path: "/" + proto,
		},
	}
}

func (v v2rayConfig) httpInboundConfig() v4.InboundDetourConfig {
	vc := []byte(`{
        "accounts": [],
        "allowTransparent": false
      }`)
	return v4.InboundDetourConfig{
		Protocol: "http",
		PortRange: &cfgcommon.PortRange{
			From: v.port,
			To:   v.port,
		},
		Settings: (*json.RawMessage)(&vc),
		Tag:      handleTag + "http",
	}
}

func (v v2rayConfig) socks5InboundConfig() v4.InboundDetourConfig {
	vc := []byte(`{
        "auth": "noauth",
        "accounts": [],
        "udp": false
      }`)
	return v4.InboundDetourConfig{
		Protocol: "socks",
		PortRange: &cfgcommon.PortRange{
			From: v.port,
			To:   v.port,
		},
		Settings: (*json.RawMessage)(&vc),
		Tag:      handleTag + "socks5",
	}
}

func (v v2rayConfig) shadowsocksOutboundConfig() v4.OutboundDetourConfig {
	vc := []byte(fmt.Sprintf(`{
    "servers": [
        {
            "address": "%s",
            "port": %d,
			"method": "%s"
            "password": "%s"
        }
    ]
}`, v.addr, v.port, v.crypt, v.pass))
	return v4.OutboundDetourConfig{
		Protocol:      "shadowsocks",
		Settings:      (*json.RawMessage)(&vc),
		Tag:           handleTag + "shadowsocks",
		StreamSetting: v.StreamSetting("v-shadowsocks"),
	}
}

func (v v2rayConfig) shadowsocksInboundConfig() v4.InboundDetourConfig {
	vc := []byte(fmt.Sprintf(`{
	"method": "%s",
	"password": "%s",
	"network": "tcp"
}`, v.crypt, v.pass))
	return v4.InboundDetourConfig{
		Protocol: "shadowsocks",
		PortRange: &cfgcommon.PortRange{
			From: v.port,
			To:   v.port,
		},
		Settings:      (*json.RawMessage)(&vc),
		Tag:           handleTag + "shadowsocks",
		StreamSetting: v.StreamSetting("v-shadowsocks"),
	}
}

func (v v2rayConfig) vmessOutboundConfig() v4.OutboundDetourConfig {
	outSetting := []byte(fmt.Sprintf(`{
    "vnext": [
        {
            "address": "%s",
            "port": %d,
            "users": [
                {
                    "alterId": 0,
                    "id": "%s",
                    "security": "auto"
                }
            ]
        }
    ]
}`, v.addr, v.port, v.id))
	return v4.OutboundDetourConfig{
		Protocol:      "vmess",
		Tag:           handleTag + "vmess",
		Settings:      (*json.RawMessage)(&outSetting),
		StreamSetting: v.StreamSetting("vmess"),
	}
}

func (v v2rayConfig) vmessInboundConfig() v4.InboundDetourConfig {
	vc := []byte(fmt.Sprintf(`{
	"clients": [
	  {
		"id": "%s",
		"alterId": 0
	  }
	],
	"decryption":"auto"
}`, v.id))
	return v4.InboundDetourConfig{
		Protocol: "vmess",
		PortRange: &cfgcommon.PortRange{
			From: v.port,
			To:   v.port,
		},
		Settings:      (*json.RawMessage)(&vc),
		Tag:           handleTag + "vmess",
		StreamSetting: v.StreamSetting("vmess"),
	}
}

func (v v2rayConfig) vlessOutboundConfig() v4.OutboundDetourConfig {
	outSetting := []byte(fmt.Sprintf(`{
    "vnext": [
        {
            "address": "%s",
            "port": %d,
            "users": [
                {
                    "alterId": 0,
                    "id": "%s",
                    "security": "auto",
					"encryption": "none"
                }
            ]
        }
    ]
}`, v.addr, v.port, v.id))
	return v4.OutboundDetourConfig{
		Protocol:      "vless",
		Settings:      (*json.RawMessage)(&outSetting),
		Tag:           handleTag + "vless",
		StreamSetting: v.StreamSetting("vless"),
	}
}

func (v v2rayConfig) vlessInboundConfig() v4.InboundDetourConfig {
	vc := []byte(fmt.Sprintf(`{
	"clients": [
	  {
		"id": "%s",
		"alterId": 0
	  }
	],
	"decryption":"none"
}`, v.id))
	return v4.InboundDetourConfig{
		Protocol: "vless",
		PortRange: &cfgcommon.PortRange{
			From: v.port,
			To:   v.port,
		},
		Settings:      (*json.RawMessage)(&vc),
		Tag:           handleTag + "vless",
		StreamSetting: v.StreamSetting("vless"),
	}
}

func (v v2rayConfig) torOutboundConfig() v4.OutboundDetourConfig {
	torSetting := []byte(`{
                "servers": [
                    {
                        "address": "127.0.0.1",
                        "port": 9050
                    }
                ]
            }`)
	return v4.OutboundDetourConfig{
		Protocol: "socks",
		Settings: (*json.RawMessage)(&torSetting),
	}
}

func (v v2rayConfig) freedomOutboundConfig() v4.OutboundDetourConfig {
	return v4.OutboundDetourConfig{
		Protocol: "freedom",
		Settings: (*json.RawMessage)(&empty),
	}
}

// parseCascadeAddr 解析级联代理地址，返回地址、端口和是否使用 TLS
func (v v2rayConfig) parseCascadeAddr() (string, uint32, bool) {
	// cascadeAddr 格式可能是: "example.com" 或 "example.com:443" 或 "wss://example.com" 或 "ws://example.com"
	addr := v.cascadeAddr
	port := uint32(443) // 默认端口
	useTLS := false

	// 移除协议前缀并判断 TLS
	if strings.HasPrefix(addr, "wss://") {
		addr = strings.TrimPrefix(addr, "wss://")
		port = 443
		useTLS = true
	} else if strings.HasPrefix(addr, "ws://") {
		addr = strings.TrimPrefix(addr, "ws://")
		port = 80
		useTLS = false
	} else if strings.HasPrefix(addr, "https://") {
		addr = strings.TrimPrefix(addr, "https://")
		port = 443
		useTLS = true
	} else if strings.HasPrefix(addr, "http://") {
		addr = strings.TrimPrefix(addr, "http://")
		port = 80
		useTLS = false
	} else if strings.HasPrefix(addr, "grpcs://") {
		addr = strings.TrimPrefix(addr, "grpcs://")
		port = 443
		useTLS = true
	} else if strings.HasPrefix(addr, "grpc://") {
		addr = strings.TrimPrefix(addr, "grpc://")
		port = 80
		useTLS = false
	}

	// 检查是否包含端口
	if strings.Contains(addr, ":") {
		parts := strings.Split(addr, ":")
		if len(parts) == 2 {
			addr = parts[0]
			if p, err := strconv.ParseUint(parts[1], 10, 32); err == nil {
				port = uint32(p)
			}
		}
	}

	return addr, port, useTLS
}

// cascadeVlessOutboundConfig 创建级联代理的 vless outbound 配置
func (v v2rayConfig) cascadeVlessOutboundConfig(addr string, port uint32, cascadeTLS bool) v4.OutboundDetourConfig {
	outSetting := []byte(fmt.Sprintf(`{
    "vnext": [
        {
            "address": "%s",
            "port": %d,
            "users": [
                {
                    "alterId": 0,
                    "id": "%s",
                    "security": "auto",
					"encryption": "none"
                }
            ]
        }
    ]
}`, addr, port, v.cascadeUid))
	
	// 为级联代理创建 StreamSetting，使用 websocket 传输（级联代理通常使用 websocket）
	tunnelType := enum.TunnelType(enum.TunnelTypeWST)
	streamSetting := &v4.StreamConfig{
		Network: (*v4.TransportProtocol)(tunnelType.ToPtr()),
		Security: func() string {
			if cascadeTLS {
				return "tls"
			}
			return ""
		}(),
		WSSettings: &v4.WebSocketConfig{
			Path: "/vless",
		},
	}
	
	return v4.OutboundDetourConfig{
		Protocol:      "vless",
		Settings:      (*json.RawMessage)(&outSetting),
		Tag:           handleTag + "cascade-vless",
		StreamSetting: streamSetting,
	}
}
