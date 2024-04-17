package tunnel

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	proto2 "github.com/DVKunion/SeaMoon/pkg/network/tunnel/service/proto"
)

type grpcConn struct {
	cc grpc.Stream

	rb    []byte
	lAddr net.Addr
	rAddr net.Addr
}

func GRPCWrapConn(addr net.Addr, cc grpc.Stream) Tunnel {
	return &grpcConn{
		cc:    cc,
		lAddr: addr,
		rAddr: &net.TCPAddr{},
	}
}

func (c *grpcConn) Delay() int64 {
	return 0
}

func (c *grpcConn) Read(b []byte) (n int, err error) {
	if len(c.rb) == 0 {
		chunk, err := c.recv()
		if err != nil {
			return 0, err
		}
		c.rb = chunk.Body
	}

	n = copy(b, c.rb)
	c.rb = c.rb[n:]
	return
}

func (c *grpcConn) Write(b []byte) (n int, err error) {
	chunk := &proto2.Chunk{
		Body: b,
		Size: int32(len(b)),
	}

	if err = c.send(chunk); err != nil {
		return
	}

	n = int(chunk.Size)
	return
}

func (c *grpcConn) Close() error {
	switch cost := c.cc.(type) {
	case proto2.Tunnel_HttpClient:
	case proto2.Tunnel_Socks5Client:
		return cost.CloseSend()
	}
	return nil
}

func (c *grpcConn) LocalAddr() net.Addr {
	return c.lAddr
}

func (c *grpcConn) RemoteAddr() net.Addr {
	return c.rAddr
}

func (c *grpcConn) SetDeadline(t time.Time) error {
	return &net.OpError{Op: "set", Net: "grpc", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (c *grpcConn) SetReadDeadline(t time.Time) error {
	return &net.OpError{Op: "set", Net: "grpc", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (c *grpcConn) SetWriteDeadline(t time.Time) error {
	return &net.OpError{Op: "set", Net: "grpc", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (c *grpcConn) context() context.Context {
	if c.cc != nil {
		return c.cc.Context()
	}
	return context.Background()
}

func (c *grpcConn) send(data *proto2.Chunk) error {
	sender, ok := c.cc.(interface {
		Send(*proto2.Chunk) error
	})
	if !ok {
		// todo
		return fmt.Errorf("unsupported type: %T", c.cc)
	}
	return sender.Send(data)
}

func (c *grpcConn) recv() (*proto2.Chunk, error) {
	receiver, ok := c.cc.(interface {
		Recv() (*proto2.Chunk, error)
	})
	if !ok {
		// todo
		return nil, fmt.Errorf("unsupported type: %T", c.cc)
	}
	return receiver.Recv()
}
