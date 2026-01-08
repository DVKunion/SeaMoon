package signal

import (
	"context"
	"sync"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func (sb *Bus) SendTunnelSignal(p uint, tp enum.TunnelStatus) {
	sb.tunnelChannel <- &tunnelSignal{
		id:   p,
		next: tp,
		wg:   nil,
	}
}

func (sb *Bus) SendTunnelSignalSync(p uint, tp enum.TunnelStatus) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	sb.tunnelChannel <- &tunnelSignal{
		id:   p,
		next: tp,
		wg:   wg,
	}
	wg.Wait()
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
	if *tun.Status == ts.next {
		xlog.Warn(xlog.SignalMissOperationWarn, "id", ts.id, "type", "tunnel", "status", ts.next)
		return
	}
	service.SVC.UpdateTunnelStatus(ctx, tun.ID, ts.next, "")
	switch ts.next {
	case enum.TunnelInitializing, enum.TunnelActive:
		if addr, uid, err := service.SVC.DeployTunnel(ctx, tun); err != nil {
			xlog.Error(xlog.SignalDeployTunError, "obj", "tunnel", "err", err)
			service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
			return
		} else {
			service.SVC.UpdateTunnelDetail(ctx, tun.ID, addr, uid)
			// 更新 tun 的地址信息，用于健康检查
			tun.Addr = &addr
		}
		xlog.Info(xlog.SignalDeployTunnel, "id", tun.ID, "type", tun.Type)
		service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelActive, "")

		// 部署成功后进行健康检查
		go service.SVC.CheckAndUpdateTunnelHealth(ctx, tun)
	case enum.TunnelInactive:
		sb.stopTunnel(ctx, tun)
	case enum.TunnelDelete:
		sb.deleteTunnel(ctx, tun)
	}
}

func (sb *Bus) stopTunnel(ctx context.Context, tun *models.Tunnel) {
	if err := service.SVC.StopTunnel(ctx, tun); err != nil {
		xlog.Error(xlog.SignalStopTunError, "obj", "tunnel", "err", err)
		service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
		return
	}
	xlog.Info(xlog.SignalStopTunnel, "id", tun.ID, "type", tun.Type)
	return
}

func (sb *Bus) deleteTunnel(ctx context.Context, tun *models.Tunnel) {
	// 先检查是否有依赖该隧道的级联代理，如果有则先删除它们
	dependents, err := service.SVC.GetCascadeDependents(ctx, tun.ID)
	if err != nil {
		xlog.Error(xlog.SignalDeleteTunError, "obj", "tunnel", "err", err, "reason", "failed to get cascade dependents")
	}
	for _, dep := range dependents {
		xlog.Info(xlog.SignalDeleteTunnel, "id", dep.ID, "reason", "cascade dependent of tunnel", "parent", tun.ID)
		sb.deleteTunnel(ctx, dep)
	}

	// 先停掉本地的服务
	for _, py := range tun.Proxies {
		sb.deleteProxy(ctx, &py)
	}

	sb.stopTunnel(ctx, tun)

	// 最后删除服务即可
	if err := service.SVC.DeleteTunnel(ctx, tun.ID); err != nil {
		xlog.Error(xlog.SignalDeleteTunError, "obj", "tunnel", "err", err)
		service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
		return
	}
	xlog.Info(xlog.SignalDeleteTunnel, "id", tun.ID, "type", tun.Type)
}
