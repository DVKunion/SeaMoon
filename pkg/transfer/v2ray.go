package transfer

import (
	"context"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/app/dispatcher"
	pinboud "github.com/v2fly/v2ray-core/v5/app/proxyman/inbound"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/features/inbound"
	"github.com/v2fly/v2ray-core/v5/features/routing"
	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
)

func Init(cfg *v2rayConfig) error {
	config, err := cfg.Build()
	if err != nil {
		return err
	}
	v2ray, err = core.New(config)
	if err != nil {
		return err
	}
	return nil
}

// V2rayTransport v2ray 相关协议支持: vmess / vless / shadowsock
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
	return worker.Process(ctx, net.Network_TCP, conn, dispatch)
}
