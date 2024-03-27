package enum

import "strings"

type ProxyStatus int8

const (
	ProxyStatusInitializing ProxyStatus = iota + 1
	ProxyStatusActive
	ProxyStatusInactive
	ProxyStatusError
	ProxyStatusSpeeding
	ProxyStatusRecover
	ProxyStatusDelete
)

type ProxyType string

const (
	ProxyTypeAUTO        ProxyType = "auto"
	ProxyTypeHTTP        ProxyType = "http"
	ProxyTypeSOCKS5      ProxyType = "socks5"
	ProxyTypeSOCKS5Ssr   ProxyType = "socks5+ssr"
	ProxyTypeSOCKS5Vmess ProxyType = "socks5+vmess"
	ProxyTypeSOCKS5Vless ProxyType = "socks5+vless"
)

func (t ProxyType) String() string {
	return string(t)
}

func (t ProxyType) ProtoString() string {
	return t.String() + "://"
}

func (t ProxyType) Path(p string) string {
	if strings.HasSuffix(p, "/") {
		return p + t.String()
	} else {
		return p + "/" + t.String()
	}
}
