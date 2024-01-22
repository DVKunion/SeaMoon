package sdk

import (
	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
)

type CloudSDK interface {
	// Auth 判断是否存在权限动作, 并获取账户钱包信息
	Auth(providerId uint) error
	// Deploy 部署隧道函数
	Deploy(providerId uint, tun *models.Tunnel) error
	// Destroy 删除隧道函数
	Destroy(providerId uint, tun *models.Tunnel) error
	// SyncFC 同步函数
	SyncFC(providerId uint) error

	// Billing(provider models.CloudProvider, tunnel models.Tunnel) error

	// UpdateVersion 一键更新: 用本地的版本号请求远端服务更新至
	// UpdateVersion(auth models.CloudAuth) error
}

var cloudFactory = map[types.CloudType]CloudSDK{}

func GetSDK(t types.CloudType) CloudSDK {
	return cloudFactory[t]
}

func init() {
	cloudFactory[types.ALiYun] = &ALiYunSDK{}
	cloudFactory[types.Sealos] = &SealosSDK{}
}
