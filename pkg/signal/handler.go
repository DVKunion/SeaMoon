package signal

import (
	"context"
	"sync"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/listener"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func (sb *Bus) proxyHandler(ctx context.Context, pys *proxySignal) {
	// proxy sync change task
	// 如果是需要同步的，记得释放锁
	defer func() {
		if pys.wg != nil {
			pys.wg.Done()
		}
	}()
	proxy, err := service.SVC.GetProxyById(ctx, pys.id)
	if err != nil {
		xlog.Error(xlog.SignalGetObjError, "obj", "proxy", "err", err)
		service.SVC.UpdateProxyStatus(ctx, pys.id, enum.ProxyStatusError, err.Error())
		return
	}
	// 缓冲逻辑：状态没改变时候，不需要处理
	if proxy.Status == &pys.next {
		xlog.Warn(xlog.SignalMissOperationWarn, "id", pys.id, "type", "proxy", "status", pys.next)
		return
	}
	service.SVC.UpdateProxyStatus(ctx, pys.id, pys.next, "")
	switch pys.next {
	case enum.ProxyStatusActive, enum.ProxyStatusRecover:
		sigCtx, cancel := context.WithCancel(ctx)
		if server, err := listener.TCPListen(sigCtx, proxy); err != nil {
			xlog.Error(xlog.SignalListenerError, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr(), "err", err)
			service.SVC.UpdateProxyStatus(ctx, pys.id, enum.ProxyStatusError, err.Error())
			cancel()
			return
		} else {
			sb.canceler[pys.id] = cancel
			sb.listener[pys.id] = server
		}
		xlog.Info(xlog.SignalStartProxy, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr())
		service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusActive, "")
	case enum.ProxyStatusInactive:
		if cancel, ok := sb.canceler[pys.id]; ok {
			// 先调一下 cancel
			cancel()
			if ln, exist := sb.listener[pys.id]; exist {
				// 尝试着去停一下 ln, 防止泄漏
				err := ln.Close()
				if err != nil {
					// 错了就错了吧，说明 ctx 挂了一般 goroutines 也跟着挂了
					xlog.Error(xlog.SignalListenerError, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr(), "err", err)
				}
			}
		}
		xlog.Info(xlog.SignalStopProxy, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr())
	case enum.ProxyStatusDelete:
		// 先同步停止服务
		wg := &sync.WaitGroup{}
		wg.Add(1)
		sb.SendProxySignal(pys.id, enum.ProxyStatusInactive, wg)
		wg.Wait()
		// 最后删除数据
		if err = service.SVC.SpeedProxy(ctx, proxy); err != nil {
			xlog.Error(xlog.SignalSpeedProxyError, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr(), "err", err)
			service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusError, err.Error())
			return
		}
		xlog.Info(xlog.SignalDeleteProxy, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr())
	case enum.ProxyStatusSpeeding:
		if err = service.SVC.SpeedProxy(ctx, proxy); err != nil {
			xlog.Error(xlog.SignalSpeedProxyError, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr(), "err", err)
			service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusError, err.Error())
			return
		}
		xlog.Info(xlog.SignalSpeedProxy, "id", pys.id, "type", proxy.Type, "addr", proxy.Addr())
		service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusActive, "")
	}
}

