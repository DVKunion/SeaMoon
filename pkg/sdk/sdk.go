package sdk

import (
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/sdk/aliyun"
	"github.com/DVKunion/SeaMoon/pkg/sdk/sealos"
	"github.com/DVKunion/SeaMoon/pkg/sdk/tencent"
)

type CloudSDK interface {
	// Auth 判断是否认证信息有效，并尝试增添最低的权限角色
	// 返回认证后查询的账户信息
	Auth(ca *models.CloudAuth, region string) (*models.ProviderInfo, error)
	// Deploy 部署隧道函数
	Deploy(ca *models.CloudAuth, tun *models.Tunnel) (string, string, error)
	// Destroy 删除隧道函数
	Destroy(ca *models.CloudAuth, tun *models.Tunnel) error
	// SyncFC 同步函数
	SyncFC(ca *models.CloudAuth, regions []string) (models.TunnelCreateApiList, error)

	// UpdateVersion 一键更新: 用本地的版本号请求远端服务更新至客户端版本
	// UpdateVersion(auth models.CloudAuth) error
}

var cloudFactory = map[enum.ProviderType]CloudSDK{}

func GetSDK(t enum.ProviderType) CloudSDK {
	return cloudFactory[t]
}

func init() {
	cloudFactory[enum.ProvTypeALiYun] = &aliyun.SDK{}
	cloudFactory[enum.ProvTypeTencentYun] = &tencent.SDK{}
	cloudFactory[enum.ProvTypeSealos] = &sealos.SDK{}
}
