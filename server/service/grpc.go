package service

import (
	"log/slog"
	"net"
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
	server *grpc.Server

	pb.UnimplementedTunnelServer
}

func init() {
	register(tunnel.GRT, &GRPCService{})
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
