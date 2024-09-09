package config

import (
	"fmt"

	"github.com/xtls/xray-core/infra/conf"

	"github.com/DVKunion/SeaMoon/pkg/api/models/abstract"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
	"github.com/DVKunion/SeaMoon/plugins/xray/net"
)

var inboundProtoMaps = map[string]string{
	// http 默认配置了超时时间为300，0为不限制
	"http": `{"timeout": 300, "accounts": %s, "allowTransparent": false, "userLevel": 0}`,
	// socks5 暂时不支持 udp
	"socks": `{"auth": "%s", "accounts": %s, "udp": false, "ip": "127.0.0.1", "userLevel": 0}`,
	// shadowsocks 暂时也不开 udp 和 ota, xray 实际上支持了多租户模式的socks, 本期暂时不开放
	"shadowsocks": `{"email": "%s", "method": "%s", "password": "%s", "level": 0, "ota": false, "network": "tcp"}`,
	// torjan 暂时不支持 xray fallbacks 特性
	"torjan": `{"clients": %s, "fallbacks": []}`,
	//"vmess": `{"clients": %s, "default": { "level": 0}, "detour": {"to": "tag_to_detour"}}`,
	"vmess": `{"clients": %s, "default": {"level": 0}}`,
	// vless 暂时不支持 xray fallbacks 特性, decryption 目前一直只能为空
	"vless": `{"clients": %s, "decryption": "none", "fallbacks": []}`,
}

// WithHttpInbound http 入站配置
func WithHttpInbound(bc *BoundConfig, auths abstract.AuthList) Options {
	return withInbound("http", bc,
		[]byte(fmt.Sprintf(inboundProtoMaps["http"], auths.Marshal())))
}

// WithSocksInbound socks5 入站配置
func WithSocksInbound(bc *BoundConfig, auths abstract.AuthList) Options {
	if auths.IsEmpty() {
		return withInbound("socks", bc,
			[]byte(fmt.Sprintf(inboundProtoMaps["socks"], "noauth", "")))
	}
	return withInbound("socks", bc,
		[]byte(fmt.Sprintf(inboundProtoMaps["socks"], "password", auths.Marshal())))
}

func WithShadowSocksInbound(bc *BoundConfig, au *abstract.EmailPassAuth) Options {
	return withInbound("shadowsocks", bc,
		[]byte(fmt.Sprintf(inboundProtoMaps["shadowsocks"], au.Email, au.Method, au.Password)))
}

// WithTorjanInbound torjan 入站配置
func WithTorjanInbound(bc *BoundConfig, auths abstract.AuthList) Options {
	return withInbound("torjan", bc,
		[]byte(fmt.Sprintf(inboundProtoMaps["torjan"], auths.Marshal())))
}

func WithVmessInbound(bc *BoundConfig, auths abstract.AuthList) Options {
	return withInbound("vmess", bc,
		[]byte(fmt.Sprintf(inboundProtoMaps["vmess"], auths.Marshal())))
}

func WithVlessInbound(bc *BoundConfig, auths abstract.AuthList) Options {
	return withInbound("vless", bc,
		[]byte(fmt.Sprintf(inboundProtoMaps["vless"], auths.Marshal())))
}

func withInbound(proto string, bc *BoundConfig, settings []byte) Options {
	return func(cf *conf.Config) {
		cf.InboundConfigs = append(cf.InboundConfigs, conf.InboundDetourConfig{
			Protocol: proto,
			Tag:      bc.Tag,
			ListenOn: net.A2Cfg(bc.Addr),
			PortList: net.P2Cfg(bc.Port),
			Settings: tools.JsonRawPtr(settings),
		})
	}
}
