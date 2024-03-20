package service

import (
	"context"
	"errors"

	"github.com/showwin/speedtest-go/speedtest"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
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

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(obj); err != nil {
		return nil, err
	}

	return p.GetProxyById(ctx, id)
}

func (p *proxy) UpdateProxyConn(ctx context.Context, id uint, op int) error {
	query := dao.Q.Proxy

	if op == 1 {
		if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).UpdateSimple(query.Conn.Add(1)); err != nil {
			return err
		}
	}

	if op == -1 {
		if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).UpdateSimple(query.Conn.Sub(1)); err != nil {
			return err
		}
	}

	return nil
}

func (p *proxy) UpdateProxyNetFlow(ctx context.Context, id uint, in int64, out int64) error {
	query := dao.Q.Proxy

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).UpdateSimple(
		query.InBound.Add(in),
		query.OutBound.Add(out),
	); err != nil {
		return err
	}

	return nil
}

func (p *proxy) UpdateProxyStatus(ctx context.Context, id uint, status enum.ProxyStatus, msg string) error {
	query := dao.Q.Proxy

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(&models.Proxy{
		Status:        &status,
		StatusMessage: &msg,
	}); err != nil {
		return err
	}
	return nil
}

func (p *proxy) DeleteProxy(ctx context.Context, id uint) error {
	target, err := p.GetProxyById(ctx, id)
	if err != nil {
		return err
	}
	if *target.Status == enum.ProxyStatusActive {
		return errors.New("禁止删除运行中的服务，请先停止服务")
	}
	query := dao.Q.Proxy
	res, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Delete()
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
	err = s.PingTest(nil)
	if err != nil {
		return err
	}
	err = s.DownloadTest()
	if err != nil {
		return err
	}
	err = s.UploadTest()
	if err != nil {
		return err
	}

	_, err = p.UpdateProxy(ctx, obj.ID, &models.Proxy{
		SpeedUp:   &s.ULSpeed,
		SpeedDown: &s.DLSpeed,
	})

	return err
}
