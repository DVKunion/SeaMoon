package service

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	pb "github.com/DVKunion/SeaMoon/pkg/proto"
	"github.com/DVKunion/SeaMoon/pkg/transfer"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type GRPCService struct {
	addr   net.Addr
	cc     *grpc.ClientConn
	server *grpc.Server

	pb.UnimplementedTunnelServer
}

func init() {
	register(tunnel.GRT, &GRPCService{})
}

func (g GRPCService) Conn(ctx context.Context, t transfer.Type, sOpts ...Option) (net.Conn, error) {
	var cs grpc.ClientStream
	var srvOpts = &Options{}
	var err error

	for _, o := range sOpts {
		o(srvOpts)
	}

	if strings.HasPrefix(srvOpts.addr, "grpc://") {
		srvOpts.addr = strings.TrimPrefix(srvOpts.addr, "grpc://")
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
	case transfer.HTTP:
		cs, err = client.Http(ctx)
	case transfer.SOCKS5:
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

	return server.Serve(ln)
}

func (g GRPCService) Http(server pb.Tunnel_HttpServer) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.HttpTransport(gt); err != nil {
		slog.Error("connection error", "msg", err)
		return err
	}

	return nil
}

func (g GRPCService) Socks5(server pb.Tunnel_Socks5Server) error {
	gt := tunnel.GRPCWrapConn(g.addr, server)

	if err := transfer.Socks5Transport(gt); err != nil {
		slog.Error("connection error", "msg", err)
		return err
	}
	return nil
}
