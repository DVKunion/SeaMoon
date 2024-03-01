package tencent

import (
	"errors"
	"strconv"
	"strings"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/pkg/tools"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

var (
	// 腾讯云 在 fc 上层还有一套 namespace 的概念，为了方便管理，这里硬编码了 namespace 的内容。
	serviceName = "seamoon"
	serviceDesc = "seamoon service"
)

type SDK struct {
}

func (t *SDK) Auth(providerId uint) error {
	svc := service.GetService("provider")
	provider, ok := svc.GetById(providerId).(*models.CloudProvider)
	if !ok {
		return errors.New("can not found provider")
	}
	err := t.auth(provider)
	if err != nil {
		return err
	}
	amount, err := t.billing(provider)
	if err != nil {
		return err
	}
	provider.Amount = &amount
	svc.Update(provider.ID, provider)
	return nil
}

func (t *SDK) Deploy(providerId uint, tun *models.Tunnel) error {
	provider, ok := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	if !ok {
		return errors.New("can not found provider")
	}
	addr, err := t.deploy(provider, tun)
	if err != nil {
		*tun.Status = tunnel.ERROR
		service.GetService("tunnel").Update(tun.ID, tun)
		return err
	}
	*tun.Addr = strings.Replace(addr, "https://", "", -1)
	*tun.Status = tunnel.ERROR
	service.GetService("tunnel").Update(tun.ID, tun)
	return nil
}

func (t *SDK) Destroy(providerId uint, tun *models.Tunnel) error {
	return nil
}

func (t *SDK) SyncFC(providerId uint) error {
	svc := service.GetService("provider")
	provider, ok := svc.GetById(providerId).(*models.CloudProvider)
	if !ok {
		return errors.New("can not found provider")
	}
	fcList, err := t.sync(provider)
	if err != nil {
		return err
	}
	for _, fc := range fcList {
		fcNameList := strings.Split(*fc.detail.FunctionName, "-")
		fcName := fcNameList[0]
		if len(fcNameList) > 3 {
			fcName = fcNameList[2]
		}
		var tun = models.Tunnel{
			CloudProviderId: provider.ID,
			Name:            &fcName,
			TunnelConfig: &models.TunnelConfig{
				CPU:      0,
				Memory:   tools.PtrInt32(fc.detail.MemorySize),
				Instance: 1, // 这个玩意tmd怎么也找不到，同步过来的就算他1好了。

				TLS: true, // 默认同步过来都打开
				Tor: func() bool {
					if fc.detail.Environment != nil {
						// 如果是 开启 Tor 的隧道，需要有环境变量
						return len(fc.detail.Environment.Variables) > 0
					}
					return false
				}(),
			},
		}
		// 自动填充防止空指针
		models.AutoFull(&tun)
		*tun.Type = tunnel.TranType(*fc.detail.Description)
		*tun.Port = strconv.Itoa(int(*fc.detail.ImageConfig.ImagePort))
		*tun.Status = tunnel.ACTIVE
		*tun.Addr = fc.addr
		tun.TunnelConfig.TunnelAuthType = types.TransAuthType(fc.auth)

		service.GetService("tunnel").Create(&tun)
	}
	return nil
}
