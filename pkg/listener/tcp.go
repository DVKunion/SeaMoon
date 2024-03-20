package listener

import (
	"context"
	"net"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	db_service "github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/service"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
	"github.com/DVKunion/SeaMoon/pkg/tools"
)

func TCPListen(ctx context.Context, py *models.Proxy, tun *models.Tunnel) (net.Listener, error) {
	server, err := net.Listen("tcp", py.Addr())
	if err != nil {
		return nil, err
	}
	go listen(ctx, server, py.ID, py.Type, tun)
	return server, nil
}

func listen(ctx context.Context, server net.Listener, id uint, t *enum.ProxyType, tun *models.Tunnel) {
	for {
		conn, err := server.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				// 说明是 server 被外部 close 掉了，导致了此处的 accept 报错
				// 正常现象，return 即可。
				return
			} else {
				// 除此之外，都为异常。为了保证服务正常不出现 panic 和空指针，跳过该 conn
				xlog.Error(errors.ListenerAcceptError, "err", err)
				continue
			}
		}
		if err = db_service.SVC.UpdateProxyConn(ctx, id, 1); err != nil {
			// todo: do log
		}
		if srv, ok := service.Factory.Load(*tun.Type); ok {
			destConn, err := srv.(service.Service).Conn(ctx, *t,
				service.WithAddr(tun.GetAddr()), service.WithTorFlag(tun.Config.Tor))
			if err != nil {
				xlog.Error(errors.ListenerDailError, "err", err)
				if err = db_service.SVC.UpdateProxyConn(ctx, id, -1); err != nil {
					// todo: do log
				}
				continue
			}
			go func() {
				if _, err = db_service.SVC.UpdateProxy(ctx, id, &models.Proxy{
					Lag: tools.Int64Ptr(destConn.Delay()),
				}); err != nil {
					xlog.Error(errors.ListenerLagError, "id", id, "err", err)
				}
			}()
			go func() {
				in, out, err := network.Transport(conn, destConn)
				if err != nil {
					xlog.Error(errors.NetworkTransportError, "err", err)
				}
				if err = db_service.SVC.UpdateProxyConn(ctx, id, -1); err != nil {
					// todo: do log
				}
				if err = db_service.SVC.UpdateProxyNetFlow(ctx, id, in, out); err != nil {
					// todo: do log
				}
			}()
		}
	}
}
