package sealos

import (
	"fmt"
	"strings"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/system/consts"
	"github.com/DVKunion/SeaMoon/pkg/tools"

	"github.com/DVKunion/SeaMoon/pkg/api/models"
)

type SDK struct {
}

func (s *SDK) Auth(ca *models.CloudAuth, region string) (*models.ProviderInfo, error) {
	a, c, err := getAmountAndCost(ca, region)
	if err != nil {
		return nil, err
	}
	return &models.ProviderInfo{
		Amount: &a,
		Cost:   &c,
	}, nil
}

func (s *SDK) Deploy(ca *models.CloudAuth, tun *models.Tunnel) (string, error) {

	// 拼接规则 seamoon-NAME-TYPE
	svc := "seamoon-" + *tun.Name + "-" + string(*tun.Type)
	// sealos 默认用 dockerhub 镜像
	img := "dvkunion/seamoon:" + consts.Version
	// 域名自己生成了一个随机 12 位字符串
	host := tools.GenerateRandomLetterString(12)

	addr := fmt.Sprintf("%s.%s", host, regionMap[tun.Config.Region])

	return addr, deploy(ca.KubeConfig, svc, img, host, *tun.Port, tun.Config, tun.Type)
}

func (s *SDK) Destroy(ca *models.CloudAuth, tun *models.Tunnel) error {
	// 拼接规则 seamoon-NAME-TYPE
	svcName := "seamoon-" + *tun.Name + "-" + string(*tun.Type)
	return destroy(ca.KubeConfig, svcName)
}

func (s *SDK) SyncFC(ca *models.CloudAuth, regions []string) (models.TunnelCreateApiList, error) {
	res := make(models.TunnelCreateApiList, 0)

	svcs, ingresses, err := sync(ca.KubeConfig)
	if err != nil {
		return res, err
	}

	for _, svc := range svcs.Items {
		if strings.HasPrefix(svc.Name, "seamoon-") {
			// 说明是我们的服务，继续获取对应的 label 来查找 ingress
			var tun = models.NewTunnelCreateApi()
			*tun.Name = strings.Split(svc.Name, "-")[1]
			*tun.UniqID = string(svc.ObjectMeta.UID)

			for _, condition := range svc.Status.Conditions {
				if condition.Type == "Available" && condition.Status == "True" {
					*tun.Status = enum.TunnelActive
				}
				if condition.Type == "Progressing" && condition.Status == "True" {
					*tun.Status = enum.TunnelWaiting
					*tun.StatusMessage = condition.Message
				}
				if condition.Type == "Progressing" && condition.Status == "False" {
					*tun.Status = enum.TunnelError
					*tun.StatusMessage = condition.Message
				}
				if condition.Type == "Available" && condition.Status == "False" && *tun.StatusMessage == "" {
					*tun.Status = enum.TunnelError
					*tun.StatusMessage = condition.Message
				}
			}

			*tun.Port = svc.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort
			tun.Config = &models.TunnelConfig{
				Region:     regions[0],
				CPU:        float32(svc.Spec.Template.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()) / 1000,
				Memory:     int32(svc.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().MilliValue()) / 1024 / 1024 / 1000,
				Instance:   *svc.Spec.Replicas,
				FcAuthType: enum.AuthEmpty, // sealos暂不支持认证
			}
			*tun.Type = func() enum.TunnelType {
				if strings.HasSuffix(svc.Name, "websocket") {
					return enum.TunnelTypeWST
				}
				if strings.HasSuffix(svc.Name, "grpc") {
					return enum.TunnelTypeGRT
				}
				return enum.TunnelTypeNULL
			}()
			for _, ingress := range ingresses.Items {
				if ingress.Name == svc.Name {
					*tun.Addr = ingress.Spec.Rules[0].Host
				}
			}
			res = append(res, tun)
		}
	}

	return res, nil
}
