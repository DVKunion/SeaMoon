package proxy

import (
	"context"
	"net/http"
)

type SocksAction string

const (
	CONNECT    SocksAction = "connect"
	DISCONNECT SocksAction = "disconnect"
	READ       SocksAction = "read"
	FORWARD    SocksAction = "forward"
)

// AliYunSocks5Handler 阿里云Socks5代理入口
func AliYunSocks5Handler(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	// 将 http req 解析为 socks
	action := req.Header.Get("SM-SOCKS")
	//mark := action[0:22]
	action = action[22:]
	//run := "run" + mark
	//writebuf := "writebuf" + mark
	//readbuf := "readbuf" + mark
	switch action {
	case string(CONNECT):
		break
	case string(DISCONNECT):
		break
	case string(READ):
		//readBuffer := $_SESSION[$readbuf];
		break
	case string(FORWARD):
		break
	default:
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/plain")
		_, err := w.Write([]byte("SeaMoon Listening For The Socks5 Request"))
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func NewSocks5Client(listenAddr string, proxyAddr string, verbose bool) {

}

func parseHttp2Socks(action string) {
}
