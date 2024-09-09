// Package xray
// wrapper for using in project
package xray

import (
	"context"
	"fmt"
	"time"

	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/plugins/xray/config"
	"github.com/DVKunion/SeaMoon/plugins/xray/net"
	"github.com/DVKunion/SeaMoon/plugins/xray/service"
)

type Xray struct {
	conn *grpc.ClientConn
}

var versions = core.VersionStatement()

func GetVer() string {
	if len(versions) > 1 {
		return versions[0]
	}
	return "unknown"
}

func StartServer(opts ...config.Options) (core.Server, error) {
	cfg, err := config.Render(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "xray build error")
	}

	server, err := core.New(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "xray new error")
	}

	return server, nil
}

func (x Xray) Init() error {
	if x.conn == nil {
		cc, err := grpc.NewClient(fmt.Sprintf("%s:%d", "127.0.0.1", 10085),
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			))
		if err != nil {
			return err
		}
		x.conn = cc
	}
	return nil
}

func (x Xray) StartProxy(ctx context.Context, proxy models.Proxy, tunnel models.Tunnel) error {
	if x.conn == nil {
		return errors.New("need init grpc connection first")
	}
	handleService := service.NewHandleService(x.conn)
	// 这里用一个 err chan 来处理过多的 if err 场景
	errC := make(chan error, 1)
	doneC := make(chan struct{})
	// 配置出入信息
	go func() {
		var sc *conf.StreamConfig
		var inbound = &config.BoundConfig{
			Addr: net.ParseAddress(*proxy.ListenAddr),
			Port: net.PortFromString(*proxy.ListenPort),
			Tag:  proxy.Tag(),
		}
		var outbound = &config.BoundConfig{
			Addr: net.ParseAddress(*tunnel.Addr),
			Port: net.PortFromInt(uint32(*tunnel.Port)),
			Tag:  "",
		}
		switch *tunnel.Type {
		case enum.TunnelTypeWST:
			sc = config.WithWSSettings(tunnel.Config.TLS, *tunnel.Addr)
		case enum.TunnelTypeGRT:
			// todo
			//sc = config.WithGrpcSettings(*proxy.Type)
		}
		switch *proxy.Type {
		case enum.ProxyTypeHTTP:
			errC <- handleService.AddInbound(ctx, config.WithHttpInbound(inbound, nil))
		case enum.ProxyTypeSOCKS5:
			errC <- handleService.AddInbound(ctx, config.WithSocksInbound(inbound, nil))
			errC <- handleService.AddOutbound(ctx, config.WithSocksOutbound(outbound, nil, sc))
		case enum.ProxyTypeVmess:
			errC <- handleService.AddInbound(ctx, config.WithSocksInbound(inbound, nil))
			errC <- handleService.AddOutbound(ctx, config.WithVmessOutbound(outbound, nil, sc))
		case enum.ProxyTypeVless:
			errC <- handleService.AddInbound(ctx, config.WithSocksInbound(inbound, nil))
			errC <- handleService.AddOutbound(ctx, config.WithVlessOutbound(outbound, nil, sc))
		case enum.ProxyTypeShadowSocks:
			errC <- handleService.AddInbound(ctx, config.WithSocksInbound(inbound, nil))
			errC <- handleService.AddOutbound(ctx, config.WithVlessOutbound(outbound, nil, sc))
		case enum.ProxyTypeTorjan:
			errC <- handleService.AddInbound(ctx, config.WithSocksInbound(inbound, nil))
			errC <- handleService.AddOutbound(ctx, config.WithTorjanOutbound(outbound, nil, sc))
		default:
			errC <- errors.New("un support")
		}
	}()

	for {
		select {
		case err := <-errC:
			if err != nil {
				return err
			}
		case <-doneC:
			return nil
		// 增加一个最大的超时时间 防止创建时候卡死了
		case <-time.After(30 * time.Second):
			return errors.New("start proxy timeout")
		}
	}
}
