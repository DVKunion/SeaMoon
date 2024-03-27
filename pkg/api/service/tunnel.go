package service

import (
	"context"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/sdk"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
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

	return t.GetTunnelByName(ctx, *obj.Name)
}

func (t *tunnel) UpdateTunnel(ctx context.Context, obj *models.Tunnel) (*models.Tunnel, error) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(obj.ID)).Updates(obj); err != nil {
		return nil, err
	}

	return t.GetTunnelById(ctx, obj.ID)
}

func (t *tunnel) UpdateTunnelStatus(ctx context.Context, id uint, status enum.TunnelStatus, msg string) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(&models.Tunnel{
		Status:        &status,
		StatusMessage: &msg,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateStatusError, "type", "tunnel_status", "err", err)
	}
}

func (t *tunnel) UpdateTunnelAddr(ctx context.Context, id uint, addr string) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(&models.Tunnel{
		Addr: &addr,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateFiledError, "type", "tunnel_addr", "err", err)
	}
}

func (t *tunnel) DeleteTunnel(ctx context.Context, id uint) error {
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

func (t *tunnel) DeployTunnel(ctx context.Context, tun *models.Tunnel) (string, error) {

	prv, err := SVC.GetProviderById(ctx, tun.ProviderId)
	if err != nil {
		return "", err
	}

	addr, err := sdk.GetSDK(*prv.Type).Deploy(prv.CloudAuth, tun)
	if err != nil {
		return "", err
	}
	return addr, nil
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
