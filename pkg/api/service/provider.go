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

type provider struct {
}

func (p *provider) TotalProviders(ctx context.Context) (int64, error) {
	return dao.Q.Provider.WithContext(ctx).Count()
}

func (p *provider) ListProviders(ctx context.Context, page, size int) (models.ProviderList, error) {
	return dao.Q.Provider.WithContext(ctx).Preload(dao.Q.Provider.Tunnels.Proxies).Offset(page * size).Limit(size).Find()
}

func (p *provider) ListActiveProviders(ctx context.Context) (models.ProviderList, error) {
	return dao.Q.Provider.WithContext(ctx).Preload(dao.Q.Provider.Tunnels.Proxies).Where(
		dao.Q.Provider.Status.Eq(int8(enum.ProvStatusSuccess))).Find()
}

func (p *provider) GetProviderById(ctx context.Context, id uint) (*models.Provider, error) {
	return dao.Q.Provider.WithContext(ctx).Preload(dao.Q.Provider.Tunnels.Proxies).Where(dao.Q.Provider.ID.Eq(id)).Take()
}

func (p *provider) GetProviderByName(ctx context.Context, name string) (*models.Provider, error) {
	return dao.Q.Provider.WithContext(ctx).Preload(dao.Q.Provider.Tunnels.Proxies).Where(dao.Q.Provider.Name.Eq(name)).Take()
}

func (p *provider) CreateProvider(ctx context.Context, obj *models.Provider) (*models.Provider, error) {
	if obj.Type == nil || obj.CloudAuth == nil || len(obj.Regions) <= 0 {
		return nil, errors.New(xlog.ServiceDBNeedParamsError)
	}

	if err := dao.Q.Provider.WithContext(ctx).Create(obj); err != nil {
		return nil, err
	}
	return p.GetProviderByName(ctx, *obj.Name)
}

// UpdateProvider 用于通用式更新
func (p *provider) UpdateProvider(ctx context.Context, obj *models.Provider) (*models.Provider, error) {
	if obj.ID == 0 {
		return nil, errors.New(xlog.ServiceDBNeedParamsError)
	}

	query := dao.Q.Provider

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(obj.ID)).Updates(obj); err != nil {
		return nil, err
	}

	return p.GetProviderById(ctx, obj.ID)
}

// UpdateProviderStatus 用于更新状态，通常吞掉了状态更新时的错误
func (p *provider) UpdateProviderStatus(ctx context.Context, id uint, status enum.ProviderStatus, msg string) {
	query := dao.Q.Provider

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(&models.Provider{
		Status:        &status,
		StatusMessage: &msg,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateStatusError, "type", "provider_status", "err", err)
	}
}

func (p *provider) DeleteProvider(ctx context.Context, id uint) error {
	query := dao.Q.Provider
	res, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Delete()
	if err != nil || res.Error != nil {
		return err
	}
	return nil
}

func (p *provider) SyncProvider(ctx context.Context, prov *models.Provider) error {
	// 先同步账户
	// do auth check
	info, err := sdk.GetSDK(*prov.Type).Auth(prov.CloudAuth, prov.Regions[0])
	if err != nil {
		return err
	}

	prov.Info = info

	// 自动同步函数
	tuns, err := sdk.GetSDK(*prov.Type).SyncFC(prov.CloudAuth, prov.Regions)
	if err != nil {
		return err
	}

	for _, tun := range tuns {
		// 检测是否存在
		if SVC.ExistTunnel(ctx, nil, tun.UniqID) {
			continue
		}
		tun.ProviderId = prov.ID
		if _, err = SVC.CreateTunnel(ctx, tun.ToModel(true)); err != nil {
			return err
		}
	}

	return nil
}

func (p *provider) ExistProvider(ctx context.Context, name *string) bool {
	if name == nil {
		return false
	}
	a, err := dao.Q.Provider.WithContext(ctx).Where(dao.Q.Provider.Name.Eq(*name)).Count()
	return err == nil && a != 0
}
