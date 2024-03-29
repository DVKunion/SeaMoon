// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.3.0
// source: tunnel.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TunnelClient is the client API for Tunnel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TunnelClient interface {
	Http(ctx context.Context, opts ...grpc.CallOption) (Tunnel_HttpClient, error)
	Socks5(ctx context.Context, opts ...grpc.CallOption) (Tunnel_Socks5Client, error)
}

type tunnelClient struct {
	cc grpc.ClientConnInterface
}

func NewTunnelClient(cc grpc.ClientConnInterface) TunnelClient {
	return &tunnelClient{cc}
}

func (c *tunnelClient) Http(ctx context.Context, opts ...grpc.CallOption) (Tunnel_HttpClient, error) {
	stream, err := c.cc.NewStream(ctx, &Tunnel_ServiceDesc.Streams[0], "/tunnel.Tunnel/Http", opts...)
	if err != nil {
		return nil, err
	}
	x := &tunnelHttpClient{stream}
	return x, nil
}

type Tunnel_HttpClient interface {
	Send(*Chunk) error
	Recv() (*Chunk, error)
	grpc.ClientStream
}

type tunnelHttpClient struct {
	grpc.ClientStream
}

func (x *tunnelHttpClient) Send(m *Chunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *tunnelHttpClient) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tunnelClient) Socks5(ctx context.Context, opts ...grpc.CallOption) (Tunnel_Socks5Client, error) {
	stream, err := c.cc.NewStream(ctx, &Tunnel_ServiceDesc.Streams[1], "/tunnel.Tunnel/Socks5", opts...)
	if err != nil {
		return nil, err
	}
	x := &tunnelSocks5Client{stream}
	return x, nil
}

type Tunnel_Socks5Client interface {
	Send(*Chunk) error
	Recv() (*Chunk, error)
	grpc.ClientStream
}

type tunnelSocks5Client struct {
	grpc.ClientStream
}

func (x *tunnelSocks5Client) Send(m *Chunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *tunnelSocks5Client) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TunnelServer is the server API for Tunnel service.
// All implementations must embed UnimplementedTunnelServer
// for forward compatibility
type TunnelServer interface {
	Http(Tunnel_HttpServer) error
	Socks5(Tunnel_Socks5Server) error
	mustEmbedUnimplementedTunnelServer()
}

// UnimplementedTunnelServer must be embedded to have forward compatible implementations.
type UnimplementedTunnelServer struct {
}

func (UnimplementedTunnelServer) Http(Tunnel_HttpServer) error {
	return status.Errorf(codes.Unimplemented, "method Http not implemented")
}
func (UnimplementedTunnelServer) Socks5(Tunnel_Socks5Server) error {
	return status.Errorf(codes.Unimplemented, "method Socks5 not implemented")
}
func (UnimplementedTunnelServer) mustEmbedUnimplementedTunnelServer() {}

// UnsafeTunnelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TunnelServer will
// result in compilation errors.
type UnsafeTunnelServer interface {
	mustEmbedUnimplementedTunnelServer()
}

func RegisterTunnelServer(s grpc.ServiceRegistrar, srv TunnelServer) {
	s.RegisterService(&Tunnel_ServiceDesc, srv)
}

func _Tunnel_Http_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TunnelServer).Http(&tunnelHttpServer{stream})
}

type Tunnel_HttpServer interface {
	Send(*Chunk) error
	Recv() (*Chunk, error)
	grpc.ServerStream
}

type tunnelHttpServer struct {
	grpc.ServerStream
}

func (x *tunnelHttpServer) Send(m *Chunk) error {
	return x.ServerStream.SendMsg(m)
}

func (x *tunnelHttpServer) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Tunnel_Socks5_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TunnelServer).Socks5(&tunnelSocks5Server{stream})
}

type Tunnel_Socks5Server interface {
	Send(*Chunk) error
	Recv() (*Chunk, error)
	grpc.ServerStream
}

type tunnelSocks5Server struct {
	grpc.ServerStream
}

func (x *tunnelSocks5Server) Send(m *Chunk) error {
	return x.ServerStream.SendMsg(m)
}

func (x *tunnelSocks5Server) Recv() (*Chunk, error) {
	m := new(Chunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Tunnel_ServiceDesc is the grpc.ServiceDesc for Tunnel service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Tunnel_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tunnel.Tunnel",
	HandlerType: (*TunnelServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Http",
			Handler:       _Tunnel_Http_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Socks5",
			Handler:       _Tunnel_Socks5_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "tunnel.proto",
}
