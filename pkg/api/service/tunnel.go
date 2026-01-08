package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/sdk"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
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

	// 手动填充账户与密码
	if obj.Config.V2rayUid == "" {
		obj.Config.V2rayUid = tools.GenerateUUID()
	}
	if obj.Config.SSRPass == "" {
		obj.Config.SSRPass = tools.GenerateRandomString(12)
	}
	// todo: 开放认证
	if obj.Config.SSRCrypt == "" {
		obj.Config.SSRCrypt = "aes-256-gcm"
	}

	// 如果启用了级联代理，从选中的隧道获取信息
	if obj.Config.CascadeProxy && obj.Config.CascadeTunnelId > 0 {
		cascadeTunnel, err := t.GetTunnelById(ctx, obj.Config.CascadeTunnelId)
		if err != nil || cascadeTunnel == nil {
			return nil, errors.New("cascade tunnel not found")
		}
		// 填充级联代理信息
		if cascadeTunnel.Addr != nil {
			obj.Config.CascadeAddr = cascadeTunnel.GetAddr()
		}
		if cascadeTunnel.Config.V2rayUid != "" {
			obj.Config.CascadeUid = cascadeTunnel.Config.V2rayUid
		}
		if cascadeTunnel.Config.SSRPass != "" {
			obj.Config.CascadePassword = cascadeTunnel.Config.SSRPass
		}
	}

	if err = dao.Q.Tunnel.WithContext(ctx).Create(obj); err != nil {
		return nil, err
	}

	return t.GetTunnelByName(ctx, *obj.Name)
}

func (t *tunnel) UpdateTunnel(ctx context.Context, obj *models.Tunnel) (*models.Tunnel, error) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Omit(query.Status).Where(query.ID.Eq(obj.ID)).Updates(obj); err != nil {
		return nil, err
	}

	return t.GetTunnelById(ctx, obj.ID)
}

func (t *tunnel) UpdateTunnelStatus(ctx context.Context, id uint, status enum.TunnelStatus, msg string) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(models.Tunnel{
		Status:        &status,
		StatusMessage: &msg,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateStatusError, "type", "tunnel_status", "err", err)
	}
}

func (t *tunnel) UpdateTunnelStatusByUid(ctx context.Context, uid string, status enum.TunnelStatus, msg string) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.UniqID.Eq(uid)).Updates(models.Tunnel{
		Status:        &status,
		StatusMessage: &msg,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateStatusError, "type", "tunnel_status_uid", "err", err)
	}
}

func (t *tunnel) UpdateTunnelDetail(ctx context.Context, id uint, addr string, uid string) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(models.Tunnel{
		UniqID: &uid,
		Addr:   &addr,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateFiledError, "type", "tunnel_addr", "err", err)
	}
}

