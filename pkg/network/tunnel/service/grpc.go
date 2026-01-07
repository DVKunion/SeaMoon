package service

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/network/transfer"
	"github.com/DVKunion/SeaMoon/pkg/network/tunnel"
	"github.com/DVKunion/SeaMoon/pkg/network/tunnel/service/proto"
	"github.com/DVKunion/SeaMoon/pkg/network/tunnel/service/proto/gost"
	"github.com/DVKunion/SeaMoon/pkg/system/version"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

type GRPCService struct {
	addr    net.Addr
	cc      *grpc.ClientConn
	server  *grpc.Server
	startAt time.Time
	udpAddr string
	proto.UnimplementedTunnelServer
	gost.UnimplementedGostTunelServer
}

func init() {
	register(enum.TunnelTypeGRT, &GRPCService{})
}

func (g GRPCService) Conn(ctx context.Context, t enum.ProxyType, sOpts ...Option) (tunnel.Tunnel, error) {
	var cs grpc.ClientStream
	var srvOpts = &Options{}
	var err error

	for _, o := range sOpts {
		o(srvOpts)
	}

	if after, ok := strings.CutPrefix(srvOpts.addr, "grpc://"); ok {
		srvOpts.addr = after
	}

	if after, ok := strings.CutPrefix(srvOpts.addr, "grpcs://"); ok {
		srvOpts.addr = after
	}

	if srvOpts.udpAddr != "" {
		g.udpAddr = srvOpts.udpAddr
		_ = g.udpAddr
	}

	nAddr, err := net.ResolveTCPAddr("tcp", srvOpts.addr)
	if err != nil {
		return nil, err
	}

	if g.cc == nil {
		// do connect
		grpcOpts := []grpc.DialOption{
			//grpc.WithAuthority(host),
			//grpc.WithConnectParams(grpc.ConnectParams{
			//	Backoff: backoff.DefaultConfig,
			//MinConnectTimeout: d.md.minConnectTimeout,
			//}),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
				Timeout:             3 * time.Second,  // wait 1 second for ping ack before considering the c
				PermitWithoutStream: false,            // send pings even without active streams
			}),
			//grpc.FailOnNonTempDialError(true),
		}

		//if !d.md.insecure {
		//	grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(d.options.TLSConfig)))
		//} else {
		//grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})))
		//}
		g.cc, err = grpc.DialContext(ctx, srvOpts.addr, grpcOpts...)
		if err != nil {
			return nil, err
		}
	}

	client := proto.NewTunnelClient(g.cc)

	switch t {
	case enum.ProxyTypeHTTP:
		cs, err = client.Http(ctx)
	case enum.ProxyTypeSOCKS5:
		cs, err = client.Socks5(ctx)
	}

	if err != nil {
		return nil, err
	}

	return tunnel.GRPCWrapConn(nAddr, cs), nil

}

func (g GRPCService) Serve(ln net.Listener, srvOpt ...Option) error {
	var srvOpts = &Options{}
	for _, o := range srvOpt {
		o(srvOpts)
	}
	var gRPCOpts []grpc.ServerOption
	if srvOpts.tlsConf != nil {
		gRPCOpts = append(gRPCOpts, grpc.Creds(credentials.NewTLS(srvOpts.tlsConf)))
	}

	if srvOpts.keepalive != nil {
		gRPCOpts = append(gRPCOpts,
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Time:              10 * time.Second, // send pings every 10 seconds if there is no activity
				Timeout:           3 * time.Second,  // wait 1 second for ping ack before considering the connection dead
				MaxConnectionIdle: 30 * time.Second,
				//MaxConnectionIdle: srvOpts.keepalive.MaxConnectionIdle,
				//Time:              srvOpts.keepalive.MaxTime,
				//Timeout:           srvOpts.keepalive.Timeout,
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				//MinTime:             srvOpts.keepalive.MinTime,
				PermitWithoutStream: false,
			}),
		)
	}

	server := grpc.NewServer(gRPCOpts...)

	addr := strings.Split(ln.Addr().String(), ":")
	var port = 443 // 默认端口
	if len(addr) > 1 {
		if p, err := strconv.Atoi(addr[1]); err == nil {
			port = p
		}
	}

	// init config
	// 构建 v2ray 配置选项
	configOpts := []transfer.ConfigOpt{
		transfer.WithServerMod(),
		transfer.WithNetAddr("0.0.0.0", uint32(port)),
		transfer.WithTunnelType("", enum.TunnelTypeWST),
		transfer.WithAuthInfo(srvOpts.uid, srvOpts.crypt, srvOpts.pass),
		transfer.WithExtra(srvOpts.tor, srvOpts.tlsConf != nil),
	}

	// 如果配置了级联代理，添加到配置中，并设置全局配置供 socks5/http 原生转发使用
	if srvOpts.cascadeAddr != "" && srvOpts.cascadeUid != "" {
		configOpts = append(configOpts, transfer.WithCascadeProxy(srvOpts.cascadeAddr, srvOpts.cascadeUid, srvOpts.cascadePassword))
		// 设置全局级联代理配置
		transfer.SetCascadeConfig(srvOpts.cascadeAddr, srvOpts.cascadeUid, srvOpts.cascadePassword)
	}

	config := transfer.NewV2rayConfig(configOpts...)

	if err := transfer.Init(config); err != nil {
		xlog.Error(xlog.ServiceV2rayInitError, "err", err)
	}

	proto.RegisterTunnelServer(server, &g)
	gost.RegisterGostTunelServer(server, &g)

	g.startAt = time.Now()
	return server.Serve(ln)
}

func (g GRPCService) Auto(server proto.Tunnel_AutoServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.AutoTransport(gt); err != nil {
		xlog.Error(xlog.ServiceTransportError, "type", "socks5", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) Http(server proto.Tunnel_HttpServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.HttpTransport(gt); err != nil {
		xlog.Error(xlog.ServiceTransportError, "type", "http", "err", err)
		return err
	}

	return nil
}

func (g GRPCService) Socks5(server proto.Tunnel_Socks5Server) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.Socks5Transport(gt, false, g.udpAddr); err != nil {
		xlog.Error(xlog.ServiceTransportError, "type", "socks5", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) V2RaySsr(server proto.Tunnel_V2RaySsrServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.V2rayTransport(gt, "shadowsocks"); err != nil {
		xlog.Error(xlog.ServiceTransportError, "type", "v2ray-ssr", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) V2RayVmess(server proto.Tunnel_V2RayVmessServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.V2rayTransport(gt, "vmess"); err != nil {
		xlog.Error(xlog.ServiceTransportError, "type", "vmess", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) V2RayVless(server proto.Tunnel_V2RayVlessServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.V2rayTransport(gt, "vless"); err != nil {
		xlog.Error(xlog.ServiceTransportError, "type", "vless", "err", err)
		return err
	}
	return nil
}

// Tunnel gost grpc 适配, 实际上直接做一个 auto 协议就好了
func (g GRPCService) Tunnel(server gost.GostTunel_TunnelServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.AutoTransport(gt); err != nil {
		xlog.Error(xlog.ServiceTransportError, "type", "socks5", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) Health(ctx context.Context, p *proto.Ping) (*proto.Pong, error) {
	return &proto.Pong{
		Status:  "OK",
		Time:    g.startAt.Format("2006-01-02 15:04:05"),
		Version: version.Version,
		Commit:  version.Commit,
	}, nil
}
