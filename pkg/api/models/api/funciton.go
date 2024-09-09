package api

import (
	"fmt"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models/external"
)

type TunnelApi struct {
	ProviderId   uint               `json:"provider_id"`
	ProviderType *enum.ProviderType `json:"provider_type"`

	ID        uint      `json:"id"`
	UniqID    *string   `json:"uniq_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name          *string          `json:"name"`
	Addr          *string          `json:"address"`
	Port          *int32           `json:"port"`
	Type          *enum.TunnelType `json:"type"`
	StatusMessage *string          `json:"status_message"`

	Config *TunnelConfig `json:"tunnel_config"`
}

type TunnelCreateApi struct {
	ID            uint               `json:"id"`
	ProviderId    uint               `json:"provider_id"`
	UniqID        *string            `json:"uniq_id"`
	Name          *string            `json:"name"`
	Port          *int32             `json:"port"`
	Type          *enum.TunnelType   `json:"type"`
	Status        *enum.TunnelStatus `json:"status"`
	StatusMessage *string            `json:"status_message"`
	Addr          *string            `json:"address"`
	Config        *TunnelConfig      `json:"tunnel_config"`
}

type TunnelCreateApiList []*TunnelCreateApi

func (t Tunnel) GetAddr() string {
	switch *t.Type {
	case enum.TunnelTypeWST:
		if t.Config.TLS {
			return "wss://" + *t.Addr
		}
		return "ws://" + *t.Addr
	case enum.TunnelTypeGRT:
		if t.Config.TLS {
			return "grpcs://" + *t.Addr
		}
		return "grpc://" + *t.Addr
	}
	return ""
}

func (t Tunnel) ToApi(extra ...func(api interface{})) *TunnelApi {
	return toApi(t, &TunnelApi{}, extra...).(*TunnelApi)
}

func (tl TunnelList) ToApi(extra ...func(api interface{})) []*TunnelApi {
	res := make([]*TunnelApi, 0)
	for _, t := range tl {
		api := toApi(t, &TunnelApi{}, extra...)
		res = append(res, api.(*TunnelApi))
	}
	return res
}

func (tl TunnelList) ToConfig(p string) []byte {
	switch p {
	case "clash":
		cc := external.ClashConfig{
			MixedPort:          7890,
			AllowLan:           false,
			LogLevel:           "info",
			ExternalController: "127.0.0.1:9090",
			Secret:             "",
			//DNS: ClashDNS{
			//	Enable:       true,
			//	Ipv6:         false,
			//	Listen:       "127.0.0.1:5353",
			//	EnhancedMode: "fake-ip",
			//	FakeIPFilter: []string{"*.lan"},
			//	Nameserver:
			//},
			Proxies: make([]external.ClashProxies, 0),
			ProxyGroups: []external.ClashProxyGroups{
				{
					Name:    "Proxies",
					Type:    "select",
					Proxies: make([]string, 0),
				},
				{
					Name:    "Direct",
					Type:    "select",
					Proxies: []string{"DIRECT"},
				},
			},
			Rules: external.BindingRules,
		}
		for _, t := range tl {
			cc.Proxies = append(cc.Proxies, external.ClashProxies{
				Name:   *t.Name + "-" + t.Config.Region + "-vmess",
				Type:   "vmess",
				Server: *t.Addr,
				Port: func() int {
					if t.Config.TLS {
						return 443
					}
					return 80
				}(),
				UUID:           t.Config.V2rayUid,
				NetWork:        "ws",
				TLS:            t.Config.TLS,
				SkipCertVerify: !t.Config.TLS,
				Cipher:         "auto",
				AlterId:        0,
				WsOpts: struct {
					Path string `yaml:"path,omitempty"`
				}(struct{ Path string }{Path: "/vmess"}),
			})
			//cc.Proxies = append(cc.Proxies, ClashProxies{
			//	Name:   *t.Name + "-" + t.Config.Region + "-vless",
			//	Type:   "vless",
			//	Server: *t.Addr,
			//	Port: func() int {
			//		if t.Config.TLS {
			//			return 443
			//		}
			//		return 80
			//	}(),
			//	UUID:           t.Config.V2rayUid,
			//	UDP:            false,
			//	NetWork:        "ws",
			//	TLS:            t.Config.TLS,
			//	SkipCertVerify: !t.Config.TLS,
			//	Cipher:         "auto",
			//	AlterId:        0,
			//	WsOpts: struct {
			//		Path string `yaml:"path,omitempty"`
			//	}(struct{ Path string }{Path: "/vless"}),
			//})
			cc.ProxyGroups[0].Proxies = append(cc.ProxyGroups[0].Proxies, *t.Name+"-"+t.Config.Region+"-vmess")
			//cc.ProxyGroups[0].Proxies = append(cc.ProxyGroups[0].Proxies, *t.Name+"-"+t.Config.Region+"-vless")
		}
		data, err := yaml.Marshal(&cc)
		if err != nil {

		}
		return data
	case "shadowrocket":
		res := ""
		for _, t := range tl {
			res += fmt.Sprintf("vmess://%s?remarks=%s&path=/vmess&obfs=%s&tls=%d&alterId=0\n",
				base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("auto:%s@%s:%v", t.Config.V2rayUid, *t.Addr, func() int {
					if t.Config.TLS {
						return 443
					}
					return 80
				}()))),
				fmt.Sprintf("%s-%s-%s", *t.Name, t.Config.Region, "vmess"),
				t.Type.String(),
				func() int {
					if t.Config.TLS {
						return 1
					}
					return 0
				}(),
			)
			res += fmt.Sprintf("vless://%s?remarks=%s&path=/vless&obfs=%s&tls=%d&alterId=0\n",
				base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("auto:%s@%s:%v", t.Config.V2rayUid, *t.Addr, func() int {
					if t.Config.TLS {
						return 443
					}
					return 80
				}()))),
				fmt.Sprintf("%s-%s-%s", *t.Name, t.Config.Region, "vless"),
				t.Type.String(),
				func() int {
					if t.Config.TLS {
						return 1
					}
					return 0
				}(),
			)
		}
		return []byte(base64.URLEncoding.EncodeToString([]byte(res)))
	}
	return nil
}

func (ta TunnelCreateApi) ToModel(full bool) *Tunnel {
	return toModel(ta, &Tunnel{}, full).(*Tunnel)
}

func NewTunnelCreateApi() *TunnelCreateApi {
	res := &TunnelCreateApi{}
	autoFull(res)
	return res
}
