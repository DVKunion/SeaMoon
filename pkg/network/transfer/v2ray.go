package transfer

import (
	"context"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	pinboud "github.com/v2fly/v2ray-core/v5/app/proxyman/inbound"
	"github.com/v2fly/v2ray-core/v5/common/errors"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/features/inbound"
	"github.com/v2fly/v2ray-core/v5/features/routing"

	_ "github.com/v2fly/v2ray-core/v5/app/dispatcher"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/inbound"
	_ "github.com/v2fly/v2ray-core/v5/app/proxyman/outbound"

	_ "github.com/v2fly/v2ray-core/v5/proxy/freedom"
	_ "github.com/v2fly/v2ray-core/v5/proxy/http"
	_ "github.com/v2fly/v2ray-core/v5/proxy/shadowsocks"
	_ "github.com/v2fly/v2ray-core/v5/proxy/socks"
	_ "github.com/v2fly/v2ray-core/v5/proxy/vless/inbound"
	_ "github.com/v2fly/v2ray-core/v5/proxy/vless/outbound"
	_ "github.com/v2fly/v2ray-core/v5/proxy/vmess/inbound"
	_ "github.com/v2fly/v2ray-core/v5/proxy/vmess/outbound"

	_ "github.com/v2fly/v2ray-core/v5/proxy/shadowsocks2022"

	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func Init(cfg *v2rayConfig) error {
	config, err := cfg.Build()
	if err != nil {
		return err
	}
	v2ray, err = core.New(config)
	return err
}

// V2rayTransport v2ray 相关协议支持: vmess / vless / shadowsock
func V2rayTransport(conn net.Conn, proto string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	manager := v2ray.GetFeature(inbound.ManagerType()).(inbound.Manager)
	handler, err := manager.GetHandler(ctx, handleTag+proto)
	if err != nil {
		return err
	}
	worker := handler.(*pinboud.AlwaysOnInboundHandler).GetInbound()

	sid := session.NewID()
	ctx = session.ContextWithID(ctx, sid)
	ctx = session.ContextWithInbound(ctx, &session.Inbound{
		Source: net.DestinationFromAddr(conn.RemoteAddr()),
		Tag:    handleTag + proto,
	})

	content := new(session.Content)
	ctx = session.ContextWithContent(ctx, content)

	dispatch := v2ray.GetFeature(routing.DispatcherType()).(*dispatcher.DefaultDispatcher)
	if err = worker.Process(ctx, net.Network_TCP, conn, dispatch); err != nil {
		// 这个 err 在官方的代码里面也就那个样子，在 go 的携程里面没人管
		// 大概率是 context canceled 或者  io: read/write on closed pipe
		// 应该也是没办法处理 当远程服务断开链接时候的通信
		// 我们就处理的简单一些好了。
		if err.(*errors.Error).Severity() > 1 {
			xlog.Debug(xlog.ServiceTransportError, err)
		} else {
			return err
		}
	}
	if err = conn.Close(); err != nil {
		return err
	}
	return nil
}
