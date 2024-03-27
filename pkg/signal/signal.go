package signal

import (
	"context"
	"net"
	"sync"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
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
	wg   *sync.WaitGroup
}

type providerSignal struct {
	id   uint
	next enum.ProviderStatus
	wg   *sync.WaitGroup
}

type tunnelSignal struct {
	id   uint
	next enum.TunnelStatus
	wg   *sync.WaitGroup
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
			sb.proxyHandler(ctx, pys)
		case prs := <-sb.providerChannel:
			sb.providerHandler(ctx, prs)
		case ts := <-sb.tunnelChannel:
			sb.tunnelHandler(ctx, ts)
		}
	}
}

func (sb *Bus) Recover(ctx context.Context, recover string) {
	if recover == "true" {
		// 首先看一下是否需要恢复运行状态的服务
		proxies, err := service.SVC.ListActiveProxies(ctx)
		if err != nil {
			xlog.Error(xlog.SignalRecoverProxyError, "err", err)
		}
		for _, p := range proxies {
			sb.SendProxySignal(p.ID, enum.ProxyStatusRecover, nil)
		}
	}
}

func (sb *Bus) SendProxySignal(p uint, tp enum.ProxyStatus, wg *sync.WaitGroup) {
	sb.proxyChannel <- &proxySignal{
		id:   p,
		next: tp,
		wg:   wg,
	}
}

func (sb *Bus) SendProviderSignal(p uint, tp enum.ProviderStatus, wg *sync.WaitGroup) {
	sb.providerChannel <- &providerSignal{
		id:   p,
		next: tp,
		wg:   wg,
	}
}

func (sb *Bus) SendTunnelSignal(p uint, tp enum.TunnelStatus, wg *sync.WaitGroup) {
	sb.tunnelChannel <- &tunnelSignal{
		id:   p,
		next: tp,
		wg:   wg,
	}
}
