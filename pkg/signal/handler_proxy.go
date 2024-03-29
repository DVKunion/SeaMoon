package signal

import (
	"context"
	"sync"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/listener"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func (sb *Bus) SendProxySignal(p uint, tp enum.ProxyStatus) {
	sb.proxyChannel <- &proxySignal{
		id:   p,
		next: tp,
		wg:   nil,
	}
}

func (sb *Bus) SendProxySignalSync(p uint, tp enum.ProxyStatus) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	sb.proxyChannel <- &proxySignal{
		id:   p,
		next: tp,
		wg:   wg,
	}
	wg.Wait()
}

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
	if *proxy.Status == pys.next {
		xlog.Warn(xlog.SignalMissOperationWarn, "id", pys.id, "type", "proxy", "status", pys.next)
		return
	}
	service.SVC.UpdateProxyStatus(ctx, pys.id, pys.next, "")
	switch pys.next {
	case enum.ProxyStatusActive, enum.ProxyStatusRecover:
		sigCtx, cancel := context.WithCancel(ctx)
		if server, err := listener.TCPListen(sigCtx, proxy); err != nil {
			xlog.Error(xlog.SignalListenerError, "id", pys.id, "type", *proxy.Type, "addr", proxy.Addr(), "err", err)
			service.SVC.UpdateProxyStatus(ctx, pys.id, enum.ProxyStatusError, err.Error())
			cancel()
			return
		} else {
			sb.canceler[pys.id] = cancel
			sb.listener[pys.id] = server
		}
		xlog.Info(xlog.SignalStartProxy, "id", pys.id, "type", *proxy.Type, "addr", proxy.Addr())
		service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusActive, "")
	case enum.ProxyStatusInactive:
		sb.stopProxy(proxy)
	case enum.ProxyStatusDelete:
		sb.deleteProxy(ctx, proxy)
	case enum.ProxyStatusSpeeding:
		if err = service.SVC.SpeedProxy(ctx, proxy); err != nil {
			xlog.Error(xlog.SignalSpeedProxyError, "id", pys.id, "type", *proxy.Type, "addr", proxy.Addr(), "err", err)
			service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusError, err.Error())
			return
		}
		xlog.Info(xlog.SignalSpeedProxy, "id", pys.id, "type", *proxy.Type, "addr", proxy.Addr())
		service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusActive, "")
	}
}

func (sb *Bus) stopProxy(proxy *models.Proxy) {
	if cancel, ok := sb.canceler[proxy.ID]; ok {
		// 先调一下 cancel
		cancel()
		if ln, exist := sb.listener[proxy.ID]; exist {
			// 尝试着去停一下 ln, 防止泄漏
			err := ln.Close()
			if err != nil {
				// 错了就错了吧，说明 ctx 挂了一般 goroutines 也跟着挂了
				xlog.Error(xlog.SignalListenerError, "id", proxy.ID, "type", *proxy.Type, "addr", proxy.Addr(), "err", err)
			}
		}
	}
	xlog.Info(xlog.SignalStopProxy, "id", proxy.ID, "type", *proxy.Type, "addr", proxy.Addr())
}

func (sb *Bus) deleteProxy(ctx context.Context, proxy *models.Proxy) {
	sb.stopProxy(proxy)
	// 最后删除数据
	if err := service.SVC.DeleteProxy(ctx, proxy.ID); err != nil {
		xlog.Error(xlog.ServiceDBDeleteError, "id", proxy.ID, "type", *proxy.Type, "addr", proxy.Addr(), "err", err)
		service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusError, err.Error())
		return
	}
	xlog.Info(xlog.SignalDeleteProxy, "id", proxy.ID, "type", *proxy.Type, "addr", proxy.Addr())

}