func (t *tunnel) DeleteTunnel(ctx context.Context, id uint) error {
	query := dao.Q.Tunnel
	res, err := query.WithContext(ctx).Unscoped().Where(query.ID.Eq(id)).Delete()
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

// GetCascadeDependents 获取依赖指定隧道的所有级联代理隧道
func (t *tunnel) GetCascadeDependents(ctx context.Context, tunnelId uint) (models.TunnelList, error) {
	var tunnels models.TunnelList
	err := dao.Q.Tunnel.WithContext(ctx).UnderlyingDB().
		Where("cascade_tunnel_id = ?", tunnelId).
		Find(&tunnels).Error
	return tunnels, err
}

func (t *tunnel) DeployTunnel(ctx context.Context, tun *models.Tunnel) (string, string, error) {

	prv, err := SVC.GetProviderById(ctx, tun.ProviderId)
	if err != nil {
		return "", "", err
	}

	return sdk.GetSDK(*prv.Type).Deploy(prv.CloudAuth, tun)
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

// HealthCheckInfo 健康检查结果
type HealthCheckInfo struct {
	OK            bool
	StartTime     string
	Version       string
	V2rayVersion  string
	ErrorMessage  string
}

// HealthCheck 对隧道进行健康检查
func (t *tunnel) HealthCheck(ctx context.Context, tun *models.Tunnel) (*HealthCheckInfo, error) {
	if tun.Addr == nil || *tun.Addr == "" {
		return nil, errors.New("tunnel address is empty")
	}

	// 构建健康检查 URL
	var healthURL string
	switch *tun.Type {
	case enum.TunnelTypeWST:
		if tun.Config.TLS {
			healthURL = fmt.Sprintf("https://%s/_health", *tun.Addr)
		} else {
			healthURL = fmt.Sprintf("http://%s/_health", *tun.Addr)
		}
	case enum.TunnelTypeGRT:
		// gRPC 也使用 HTTP 健康检查端点
		if tun.Config.TLS {
			healthURL = fmt.Sprintf("https://%s/_health", *tun.Addr)
		} else {
			healthURL = fmt.Sprintf("http://%s/_health", *tun.Addr)
		}
	default:
		return nil, errors.New("unsupported tunnel type")
	}

	xlog.Info(xlog.ServiceHealthCheck, "tunnel", *tun.Name, "url", healthURL)

	// 创建 HTTP 客户端，设置超时和跳过证书验证（用于自签名证书）
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(healthURL)
	if err != nil {
		return &HealthCheckInfo{
			OK:           false,
			ErrorMessage: err.Error(),
		}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &HealthCheckInfo{
			OK:           false,
			ErrorMessage: fmt.Sprintf("health check failed with status: %d", resp.StatusCode),
		}, errors.New(fmt.Sprintf("health check failed with status: %d", resp.StatusCode))
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &HealthCheckInfo{
			OK:           false,
			ErrorMessage: err.Error(),
		}, err
	}

	// 解析响应内容
	// 格式: OK\n2026-01-07 08:59:39\n2.0.3-c4efd4c670e7e278aa7d3a6b3fc2d78601263e27\nv2ray-core:-5.16.1
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	info := &HealthCheckInfo{OK: true}

	if len(lines) >= 1 && strings.TrimSpace(lines[0]) == "OK" {
		info.OK = true
	} else {
		info.OK = false
		info.ErrorMessage = "unexpected response format"
	}

	if len(lines) >= 2 {
		info.StartTime = strings.TrimSpace(lines[1])
	}
	if len(lines) >= 3 {
		info.Version = strings.TrimSpace(lines[2])
	}
	if len(lines) >= 4 {
		info.V2rayVersion = strings.TrimSpace(lines[3])
	}

	return info, nil
}

// UpdateTunnelHealthInfo 更新隧道的健康检查信息
func (t *tunnel) UpdateTunnelHealthInfo(ctx context.Context, id uint, version, v2rayVersion, lastCheckTime string) {
	query := dao.Q.Tunnel

	if _, err := query.WithContext(ctx).Where(query.ID.Eq(id)).Updates(models.Tunnel{
		Version:       &version,
		V2rayVersion:  &v2rayVersion,
		LastCheckTime: &lastCheckTime,
	}); err != nil {
		xlog.Error(xlog.ServiceDBUpdateFiledError, "type", "tunnel_health", "err", err)
	}
}

// CheckAndUpdateTunnelHealth 检查并更新隧道健康状态
func (t *tunnel) CheckAndUpdateTunnelHealth(ctx context.Context, tun *models.Tunnel) {
	info, err := t.HealthCheck(ctx, tun)
	checkTime := time.Now().Format("2006-01-02 15:04:05")

	if err != nil || !info.OK {
		// 健康检查失败，更新状态为异常
		errMsg := "health check failed"
		if err != nil {
			errMsg = err.Error()
		} else if info.ErrorMessage != "" {
			errMsg = info.ErrorMessage
		}
		xlog.Warn(xlog.ServiceHealthCheckFailed, "tunnel", tun.ID, "err", errMsg)
		t.UpdateTunnelStatus(ctx, tun.ID, enum.TunnelError, errMsg)
		t.UpdateTunnelHealthInfo(ctx, tun.ID, "", "", checkTime)
		return
	}

	// 健康检查成功，更新版本信息
	xlog.Info(xlog.ServiceHealthCheckSuccess, "tunnel", tun.ID, "version", info.Version)
	t.UpdateTunnelHealthInfo(ctx, tun.ID, info.Version, info.V2rayVersion, checkTime)
}
