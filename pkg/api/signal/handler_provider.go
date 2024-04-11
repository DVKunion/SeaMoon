package signal

import (
	"context"
	"sync"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func (sb *Bus) SendProviderSignal(p uint, tp enum.ProviderStatus) {
	sb.providerChannel <- &providerSignal{
		id:   p,
		next: tp,
		wg:   nil,
	}
}

func (sb *Bus) SendProviderSignalSync(p uint, tp enum.ProviderStatus) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	sb.providerChannel <- &providerSignal{
		id:   p,
		next: tp,
		wg:   wg,
	}
	wg.Wait()
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
	if *provider.Status == prs.next {
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
		xlog.Info(xlog.SignalSyncProvider, "id", provider.ID, "type", *provider.Type)
		service.SVC.UpdateProviderStatus(ctx, provider.ID, enum.ProvStatusSuccess, "")
	case enum.ProvStatusDelete:
		for _, tun := range provider.Tunnels {
			sb.deleteTunnel(ctx, &tun)
		}
		// 然后删除数据
		if err = service.SVC.DeleteProvider(ctx, provider.ID); err != nil {
			xlog.Error(xlog.SignalSyncProviderError, "obj", "provider", "err", err)
			service.SVC.UpdateProviderStatus(ctx, provider.ID, enum.ProvStatusSyncError, err.Error())
			return
		}
		xlog.Info(xlog.SignalDeleteProvider, "id", provider.ID, "type", provider.Type)
	}
}
