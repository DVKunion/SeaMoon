package service

import (
	"context"

	"github.com/showwin/speedtest-go/speedtest"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

type proxy struct {
}

func (p *proxy) TotalProxies(ctx context.Context) (int64, error) {
	return dao.Q.Proxy.WithContext(ctx).Count()
}

func (p *proxy) ListProxies(ctx context.Context, page, size int) (models.ProxyList, error) {
	return dao.Q.Proxy.WithContext(ctx).Offset(page * size).Limit(size).Find()
}

func (p *proxy) ListActiveProxies(ctx context.Context) (models.ProxyList, error) {
	return dao.Q.Proxy.WithContext(ctx).Where(
		dao.Q.Proxy.Status.Eq(int8(enum.ProxyStatusActive))).Find()
}

func (p *proxy) GetProxyById(ctx context.Context, id uint) (*models.Proxy, error) {
	return dao.Q.Proxy.WithContext(ctx).Where(dao.Q.Proxy.ID.Eq(id)).Take()
}

func (p *proxy) GetProxyByName(ctx context.Context, name string) (*models.Proxy, error) {
	return dao.Q.Proxy.WithContext(ctx).Where(dao.Q.Proxy.Name.Eq(name)).First()
}

func (p *proxy) CreateProxy(ctx context.Context, obj *models.Proxy) (*models.Proxy, error) {
	if err := dao.Q.Proxy.WithContext(ctx).Create(obj); err != nil {
		return nil, err
	}
	return p.GetProxyByName(ctx, *obj.Name)
}

func (p *proxy) UpdateProxy(ctx context.Context, id uint, obj *models.Proxy) (*models.Proxy, error) {
	query := dao.Q.Proxy

	if _, err := query.WithContext(ctx).Omit(query.Status).Where(query.ID.Eq(id)).Updates(obj); err != nil {
		return nil, err
	}

	return p.GetProxyById(ctx, id)
}

func (p *proxy) UpdateProxyConn(ctx context.Context, id uint, op int) {
	query := dao.Q.Proxy

	switch op {
	case 1:
		if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).UpdateSimple(query.Conn.Add(1)); err != nil {
			xlog.Error(xlog.ServiceDBUpdateFiledError, "type", "proxy_conn", "err", err)
		}
	case -1:
		if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).UpdateSimple(query.Conn.Sub(1)); err != nil {
			xlog.Error(xlog.ServiceDBUpdateFiledError, "type", "proxy_conn", "err", err)
		}
	}
}

func (p *proxy) UpdateProxyNetworkInfo(ctx context.Context, id uint, in int64, out int64) {
	query := dao.Q.Proxy

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).UpdateSimple(
		query.InBound.Add(in),
		query.OutBound.Add(out),
	); err != nil {
		xlog.Error(xlog.ServiceDBUpdateFiledError, "type", "proxy_network", "err", err)
	}
}

func (p *proxy) UpdateProxyNetworkLag(ctx context.Context, id uint, lag int64) {
	query := dao.Q.Proxy

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Update(
		query.Lag, lag,
	); err != nil {
		xlog.Error(xlog.ServiceDBUpdateFiledError, "type", "proxy_lag", "err", err)
	}
}

func (p *proxy) UpdateProxyStatus(ctx context.Context, id uint, status enum.ProxyStatus, msg string) {
	query := dao.Q.Proxy

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(models.Proxy{
		Status:        &status,
		StatusMessage: &msg,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateStatusError, "type", "proxy_status", "err", err)
	}
}

func (p *proxy) DeleteProxy(ctx context.Context, id uint) error {
	query := dao.Q.Proxy
	res, err := query.WithContext(ctx).Unscoped().Where(query.ID.Eq(id)).Delete()
	if err != nil || res.Error != nil {
		return err
	}
	return nil
}

func (p *proxy) ExistProxy(ctx context.Context, name *string) bool {
	if name == nil {
		return false
	}
	a, err := dao.Q.Proxy.WithContext(ctx).Where(dao.Q.Proxy.Name.Eq(*name)).Count()
	return err == nil && a != 0
}

func (p *proxy) SpeedProxy(ctx context.Context, obj *models.Proxy) error {
	sClient := speedtest.New()
	speedtest.WithUserConfig(&speedtest.UserConfig{Proxy: obj.ProtoAddr()})(sClient)

	serverList, err := sClient.FetchServers()
	if err != nil {
		return err
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		return err
	}
	if len(targets) < 1 {
		return err
	}
	s := targets[0]
	if err = s.PingTest(nil); err != nil {
		return err
	}
	if err = s.DownloadTest(); err != nil {
		return err
	}
	if err = s.UploadTest(); err != nil {
		return err
	}

	// 将 speedtest.ByteRate 转换为 float64
	ulSpeed := float64(s.ULSpeed)
	dlSpeed := float64(s.DLSpeed)

	_, err = p.UpdateProxy(ctx, obj.ID, &models.Proxy{
		SpeedUp:   &ulSpeed,
		SpeedDown: &dlSpeed,
	})

	return err
}
