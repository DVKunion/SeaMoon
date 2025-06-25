package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/sdk/aliyun/api"
	apimodels "github.com/DVKunion/SeaMoon/pkg/sdk/aliyun/models"
	"github.com/DVKunion/SeaMoon/pkg/system/version"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"

	"github.com/alibabacloud-go/tea/tea"
)

func getBilling(ca *models.CloudAuth) (float64, error) {
	client, err := api.NewClient(ca.AccessKey, ca.AccessSecret, "business.aliyuncs.com")
	if err != nil {
		return 0, err
	}
	var r = &apimodels.BillingResponse{}
	if err = api.Call(client, api.NewGetBillingsParams(), nil, r); err != nil {
		return 0, err
	}
	if r.StatusCode != 200 || r.Body.Code != "200" {
		return 0, errors.New(r.Body.Message)
	}
	return strconv.ParseFloat(strings.Replace(r.Body.Data["AvailableAmount"].(string), ",", "", -1), 64)
}

func deploy(ca *models.CloudAuth, tun *models.Tunnel) (string, string, error) {
	uid := ""
	// 原生的库是真tm的难用，
	client, err := api.NewClient(ca.AccessKey, ca.AccessSecret, fmt.Sprintf("%s.%s.fc.aliyuncs.com", ca.AccessId, tun.Config.Region))
	if err != nil {
		return "", "", err
	}
	body := map[string]interface{}{
		"serviceName": serviceName,
		"desc":        serviceDesc,
	}
	if err := api.Call(client, api.NewCreateServiceParams(), body, nil); err != nil {
		var e *tea.SDKError
		if errors.As(err, &e) && *e.Code != "ServiceAlreadyExists" {
			return "", "", err
		}
	}
	funcName := *tun.Name
	// 有了服务了，现在来创建函数
	body = map[string]interface{}{
		"functionName":        funcName,
		"description":         string(*tun.Type),
		"runtime":             "custom-container",
		"handler":             "main",
		"timeout":             300,
		"diskSize":            512,
		"cpu":                 tun.Config.CPU,
		"memorySize":          tun.Config.Memory,
		"caPort":              *tun.Port,
		"instanceConcurrency": tun.Config.Instance,
		"instanceType":        "e1",
		"environmentVariables": map[string]*string{
			"SM_SS_PASS":  tea.String(tun.Config.SSRPass),
			"SM_SS_CRYPT": tea.String(tun.Config.SSRCrypt),
			"SM_UID":      tea.String(tun.Config.V2rayUid),
		},
		"customContainerConfig": map[string]*string{
			"image": tea.String(fmt.Sprintf("%s:%s", registryEndPoint[tun.Config.Region], version.Version)),
			"args": tea.String(func() string {
				switch *tun.Type {
				case enum.TunnelTypeWST:
					return "[\"server\", \"-p\", \"9000\", \"-t\", \"websocket\"]"
				case enum.TunnelTypeGRT:
					return "[\"server\", \"-p\", \"8089\", \"-t\", \"grpc\"]"
				}
				return ""
			}()),
		},
	}
	resp := &apimodels.FunctionCreateResponse{}
	if err := api.Call(client, api.NewCreateFCParams(), body, resp); err != nil {
		return "", "", err
	} else {
		uid = resp.Body.FunctionId
	}
	conf := apimodels.TriggerConfig{
		Methods:            []string{"GET", "POST"},
		AuthType:           "anonymous",
		DisableURLInternet: false,
	}
	bytes, _ := json.Marshal(&conf)
	body = map[string]interface{}{
		"triggerName":   string(*tun.Type),
		"triggerType":   "http",
		"triggerConfig": string(bytes),
	}
	respT := &apimodels.TriggerCreateResponse{}
	if err = api.Call(client, api.NewCreateTriggerParams(funcName), body, respT); err != nil {
		return "", "", err
	}

	return strings.Replace(respT.Body.UrlInternet, "https://", "", -1), uid, nil
}

func destroy(ca *models.CloudAuth, tun *models.Tunnel) error {
	client, err := api.NewClient(ca.AccessKey, ca.AccessSecret, fmt.Sprintf("%s.%s.fc.aliyuncs.com", ca.AccessId, tun.Config.Region))
	if err != nil {
		return err
	}
	// 先删除 trigger
	var respT = &apimodels.TriggerListResponse{}
	if err = api.Call(client, api.NewListTriggerParams(*tun.Name), nil, respT); err != nil {
		return err
	}

	for _, t := range respT.Body.Triggers {
		if err = api.Call(client, api.NewDeleteTriggerParams(*tun.Name, t.TriggerName), nil, nil); err != nil {
			return err
		}
	}
	if err = api.Call(client, api.NewDeleteFCParams(*tun.Name), nil, nil); err != nil {
		return err
	}
	return nil
}

func sync(ca *models.CloudAuth, regions []string) (models.TunnelCreateApiList, error) {
	res := make(models.TunnelCreateApiList, 0)

	for _, rg := range regions {
		client, err := api.NewClient(ca.AccessKey, ca.AccessSecret, fmt.Sprintf("%s.%s.fc.aliyuncs.com", ca.AccessId, rg))
		if err != nil {
			return res, err
		}
		var r = &apimodels.FunctionListResponse{}
		if err := api.Call(client, api.NewListFCParams(), nil, r); err != nil {
			var e *tea.SDKError
			if errors.As(err, &e) && *e.Code == "ServiceNotFound" {
				// 说明没有 service，继续下一个地区好了
				continue
			} else {
				// todo: 处理这个 非 200
				return res, err
			}
		}

		for _, f := range r.Body.Functions {
			if f.Description == "" {
				xlog.Warn("sdk", "发现了不正确的隧道", "fc_name", f.FunctionName, "fc_type", f.Description)
				continue
			}
			var tun = models.NewTunnelCreateApi()
			tun.UniqID = tea.String(f.FunctionId)
			tun.Name = tea.String(f.FunctionName)
			tun.Config = &models.TunnelConfig{
				CPU:      f.Cpu,
				Region:   rg,
				Memory:   f.MemorySize,
				Instance: f.InstanceConcurrency,
				TLS:      true, // 默认同步过来都打开
				Tor:      false,
			}
			if len(f.EnvironmentVariables) > 0 {
				for key, value := range f.EnvironmentVariables {
					if key == "SEAMOON_TOR" {
						tun.Config.Tor = true
					}
					if key == "SM_SS_CRYPT" {
						tun.Config.SSRCrypt = value
					}
					if key == "SM_SS_PASS" {
						tun.Config.SSRPass = value
					}
					if key == "SM_UID" {
						tun.Config.V2rayUid = value
					}
				}
			}
			*tun.Type = enum.TransTunnelType(f.Description)
			tun.Port = tea.Int32(f.CaPort)
			var respT = &apimodels.TriggerListResponse{}
			if err := api.Call(client, api.NewListTriggerParams(f.FunctionName), nil, respT); err == nil {
				for _, t := range respT.Body.Triggers {
					if t.TriggerType == "http" {
						*tun.Addr = strings.Replace(t.UrlInternet, "https://", "", -1)
						tun.Config.FcAuthType = enum.TransAuthType(t.TriggerConfig.AuthType)
					}
				}
				*tun.Status = enum.TunnelActive
			} else {
				*tun.Status = enum.TunnelError
				*tun.StatusMessage = err.Error()
			}
			res = append(res, tun)
		}
	}

	return res, nil
}
