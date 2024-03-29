package transfer

import (
	"context"
	"encoding/json"
	"fmt"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	pinboud "github.com/v2fly/v2ray-core/v5/app/proxyman/inbound"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/features/inbound"
	"github.com/v2fly/v2ray-core/v5/features/routing"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon"
	"github.com/v2fly/v2ray-core/v5/infra/conf/synthetic/log"
	v4 "github.com/v2fly/v2ray-core/v5/infra/conf/v4"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

const handleTag = "seamoon-"

var v2ray *core.Instance

func InitV2rayServer(port uint32, id, pass, crypt string, tp enum.TunnelType, tor bool, tls bool) error {
	config, err := renderConfig(port, id, pass, crypt, tp, tor, tls)
	if err != nil {
		return err
	}
	v2ray, err = core.New(config)
	if err != nil {
		return err
	}
	return nil
}

// V2rayTransport v2ray 相关协议支持: vmess / vless
// 这是一个偷懒的版本，并没有详细的研究对应协议的具体通信解析方案, 直接集成了 v2ray-core, 并且实现的相当的简陋。
// 还是期望能够和 socks5 一样保持一致是最好的
// proto 来自己做 dispatch
func V2rayTransport(conn net.Conn, proto string) error {
	ctx, _ := context.WithCancel(context.Background())
	manager := v2ray.GetFeature(inbound.ManagerType()).(inbound.Manager)
	handler, err := manager.GetHandler(ctx, handleTag+proto)
	if err != nil {
		return err
	}
	worker := handler.(*pinboud.AlwaysOnInboundHandler).GetInbound()
	if err != nil {
		return err
	}

	sid := session.NewID()
	ctx = session.ContextWithID(ctx, sid)
	ctx = session.ContextWithInbound(ctx, &session.Inbound{
		Source:  net.DestinationFromAddr(conn.RemoteAddr()),
		Gateway: net.TCPDestination(net.ParseAddress("0.0.0.0"), net.PortFromBytes([]byte("8900"))),
		Tag:     handleTag + proto,
	})

	content := new(session.Content)
	ctx = session.ContextWithContent(ctx, content)

	dispatch := v2ray.GetFeature(routing.DispatcherType()).(*dispatcher.DefaultDispatcher)
	err = worker.Process(ctx, net.Network_TCP, conn, dispatch)
	return err
}

func renderSetting(proto string, crypt string, pass string) *[]byte {
	str := ""
	switch proto {
	case "vmess":
		str = fmt.Sprintf(`{
	"clients": [
	  {
		"id": "%s",
		"alterId": 0
	  }
	],
	"decryption":"auto"
}`, pass)
	case "vless":
		str = fmt.Sprintf(`{
	"clients": [
	  {
		"id": "%s",
		"alterId": 0
	  }
	],
	"decryption":"none"
}`, pass)
	case "shadowsocks":
		str = fmt.Sprintf(`{
	"method": "%s",
	"password": "%s",
	"network": "tcp"
}`, crypt, pass)
	}
	s := []byte(str)
	return &s
}

func renderConfig(port uint32, id string, pass string, crypt string, tp enum.TunnelType, tor bool, tls bool) (*core.Config, error) {
	t := v4.Config{
		LogConfig: &log.LogConfig{
			AccessLog: "v2ray_access.log",
			ErrorLog:  "v2ray_error.log",
			LogLevel:  "ERROR",
		},
		InboundConfigs: make([]v4.InboundDetourConfig, 0),
		OutboundConfigs: []v4.OutboundDetourConfig{
			outboundConfig(tor),
		},
	}
	if id != "" {
		t.InboundConfigs = append(t.InboundConfigs, v4.InboundDetourConfig{
			Protocol: "vmess",
			// 暂时不支持动态端口
			PortRange: &cfgcommon.PortRange{
				From: port,
				To:   port,
			},
			Settings: (*json.RawMessage)(renderSetting("vmess", crypt, id)),
			Tag:      handleTag + "vmess",
			StreamSetting: &v4.StreamConfig{
				Network: (*v4.TransportProtocol)(tp.ToPtr()),
				Security: func() string {
					if tls {
						return "tls"
					}
					return ""
				}(),
				WSSettings: &v4.WebSocketConfig{
					Path: "/vmess",
				},
			},
		})
		t.InboundConfigs = append(t.InboundConfigs, v4.InboundDetourConfig{
			Protocol: "vless",
			// 暂时不支持动态端口
			PortRange: &cfgcommon.PortRange{
				From: port,
				To:   port,
			},
			Settings: (*json.RawMessage)(renderSetting("vless", crypt, id)),
			Tag:      handleTag + "vless",
			StreamSetting: &v4.StreamConfig{
				Network:  (*v4.TransportProtocol)(tp.ToPtr()),
				Security: "tls",
				WSSettings: &v4.WebSocketConfig{
					Path: "/vless",
				},
			},
		})
	}

	if pass != "" && crypt != "" {
		t.InboundConfigs = append(t.InboundConfigs, v4.InboundDetourConfig{
			Protocol: "shadowsocks",
			PortRange: &cfgcommon.PortRange{
				From: port,
				To:   port,
			},
			Settings: (*json.RawMessage)(renderSetting("shadowsocks", crypt, pass)),
			Tag:      handleTag + "shadowsocks",
			StreamSetting: &v4.StreamConfig{
				Network:  (*v4.TransportProtocol)(tp.ToPtr()),
				Security: "tls",
				WSSettings: &v4.WebSocketConfig{
					Path: "/v-shadowsocks",
				},
			},
		})
	}
	return t.Build()
}

func outboundConfig(tor bool) v4.OutboundDetourConfig {
	torSetting := []byte(`{
                "servers": [
                    {
                        "address": "127.0.0.1",
                        "port": 9050
                    }
                ]
            }`)
	empty := []byte("{}")
	if tor {
		return v4.OutboundDetourConfig{
			Protocol: "socks",
			Settings: (*json.RawMessage)(&torSetting),
		}
	}
	return v4.OutboundDetourConfig{
		Protocol: "freedom",
		Settings: (*json.RawMessage)(&empty),
	}
}
