package transfer

import (
	"fmt"
	"testing"

	core "github.com/v2fly/v2ray-core/v5"
)

func TestV2rayConfig(t *testing.T) {
	config, err := core.LoadConfig("json", []byte(`{
  "inbounds": [
    {
      "port": 10000,
      "tag": "xixixi",
      "listen":"127.0.0.1",
      "protocol": "shadowsocks2022",
      "settings": {
			"method": "aes-256-gcm",
			"password": "123456",
			"network": "tcp"
		},
      "streamSettings": {
        "network": "ws",
		"security": "tls",
        "wsSettings": {
          "path": "/vlite"
        }
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ]
}`))
	fmt.Println(config, err)
}
