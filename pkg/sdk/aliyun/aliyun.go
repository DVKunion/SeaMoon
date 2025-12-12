package aliyun

import (
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
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
