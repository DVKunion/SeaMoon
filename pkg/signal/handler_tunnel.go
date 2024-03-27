package signal

import (
	"context"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

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
		_ = sb.stopTunnel(ctx, tun)
	case enum.TunnelDelete:
		sb.deleteTunnel(ctx, tun)
	}
}

func (sb *Bus) stopTunnel(ctx context.Context, tun *models.Tunnel) error {
	if err := service.SVC.StopTunnel(ctx, tun); err != nil {
		xlog.Error(xlog.SignalStopTunError, "obj", "tunnel", "err", err)
		service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
		return err
	}
	xlog.Info(xlog.SignalStopTunnel, "id", tun.ID, "type", tun.Type)
	return nil
}

func (sb *Bus) deleteTunnel(ctx context.Context, tun *models.Tunnel) {
	// 先停掉本地的服务
	for _, py := range tun.Proxies {
		sb.deleteProxy(ctx, &py)
	}
	if err := sb.stopTunnel(ctx, tun); err != nil {
		xlog.Error(xlog.SignalDeleteTunError, "obj", "tunnel", "err", err)
		service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
		return
	}
	// 最后删除服务即可
	if err := service.SVC.DeleteTunnel(ctx, tun.ID); err != nil {
		xlog.Error(xlog.SignalDeleteTunError, "obj", "tunnel", "err", err)
		service.SVC.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, err.Error())
		return
	}
	xlog.Info(xlog.SignalDeleteTunnel, "id", tun.ID, "type", tun.Type)
}
