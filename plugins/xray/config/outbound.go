package config

import (
	"fmt"

	"github.com/xtls/xray-core/infra/conf"

	"github.com/DVKunion/SeaMoon/pkg/api/models/abstract"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

var outboundProtoMap = map[string]string{
	//"freedom": `{ "domainStrategy": "AsIs", "redirect": "127.0.0.1:3366", "userLevel": 0}`,
	"http":        `{"servers": [{"address": "%s", "port": %d, "users": %s}]}`,
	"socks":       `{"servers": [{"address": "%s", "port": %d, "users": %s}]}`,
	"shadowsocks": `{"servers": [{"address": "%s", "port": %d, "method": %s, "password": %s, "email": "%s", "level": 0}]}`,
	"torjan":      `{"servers": [{"address": "%s", "port": %d, "password": "%s", "email": "%s", "level": 0}]}`,
	"vmess":       `{"vnext": [{"address": "%s", "port": %d, "users": %s}]}`,
	"vless":       `{"vnext": [{"address": "%s", "port": %d, "users": %s}]}`,
}

// WithFreedomOutbound freedom 出站配置
func WithFreedomOutbound() Options {
	return withOutbound("freedom", "", nil, empty)
}

// WithHttpOutbound http 出站配置
func WithHttpOutbound(bc *BoundConfig, auths abstract.AuthList, sc *conf.StreamConfig) Options {
	return withOutbound("http", bc.Tag, sc,
		[]byte(fmt.Sprintf(outboundProtoMap["http"], bc.Addr.String(), bc.Port, auths.Marshal())))
}

// WithSocksOutbound socks 出站配置
func WithSocksOutbound(bc *BoundConfig, auths abstract.AuthList, sc *conf.StreamConfig) Options {
	return withOutbound("socks", bc.Tag, sc,
		[]byte(fmt.Sprintf(outboundProtoMap["socks"], bc.Addr.String(), bc.Port, auths.Marshal())))
}

// WithShadowSocksOutbound shadowsocks 出站配置
func WithShadowSocksOutbound(bc *BoundConfig, au abstract.EmailPassAuth, sc *conf.StreamConfig) Options {
	return withOutbound("shadowsocks", bc.Tag, sc,
		[]byte(fmt.Sprintf(outboundProtoMap["shadowsocks"], bc.Addr.String(), bc.Port, au.Method, au.Password, au.Email)))
}

func WithTorjanOutbound(bc *BoundConfig, au *abstract.EmailPassAuth, sc *conf.StreamConfig) Options {
	return withOutbound("torjan", bc.Tag, sc,
		[]byte(fmt.Sprintf(outboundProtoMap["torjan"], bc.Addr.String(), bc.Port, au.Password, au.Email)))
}

// WithVmessOutbound vmess 出站配置
func WithVmessOutbound(bc *BoundConfig, auths abstract.AuthList, sc *conf.StreamConfig) Options {
	return withOutbound("vmess", bc.Tag, sc,
		[]byte(fmt.Sprintf(outboundProtoMap["vmess"], bc.Addr.String(), bc.Port, auths.Marshal())))
}

// WithVlessOutbound vless 出站配置
func WithVlessOutbound(bc *BoundConfig, auths abstract.AuthList, sc *conf.StreamConfig) Options {
	return withOutbound("vless", bc.Tag, sc,
		[]byte(fmt.Sprintf(outboundProtoMap["vless"], bc.Addr.String(), bc.Port, auths.Marshal())))
}

func withOutbound(proto string, tag string, sc *conf.StreamConfig, settings []byte) Options {
	return func(cf *conf.Config) {
		cf.OutboundConfigs = append(cf.OutboundConfigs, conf.OutboundDetourConfig{
			Protocol:      proto,
			Tag:           tag,
			Settings:      tools.JsonRawPtr(settings),
			StreamSetting: sc,
		})
	}
}
