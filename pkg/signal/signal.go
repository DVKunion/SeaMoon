package signal

import (
	"context"
	"net"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/listener"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

// Bus 用于控制所有需要异步处理的状态转换
type Bus struct {
	canceler map[uint]context.CancelFunc
	listener map[uint]net.Listener

	proxyChannel    chan *proxySignal
	providerChannel chan *providerSignal
	tunnelChannel   chan *tunnelSignal
}

type proxySignal struct {
	id   uint
	next enum.ProxyStatus
}

type providerSignal struct {
	id   uint
	next enum.ProviderStatus
}

type tunnelSignal struct {
	id   uint
	next enum.TunnelStatus
}

var signalBus = &Bus{
	canceler:        make(map[uint]context.CancelFunc, 0),
	listener:        make(map[uint]net.Listener, 0),
	proxyChannel:    make(chan *proxySignal, 1>>8),
	providerChannel: make(chan *providerSignal, 1>>8),
	tunnelChannel:   make(chan *tunnelSignal, 1>>8),
}

func Signal() *Bus {
	return signalBus
}

// Daemon 控制总线守护进程
func (sb *Bus) Daemon(ctx context.Context) {
	for {
		select {
		case pys := <-sb.proxyChannel:
			// proxy sync change task
			proxy, err := service.SVC.GetProxyById(ctx, pys.id)
			if err != nil {
				xlog.Error(errors.SignalGetObjError, "err", err)
				continue
			}
			switch pys.next {
			case enum.ProxyStatusActive:
				sigCtx, cancel := context.WithCancel(ctx)
				tun, err := service.SVC.GetTunnelById(ctx, proxy.TunnelID)
				if err != nil {
					xlog.Error(errors.SignalGetObjError, "err", err)
					continue
				}
				server, err := listener.TCPListen(sigCtx, proxy, tun)
				if err != nil {
					xlog.Error(errors.SignalListenerError, "id", pys.id, "addr", proxy.Addr(), "err", err)
				}
				sb.canceler[pys.id] = cancel
				sb.listener[pys.id] = server
				xlog.Info(xlog.SignalListenStart, "id", pys.id, "addr", proxy.Addr())
			case enum.ProxyStatusInactive:
				if cancel, ok := sb.canceler[pys.id]; ok {
					// 先调一下 cancel
					cancel()
					if ln, exist := sb.listener[pys.id]; exist {
						// 尝试着去停一下 ln, 防止泄漏
						err := ln.Close()
						if err != nil {
							// 错了就错了吧，说明 ctx 挂了一般 goroutines 也跟着挂了
							xlog.Error(errors.SignalListenerError, "id", pys.id, "addr", proxy.Addr(), "err", err)
						}
					}
				}
				xlog.Info(xlog.SignalListenStop, "id", pys.id, "addr", proxy.Addr())
			case enum.ProxyStatusSpeeding:
				if err = service.SVC.SpeedProxy(ctx, proxy); err != nil {
					_ = service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusError, err.Error())
					xlog.Error(errors.SignalSpeedTestError, "id", pys.id, "addr", proxy.Addr(), "err", err)
				}
				if err = service.SVC.UpdateProxyStatus(ctx, proxy.ID, enum.ProxyStatusActive, ""); err != nil {
					xlog.Error(errors.SignalUpdateObjError, "id", pys.id, "addr", proxy.Addr(), "err", err)
				}
			}
		case prs := <-sb.providerChannel:
			// todo: provider sync change task
			_, err := service.SVC.GetProviderById(ctx, prs.id)
			if err != nil {
				xlog.Error(errors.SignalGetObjError, "err", err)
			}
		case ts := <-sb.tunnelChannel:
			// todo: provider sync change task
			_, err := service.SVC.GetTunnelById(ctx, ts.id)
			if err != nil {
				xlog.Error(errors.SignalGetObjError, "err", err)
			}
		}
	}
}

func (sb *Bus) SendProxySignal(p uint, tp enum.ProxyStatus) {
	sb.proxyChannel <- &proxySignal{
		id:   p,
		next: tp,
	}
}

func (sb *Bus) SendProviderSignal(p uint, tp enum.ProviderStatus) {
	sb.providerChannel <- &providerSignal{
		id:   p,
		next: tp,
	}
}

func (sb *Bus) SendTunnelSignal(p uint, tp enum.TunnelStatus) {
	sb.tunnelChannel <- &tunnelSignal{
		id:   p,
		next: tp,
	}
}
