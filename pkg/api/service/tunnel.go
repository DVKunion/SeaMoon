package service

import (
	"context"
	"errors"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/sdk"
)

type tunnel struct {
}

func (t *tunnel) TotalTunnels(ctx context.Context) (int64, error) {
	return dao.Q.Tunnel.WithContext(ctx).Count()
}

func (t *tunnel) ListTunnels(ctx context.Context, page, size int) (models.TunnelList, error) {
	return dao.Q.Tunnel.WithContext(ctx).Preload(dao.Q.Tunnel.Proxies).Offset(page * size).Limit(size).Find()
}

func (t *tunnel) GetTunnelById(ctx context.Context, id uint) (*models.Tunnel, error) {
	return dao.Q.Tunnel.WithContext(ctx).Preload(dao.Q.Tunnel.Proxies).Where(dao.Q.Tunnel.ID.Eq(id)).Take()
}

func (t *tunnel) GetTunnelByUId(ctx context.Context, uid string) (*models.Tunnel, error) {
	return dao.Q.Tunnel.WithContext(ctx).Preload(dao.Q.Tunnel.Proxies).Where(dao.Q.Tunnel.UniqID.Eq(uid)).Take()
}

func (t *tunnel) GetTunnelByName(ctx context.Context, name string) (*models.Tunnel, error) {
	return dao.Q.Tunnel.WithContext(ctx).Preload(dao.Q.Tunnel.Proxies).Where(dao.Q.Tunnel.Name.Eq(name)).First()
}

func (t *tunnel) CreateTunnel(ctx context.Context, obj *models.Tunnel) (*models.Tunnel, error) {
	prv, err := SVC.GetProviderById(ctx, obj.ProviderId)
	if err != nil || prv == nil || prv.ID == 0 {
		return nil, err
	}
	if *prv.MaxLimit != 0 && *prv.MaxLimit < len(prv.Tunnels)+1 {
		return nil, errors.New("limit tunnel")
	}

	// sealos 特殊检查, 防止不一致
	if *prv.Type == enum.ProvTypeSealos {
		if obj.Config.CPU < 0.1 {
			obj.Config.CPU = 0.1
		}
		if obj.Config.Memory < 64 {
			obj.Config.Memory = 64
		}
	}

	if err = dao.Q.Tunnel.WithContext(ctx).Create(obj); err != nil {
		return nil, err
	}

	if *obj.Status == enum.TunnelInitializing {
		err = t.DeployTunnel(ctx, obj.ID)
		if err != nil {
			return nil, err
		}
	}

	return t.GetTunnelByName(ctx, *obj.Name)
}

func (t *tunnel) UpdateTunnel(ctx context.Context, obj *models.Tunnel) (*models.Tunnel, error) {
	tun, err := t.GetTunnelById(ctx, obj.ID)

	if err != nil {
		return nil, err
	}

	// 更新状态为激活时，要去尝试 deploy
	if tun.Status != obj.Status {
		switch *obj.Status {
		case enum.TunnelInitializing:
			// todo: deal with error
			_ = t.DeployTunnel(ctx, tun.ID)
		case enum.TunnelInactive:
			// 再停掉远端的服务
			if err := t.StopTunnel(ctx, obj); err != nil {
				*obj.Status = enum.TunnelError
				*obj.StatusMessage = err.Error()
			}
		}
	}

	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(obj.ID)).Updates(obj); err != nil {
		return nil, err
	}

	return t.GetTunnelById(ctx, obj.ID)
}

func (t *tunnel) DeleteTunnel(ctx context.Context, id uint) error {
	tun, err := t.GetTunnelById(ctx, id)
	if err != nil {
		return err
	}

	// 先停掉本地的服务
	for _, py := range tun.Proxies {
		// 需要先停掉所有的服务
		if err = SVC.UpdateProxyStatus(ctx, py.ID, enum.ProxyStatusInactive, ""); err != nil {
			return err
		}
		// 然后删除数据
		if err = SVC.DeleteProxy(ctx, py.ID); err != nil {
			return err
		}
	}
	// 再停掉远端的服务
	if err = t.StopTunnel(ctx, tun); err != nil {
		return err
	}

	// 最后清理掉数据
	query := dao.Q.Tunnel
	res, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Delete()
	if err != nil || res.Error != nil {
		return err
	}
	return nil
}

func (t *tunnel) ExistTunnel(ctx context.Context, name, uid *string) bool {
	if uid != nil {
		a, err := dao.Q.Tunnel.WithContext(ctx).Where(dao.Q.Tunnel.UniqID.Eq(*uid)).Count()
		return err == nil && a != 0
	}
	if name != nil {
		a, err := dao.Q.Tunnel.WithContext(ctx).Where(dao.Q.Tunnel.Name.Eq(*name)).Count()
		return err == nil && a != 0
	}
	return false
}

func (t *tunnel) DeployTunnel(ctx context.Context, id uint) error {
	tun, err := t.GetTunnelById(ctx, id)

	prv, err := SVC.GetProviderById(ctx, tun.ProviderId)
	if err != nil {
		*tun.Status = enum.TunnelError
		*tun.StatusMessage = err.Error()
		_, _ = t.UpdateTunnel(ctx, tun)
		return err
	}

	addr, err := sdk.GetSDK(*prv.Type).Deploy(prv.CloudAuth, tun)
	if err != nil {
		*tun.Status = enum.TunnelError
		*tun.StatusMessage = err.Error()
		_, _ = t.UpdateTunnel(ctx, tun)
		return err
	}
	*tun.Status = enum.TunnelActive
	*tun.Addr = addr
	_, err = t.UpdateTunnel(ctx, tun)
	return err
}

func (t *tunnel) StopTunnel(ctx context.Context, tun *models.Tunnel) error {
	prov, err := SVC.GetProviderById(ctx, tun.ProviderId)
	if err != nil {
		return err
	}

	if err = sdk.GetSDK(*prov.Type).Destroy(prov.CloudAuth, tun); err != nil {
		return err
	}
	return nil
}
