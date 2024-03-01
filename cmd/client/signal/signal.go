package signal

import (
	"context"
	"log/slog"
	"net"

	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/cmd/client/listener"
	"github.com/DVKunion/SeaMoon/pkg/xlog"
)

// SignalBus 用于控制所有服务类型的启动
type SignalBus struct {
	canceler map[uint]context.CancelFunc
	listener map[uint]net.Listener

	proxyChannel chan *ProxySignal
}

type ProxySignal struct {
	ProxyId uint
	Addr    string
	Next    types.ProxyStatus
}

var signalBus = &SignalBus{
	canceler:     make(map[uint]context.CancelFunc, 0),
	listener:     make(map[uint]net.Listener, 0),
	proxyChannel: make(chan *ProxySignal, 1>>8),
}

func Signal() *SignalBus {
	return signalBus
}

// Run 监听逻辑
func (sb *SignalBus) Run(ctx context.Context) {
	for {
		select {
		case t := <-sb.proxyChannel:
			switch t.Next {
			case types.ACTIVE:
				sigCtx, cancel := context.WithCancel(ctx)
				server, err := net.Listen("tcp", t.Addr)
				if err != nil {
					slog.Error(xlog.LISTEN_ERROR, "id", t.ProxyId, "addr", t.Addr, "err", err)
				} else {
					slog.Info(xlog.LISTEN_START, "id", t.ProxyId, "addr", t.Addr)
					go listener.Listen(sigCtx, server, t.ProxyId)
					sb.canceler[t.ProxyId] = cancel
					sb.listener[t.ProxyId] = server
				}
			case types.INACTIVE:
				if cancel, ok := sb.canceler[t.ProxyId]; ok {
					// 先调一下 cancel
					cancel()
					if ln, exist := sb.listener[t.ProxyId]; exist {
						// 尝试着去停一下 ln, 防止泄漏
						err := ln.Close()
						if err != nil {
							// 错了就错了吧，说明 ctx 挂了一般 gorouting 也跟着挂了
							slog.Error(xlog.LISTEN_STOP_ERROR, "id", t.ProxyId, "addr", t.Addr, "err", err)
						}
					}
				}
				slog.Info(xlog.LISTEN_STOP, "id", t.ProxyId, "addr", t.Addr)
			}
		}
	}
}

func (sb *SignalBus) SendStartProxy(p uint, addr string) {
	sb.proxyChannel <- &ProxySignal{
		ProxyId: p,
		Addr:    addr,
		Next:    types.ACTIVE,
	}
}

func (sb *SignalBus) SendStopProxy(p uint, addr string) {
	sb.proxyChannel <- &ProxySignal{
		ProxyId: p,
		Addr:    addr,
		Next:    types.INACTIVE,
	}
}
