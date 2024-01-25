package client

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strings"

	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/service"
	"github.com/DVKunion/SeaMoon/pkg/transfer"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

func Control(ctx context.Context, sg *SigGroup) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	for {
		select {
		case t := <-sg.StartChannel:
			slog.Info(consts.LISTEN_START, "type", t)
			sg.wg.Add(1)
			if err := doListen(c, t); err != nil {
				slog.Error(consts.LISTEN_ERROR, "type", t, "err", err)
			}
			sg.wg.Done()
		case t := <-sg.StopChannel:
			slog.Info(consts.LISTEN_STOP, "type", t)
			cancel()
		}
	}
}

func doListen(ctx context.Context, t transfer.Type) error {
	server, err := net.Listen("tcp", Config().Addr(t))
	if err != nil {
		return err
	}
	defer server.Close()
	var proxyAddr string
	var proxyType tunnel.Type
	for _, p := range Config().ProxyAddr {
		if strings.HasPrefix(p, "ws://") {
			proxyAddr = strings.TrimPrefix(p, "ws://")
			proxyType = tunnel.WST
			break
		}
		if strings.HasPrefix(p, "grpc://") {
			proxyAddr = p
			proxyType = tunnel.GRT
			break
		}
	}
	if proxyAddr == "" || proxyType == "" {
		return errors.New(consts.PROXY_ADDR_ERROR)
	}
	go listen(ctx, server, proxyAddr, proxyType, t)
	<-ctx.Done()
	return nil
}

func listen(ctx context.Context, server net.Listener, pa string, pt tunnel.Type, t transfer.Type) {
	for {
		conn, err := server.Accept()
		if err != nil {
			slog.Error(consts.ACCEPT_ERROR, "err", err)
		}
		if srv, ok := service.Factory.Load(pt); ok {
			destConn, err := srv.(service.Service).Conn(ctx, t, service.WithAddr(pa))
			if err != nil {
				slog.Error(consts.CONNECT_RMOET_ERROR, "err", err)
				continue
			}
			go func() {
				if err := network.Transport(conn, destConn); err != nil {
					slog.Error(consts.CONNECT_TRANS_ERROR, "err", err)
				}
			}()
		}
	}
}
