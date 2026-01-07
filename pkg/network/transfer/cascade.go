package transfer

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/DVKunion/SeaMoon/pkg/network/basic"
	"github.com/DVKunion/SeaMoon/pkg/network/tunnel"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

// CascadeTransport 级联代理转发，直接把流量转发给下一跳处理
// 参考 TorTransport 的简洁实现，入口流量直接转发给下一跳的对应服务
func CascadeTransport(conn net.Conn, path string) error {
	cfg := GetCascadeConfig()
	if cfg == nil || !cfg.Enabled {
		xlog.Error("cascade proxy: not enabled")
		return nil
	}

	// 解析级联代理地址并建立 websocket 连接
	wsAddr := parseCascadeAddress(cfg.Addr) + "/" + path

	xlog.Info(xlog.ServiceCasCadeConnectServer, "src", conn.RemoteAddr(), "dest", wsAddr, "mode", "cascade")

	// 建立 websocket 连接到下一跳
	wsDialer := &websocket.Dialer{
		HandshakeTimeout: 30 * time.Second,
	}

	wsConn, _, err := wsDialer.Dial(wsAddr, http.Header{})
	if err != nil {
		xlog.Error("cascade proxy: failed to dial next hop", "err", err, "addr", wsAddr)
		return err
	}

	// 包装 websocket 连接
	nextHop := tunnel.WsWrapConn(wsConn)
	defer nextHop.Close()

	// 直接双向转发，让下一跳处理所有协议
	if _, _, err := basic.Transport(conn, nextHop); err != nil {
		xlog.Error(xlog.NetworkTransportError, "err", err)
	}

	xlog.Info(xlog.ServiceCasCadeDisConnect, "src", conn.RemoteAddr(), "dest", wsAddr, "mode", "cascade")

	return nil
}

// parseCascadeAddress 解析级联代理地址，处理各种前缀格式
// 输入格式可能是: wss://example.com, ws://example.com, https://example.com, example.com
// 输出格式为: wss://example.com 或 ws://example.com
func parseCascadeAddress(addr string) string {
	// 如果是 https 前缀，转换为 wss
	if strings.HasPrefix(addr, "https://") {
		return "wss://" + strings.TrimPrefix(addr, "https://")
	}

	// 如果是 http 前缀，转换为 ws
	if strings.HasPrefix(addr, "http://") {
		return "ws://" + strings.TrimPrefix(addr, "http://")
	}

	return addr
}
