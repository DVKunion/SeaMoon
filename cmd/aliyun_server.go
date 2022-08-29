package main

import (
	"SeaMoon/pkg/proxy"
	"github.com/aliyun/fc-runtime-go-sdk/fc"
	"os"
)

var (
	serverMod = os.Getenv("serverMod")
)

func main() {
	switch serverMod {
	case "http":
		fc.StartHttp(proxy.AliYunHttpHandler)
		return
	case "socks5":
		fc.StartHttp(proxy.AliYunSocks5Handler)
		return
	default:
		return
	}
}
