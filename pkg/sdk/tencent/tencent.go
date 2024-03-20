package tencent

import (
	"strings"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/tools"
)

var (
	// 腾讯云 在 fc 上层还有一套 namespace 的概念，为了方便管理，这里硬编码了 namespace 的内容。
	serviceName = "seamoon"
	serviceDesc = "seamoon service"
)

type SDK struct {
}

func (t *SDK) Auth(ca *models.CloudAuth, region string) (*models.ProviderInfo, error) {
	// 先创建权限与角色
	err := createRole(ca)
	if err != nil {
		return nil, err
	}
	// 查询账户余额
	amount, err := getAmount(ca)
	if err != nil {
		return nil, err
	}
	return &models.ProviderInfo{
		Amount: &amount,
		Cost:   tools.Float64Ptr(0),
	}, nil
}

func (t *SDK) Deploy(ca *models.CloudAuth, tun *models.Tunnel) (string, error) {
	addr, err := deploy(ca, tun)
	if err != nil {
		return "", err
	}
	return strings.Replace(addr, "https://", "", -1), nil
}

func (t *SDK) Destroy(ca *models.CloudAuth, tun *models.Tunnel) error {
	return destroy(ca, tun)
}

func (t *SDK) SyncFC(ca *models.CloudAuth, regions []string) (models.TunnelCreateApiList, error) {
	res := make(models.TunnelCreateApiList, 0)
	fcList, err := sync(ca, regions)
	if err != nil {
		return res, err
	}
	for _, fc := range fcList {
		fcNameList := strings.Split(*fc.detail.FunctionName, "-")
		fcName := fcNameList[0]
		if len(fcNameList) >= 2 {
			fcName = fcNameList[1]
		}
		var tun = models.NewTunnelCreateApi()
		tun.Name = &fcName
		tun.UniqID = fc.detail.FunctionId
		tun.Config = &models.TunnelConfig{
			Region:   fc.region,
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
		}
		*tun.Type = enum.TransTunnelType(*fc.detail.Description)
		*tun.Port = int32(*fc.detail.ImageConfig.ImagePort)
		*tun.Status = enum.TunnelActive
		*tun.Addr = fc.addr
		tun.Config.FcAuthType = enum.TransAuthType(fc.auth)
		res = append(res, tun)
	}
	return res, nil
}
