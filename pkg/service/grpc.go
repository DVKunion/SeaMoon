package service

import (
	"context"
	"crypto/tls"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	pb "github.com/DVKunion/SeaMoon/pkg/service/proto"
	"github.com/DVKunion/SeaMoon/pkg/service/proto/gost"
	"github.com/DVKunion/SeaMoon/pkg/system/consts"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
	"github.com/DVKunion/SeaMoon/pkg/transfer"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type GRPCService struct {
	addr    net.Addr
	cc      *grpc.ClientConn
	server  *grpc.Server
	startAt time.Time
	pb.UnimplementedTunnelServer
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

	if strings.HasPrefix(srvOpts.addr, "grpc://") {
		srvOpts.addr = strings.TrimPrefix(srvOpts.addr, "grpc://")
	}

	if strings.HasPrefix(srvOpts.addr, "grpcs://") {
		srvOpts.addr = strings.TrimPrefix(srvOpts.addr, "grpcs://")
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

	client := pb.NewTunnelClient(g.cc)

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

	pb.RegisterTunnelServer(server, &g)
	gost.RegisterGostTunelServer(server, &g)

	g.startAt = time.Now()
	return server.Serve(ln)
}

func (g GRPCService) Auto(server pb.Tunnel_AutoServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.AutoTransport(gt); err != nil {
		xlog.Error(errors.ServiceTransportError, "type", "socks5", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) Http(server pb.Tunnel_HttpServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.HttpTransport(gt); err != nil {
		xlog.Error(errors.ServiceTransportError, "type", "http", "err", err)
		return err
	}

	return nil
}

func (g GRPCService) Socks5(server pb.Tunnel_Socks5Server) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.Socks5Transport(gt, false); err != nil {
		xlog.Error(errors.ServiceTransportError, "type", "socks5", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) V2RaySsr(server pb.Tunnel_V2RaySsrServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.V2rayTransport(gt, "shadowsocks"); err != nil {
		xlog.Error(errors.ServiceTransportError, "type", "v2ray-ssr", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) V2RayVmess(server pb.Tunnel_V2RayVmessServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.V2rayTransport(gt, "vmess"); err != nil {
		xlog.Error(errors.ServiceTransportError, "type", "vmess", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) V2RayVless(server pb.Tunnel_V2RayVlessServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.V2rayTransport(gt, "vless"); err != nil {
		xlog.Error(errors.ServiceTransportError, "type", "vless", "err", err)
		return err
	}
	return nil
}

// Tunnel gost grpc 适配, 实际上直接做一个 auto 协议就好了
func (g GRPCService) Tunnel(server gost.GostTunel_TunnelServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.AutoTransport(gt); err != nil {
		xlog.Error(errors.ServiceTransportError, "type", "socks5", "err", err)
		return err
	}
	return nil
}

func (g GRPCService) Health(ctx context.Context, p *pb.Ping) (*pb.Pong, error) {
	return &pb.Pong{
		Status:  "OK",
		Time:    g.startAt.Format("2006-01-02 15:04:05"),
		Version: consts.Version,
		Commit:  consts.Commit,
	}, nil
}
