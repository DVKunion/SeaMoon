package listener

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"sync"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	db_service "github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/pkg/network"
	"github.com/DVKunion/SeaMoon/pkg/service"
	"github.com/DVKunion/SeaMoon/pkg/xlog"
)

func Listen(ctx context.Context, server net.Listener, p uint) {
	var pro *models.Proxy
	var tun *models.Tunnel

	// 应用级别事务锁
	var m = &sync.Mutex{}

	dbProxy := db_service.GetService("proxy")
	dbTunnel := db_service.GetService("tunnel")

	objP := dbProxy.GetById(p)
	if v, ok := objP.(*models.Proxy); ok {
		pro = v
	} else {
		*pro.Status = types.ERROR
		dbProxy.Update(pro.ID, pro)
		slog.Error("proxy error")
		return
	}
	objT := dbTunnel.GetById(pro.TunnelID)
	if v, ok := objT.(*models.Tunnel); ok {
		tun = v
	} else {
		*pro.Status = types.ERROR
		dbProxy.Update(pro.ID, pro)
		slog.Error("tunnel error")
		return
	}
	*pro.Status = types.ACTIVE
	dbProxy.Update(pro.ID, pro)
	for {
		conn, err := server.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				// 说明是 server 被外部 close 掉了，导致了此处的 accept 报错
				// 正常现象，return 即可。
				return
			} else {
				// 除此之外，都为异常。为了保证服务正常不出现 panic 和空指针，跳过该 conn
				xlog.Error("SERVE", xlog.ACCEPT_ERROR, "err", err)
				continue
			}
		}
		// 说明接到了一个conn, 更新
		count(dbProxy, pro.ID, 1, m)
		if srv, ok := service.Factory.Load(*tun.Type); ok {
			destConn, err := srv.(service.Service).Conn(ctx, *pro.Type,
				service.WithAddr(tun.GetAddr()), service.WithTorFlag(tun.TunnelConfig.Tor))
			if err != nil {
				slog.Error(xlog.CONNECT_RMOET_ERROR, "err", err)
				// 说明远程连接失败了，直接把当前的这个 conn 关掉，然后数量 -1
				//_ = conn.Close()
				count(dbProxy, pro.ID, -1, m)
				continue
			}
			go func() {
				in, out, err := network.Transport(conn, destConn)
				if err != nil {
					slog.Error(xlog.CONNECT_TRANS_ERROR, "err", err)
				}
				count(dbProxy, pro.ID, -1, m, in, out)
			}()
		}
	}
}

func count(svc db_service.ApiService, id uint, cnt int, m *sync.Mutex, args ...int64) {
	// 查询 -> 原有基础上 +1/-1 -> 更新 是一个完整的事务
	// 只有这一整个操作完成后，才应该进行下一个操作。
	// 这里丑陋的用了一个 mutex 来实现这个问题，正常应该通过 orm 事务来操作。
	m.Lock()
	proxy := svc.GetById(id).(*models.Proxy)
	*proxy.Conn = *proxy.Conn + cnt
	if len(args) > 0 {
		*proxy.InBound = *proxy.InBound + args[0]
		*proxy.OutBound = *proxy.OutBound + args[1]
	}
	svc.Update(id, proxy)
	m.Unlock()
}