func (sb *Bus) providerHandler(ctx context.Context, prs *providerSignal) {
	// proxy sync change task
	// 如果是需要同步的，记得释放锁
	defer func() {
		if prs.wg != nil {
			prs.wg.Done()
		}
	}()
	provider, err := service.SVC.GetProviderById(ctx, prs.id)
	if err != nil {
		xlog.Error(xlog.SignalGetObjError, "obj", "provider", "err", err)
		service.SVC.UpdateProviderStatus(ctx, prs.id, enum.ProvStatusFailed, err.Error())
		return
	}
	// 缓冲逻辑：状态没改变时候，不需要处理
	if provider.Status == &prs.next {
		xlog.Warn(xlog.SignalMissOperationWarn, "id", prs.id, "type", "provider", "status", prs.next)
		return
	}
	service.SVC.UpdateProviderStatus(ctx, provider.ID, prs.next, "")
	switch prs.next {
	case enum.ProvStatusSync:
		if err = service.SVC.SyncProvider(ctx, provider); err != nil {
			xlog.Error(xlog.SignalSyncProviderError, "obj", "provider", "err", err)
			service.SVC.UpdateProviderStatus(ctx, provider.ID, enum.ProvStatusSyncError, err.Error())
			return
		}
		xlog.Info(xlog.SignalSyncProvider, "id", provider.ID, "type", provider.Type)
		service.SVC.UpdateProviderStatus(ctx, provider.ID, enum.ProvStatusSuccess, "")
	case enum.ProvStatusDelete:
		wg := &sync.WaitGroup{}
		for _, tun := range provider.Tunnels {
			wg.Add(1)
			sb.SendTunnelSignal(tun.ID, enum.TunnelDelete, wg)
		}
		wg.Wait()
		// 然后删除数据
		if err = service.SVC.DeleteProvider(ctx, provider.ID); err != nil {
			xlog.Error(xlog.SignalSyncProviderError, "obj", "provider", "err", err)
			service.SVC.UpdateProviderStatus(ctx, provider.ID, enum.ProvStatusSyncError, err.Error())
			return
		}
		xlog.Info(xlog.SignalDeleteProvider, "id", provider.ID, "type", provider.Type)
	}
}

func (sb *Bus) tunnelHandler(ctx context.Context, ts *tunnelSignal) {
	// proxy sync change task
	// 如果是需要同步的，记得释放锁
	defer func() {
		if ts.wg != nil {
			ts.wg.Done()
		}
	}()
	tun, err := service.SVC.GetTunnelById(ctx, ts.id)
	if err != nil {
		xlog.Error(xlog.SignalGetObjError, "obj", "tunnel", "err", err)
		service.SVC.UpdateTunnelStatus(ctx, ts.id, enum.TunnelError, err.Error())
		return
	}
	// 缓冲逻辑：状态没改变时候，不需要处理
	if tun.Status == &ts.next {
		xlog.Warn(xlog.SignalMissOperationWarn, "id", ts.id, "type", "tunnel", "status", ts.next)
		return
	}
	service.SVC.UpdateTunnelStatus(ctx, tun.ID, ts.next, "")
	switch ts.next {
	case enum.TunnelActive:
		if addr, err := service.SVC.DeployTunnel(ctx, tun); err != nil {
			xlog.Error(xlog.SignalDeployTunError, "obj", "tunnel", "err", err)
			service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
			return
		} else {
			service.SVC.UpdateTunnelAddr(ctx, tun.ID, addr)
		}
		xlog.Info(xlog.SignalDeployTunnel, "id", tun.ID, "type", tun.Type)
	case enum.TunnelInactive:
		if err := service.SVC.StopTunnel(ctx, tun); err != nil {
			xlog.Error(xlog.SignalStopTunError, "obj", "tunnel", "err", err)
			service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
			return
		}
		xlog.Info(xlog.SignalStopTunnel, "id", tun.ID, "type", tun.Type)
	case enum.TunnelDelete:
		// 先停掉本地的服务
		wg := &sync.WaitGroup{}
		for _, py := range tun.Proxies {
			wg.Add(1)
			sb.SendProxySignal(py.ID, enum.ProxyStatusDelete, wg)
		}
		wg.Wait()
		wg = &sync.WaitGroup{}
		// 再停掉远端的服务
		wg.Add(1)
		sb.SendTunnelSignal(tun.ID, enum.TunnelInactive, wg)
		wg.Wait()
		// 最后删除服务即可
		if err := service.SVC.DeleteTunnel(ctx, tun.ID); err != nil {
			xlog.Error(xlog.SignalDeleteTunError, "obj", "tunnel", "err", err)
			service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
			return
		}
		xlog.Info(xlog.SignalDeleteTunnel, "id", tun.ID, "type", tun.Type)
	}
}
