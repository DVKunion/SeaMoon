package aliyun

import (
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

var (
	// 阿里云 在 fc 上层还有一套 service 的概念，为了方便管理，这里硬编码了 service 的内容。
	serviceName = "seamoon"
	serviceDesc = "seamoon service"
)

// SDK FC
type SDK struct {
}

func (a *SDK) Auth(ca *models.CloudAuth, region string) (*models.ProviderInfo, error) {
	amount, err := getBilling(ca)
	if err != nil {
		return nil, err
	}

	return &models.ProviderInfo{
		Amount: &amount,
		Cost:   tools.Float64Ptr(0),
	}, nil
}

func (a *SDK) Deploy(ca *models.CloudAuth, tun *models.Tunnel) (string, string, error) {
	return deploy(ca, tun)
}

func (a *SDK) Destroy(ca *models.CloudAuth, tun *models.Tunnel) error {
	return destroy(ca, tun)
}

func (a *SDK) SyncFC(ca *models.CloudAuth, regions []string) (models.TunnelCreateApiList, error) {
	return sync(ca, regions)
}
