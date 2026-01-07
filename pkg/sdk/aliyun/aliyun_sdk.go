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
	fc "github.com/alibabacloud-go/fc-20230330/v4/client"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/alibabacloud-go/tea/tea"
)

func getBilling(ca *models.CloudAuth) (float64, error) {
	client, err := api.NewBillingClient(ca.AccessKey, ca.AccessSecret, "business.aliyuncs.com")
	if err != nil {
		return 0, err
	}
	var r = &apimodels.BillingResponse{}
	if err = api.CallBilling(client, api.NewGetBillingsParams(), nil, r); err != nil {
		return 0, err
	}
	if r.StatusCode != 200 || r.Body.Code != "200" {
		return 0, errors.New(r.Body.Message)
	}
	return strconv.ParseFloat(strings.Replace(r.Body.Data["AvailableAmount"].(string), ",", "", -1), 64)
}

func deploy(ca *models.CloudAuth, tun *models.Tunnel) (string, string, error) {
	uid := ""
	// 使用新的 FC SDK
	client, err := api.NewFCClient(ca.AccessKey, ca.AccessSecret, fmt.Sprintf("fcv3.%s.aliyuncs.com", tun.Config.Region))
	if err != nil {
		return "", "", err
	}

	funcName := *tun.Name

	// 创建函数
	createFuncReq := &fc.CreateFunctionRequest{}

	// 构建 Command 参数
	var commandArgs []*string
	switch *tun.Type {
	case enum.TunnelTypeWST:
		commandArgs = []*string{
			tea.String("server"),
			tea.String("-p"),
			tea.String("9000"),
			tea.String("-t"),
			tea.String("websocket"),
		}
	case enum.TunnelTypeGRT:
		commandArgs = []*string{
			tea.String("server"),
			tea.String("-p"),
			tea.String("8089"),
			tea.String("-t"),
			tea.String("grpc"),
		}
	}

	envVars := map[string]*string{
		"SM_SS_PASS":  tea.String(tun.Config.SSRPass),
		"SM_SS_CRYPT": tea.String(tun.Config.SSRCrypt),
		"SM_UID":      tea.String(tun.Config.V2rayUid),
	}

	// 如果启用了级联代理，添加级联代理环境变量
	if tun.Config.CascadeProxy && tun.Config.CascadeAddr != "" && tun.Config.CascadeUid != "" {
		envVars["SM_CASCADE_ADDR"] = tea.String(tun.Config.CascadeAddr)
		envVars["SM_CASCADE_UID"] = tea.String(tun.Config.CascadeUid)
		if tun.Config.CascadePassword != "" {
			envVars["SM_CASCADE_PASS"] = tea.String(tun.Config.CascadePassword)
		}
	}

	createFuncInput := &fc.CreateFunctionInput{
		FunctionName:        tea.String(funcName),
		Description:         tea.String(string(*tun.Type)),
		Runtime:             tea.String("custom-container"),
		Handler:             tea.String("main"),
		Timeout:             tea.Int32(300),
		DiskSize:            tea.Int32(512),
		Cpu:                 tea.Float32(tun.Config.CPU),
		MemorySize:          tea.Int32(tun.Config.Memory),
		InstanceConcurrency: tea.Int32(tun.Config.Instance),
		EnvironmentVariables: envVars,
		CustomContainerConfig: &fc.CustomContainerConfig{
			Image:   tea.String(fmt.Sprintf("%s:%s", registryEndPoint[tun.Config.Region], version.Version)),
			Port:    tea.Int32(*tun.Port),
			Command: commandArgs,
		},
	}
	createFuncReq.Body = createFuncInput

	createFuncResp, err := client.CreateFunction(createFuncReq)
	if err != nil {
		// 检查是否是镜像相关的错误
		var sdkErr *dara.SDKError
		if errors.As(err, &sdkErr) {
			errMsg := dara.StringValue(sdkErr.Message)
			// 如果是镜像不存在或 ACR 相关的错误，提供更清晰的错误信息
			if strings.Contains(errMsg, "Image not stored in ACR") ||
				strings.Contains(errMsg, "IMAGE_NOT_EXIST") ||
				strings.Contains(errMsg, "repo image is not exist") {
				imageName := fmt.Sprintf("%s:%s", registryEndPoint[tun.Config.Region], version.Version)
				return "", "", fmt.Errorf("镜像不存在于 ACR 仓库中，请确保镜像已推送到 ACR: %s. 错误详情: %s", imageName, errMsg)
			}
		}
		return "", "", err
	}

	if createFuncResp.Body != nil && createFuncResp.Body.FunctionId != nil {
		uid = dara.StringValue(createFuncResp.Body.FunctionId)
	}

	// 创建触发器
	conf := apimodels.TriggerConfig{
		Methods:            []string{"GET", "POST"},
		AuthType:           "anonymous",
		DisableURLInternet: false,
	}
	bytes, _ := json.Marshal(&conf)

	createTriggerReq := &fc.CreateTriggerRequest{}
	createTriggerInput := &fc.CreateTriggerInput{
		TriggerName:   tea.String(string(*tun.Type)),
		TriggerType:   tea.String("http"),
		TriggerConfig: tea.String(string(bytes)),
	}
	createTriggerReq.Body = createTriggerInput

	createTriggerResp, err := client.CreateTrigger(tea.String(funcName), createTriggerReq)
	if err != nil {
		return "", "", err
	}

	var urlInternet string
	if createTriggerResp.Body != nil && createTriggerResp.Body.HttpTrigger != nil && createTriggerResp.Body.HttpTrigger.UrlInternet != nil {
		urlInternet = dara.StringValue(createTriggerResp.Body.HttpTrigger.UrlInternet)
	}

	return strings.Replace(urlInternet, "https://", "", -1), uid, nil
}

func destroy(ca *models.CloudAuth, tun *models.Tunnel) error {
	client, err := api.NewFCClient(ca.AccessKey, ca.AccessSecret, fmt.Sprintf("fcv3.%s.aliyuncs.com", tun.Config.Region))
	if err != nil {
		return err
	}

	funcName := *tun.Name

	// 先列出所有触发器
	listTriggersReq := &fc.ListTriggersRequest{}
	listTriggersResp, err := client.ListTriggers(tea.String(funcName), listTriggersReq)
	if err != nil {
		return err
	}

	// 删除所有触发器
	if listTriggersResp.Body != nil && listTriggersResp.Body.Triggers != nil {
		for _, t := range listTriggersResp.Body.Triggers {
			if t.TriggerName != nil {
				_, err = client.DeleteTrigger(tea.String(funcName), t.TriggerName)
				if err != nil {
					return err
				}
			}
		}
	}

	// 删除函数
	_, err = client.DeleteFunction(tea.String(funcName))
	if err != nil {
		return err
	}

	return nil
}

func sync(ca *models.CloudAuth, regions []string) (models.TunnelCreateApiList, error) {
	res := make(models.TunnelCreateApiList, 0)

	for _, rg := range regions {
		client, err := api.NewFCClient(ca.AccessKey, ca.AccessSecret, fmt.Sprintf("fcv3.%s.aliyuncs.com", rg))
		if err != nil {
			return res, err
		}

		// 列出所有函数
		listFuncReq := &fc.ListFunctionsRequest{}
		listFuncResp, err := client.ListFunctions(listFuncReq)
		if err != nil {
			return res, err
		}

		if listFuncResp.Body == nil || listFuncResp.Body.Functions == nil {
			continue
		}

		for _, f := range listFuncResp.Body.Functions {
			if f.Description == nil || dara.StringValue(f.Description) == "" {
				xlog.Warn("sdk", "发现了不正确的隧道", "fc_name", dara.StringValue(f.FunctionName), "fc_type", dara.StringValue(f.Description))
				continue
			}

			var tun = models.NewTunnelCreateApi()
			if f.FunctionId != nil {
				tun.UniqID = tea.String(dara.StringValue(f.FunctionId))
			}
			if f.FunctionName != nil {
				tun.Name = tea.String(dara.StringValue(f.FunctionName))
			}

			var cpu float32
			if f.Cpu != nil {
				cpu = dara.Float32Value(f.Cpu)
			}
			var memory int32
			if f.MemorySize != nil {
				memory = dara.Int32Value(f.MemorySize)
			}
			var instance int32
			if f.InstanceConcurrency != nil {
				instance = dara.Int32Value(f.InstanceConcurrency)
			}

			tun.Config = &models.TunnelConfig{
				CPU:      cpu,
				Region:   rg,
				Memory:   memory,
				Instance: instance,
				TLS:      true, // 默认同步过来都打开
				Tor:      false,
			}

			if f.EnvironmentVariables != nil && len(f.EnvironmentVariables) > 0 {
				for key, value := range f.EnvironmentVariables {
					if key == "SEAMOON_TOR" {
						tun.Config.Tor = true
					}
					if key == "SM_SS_CRYPT" && value != nil {
						tun.Config.SSRCrypt = dara.StringValue(value)
					}
					if key == "SM_SS_PASS" && value != nil {
						tun.Config.SSRPass = dara.StringValue(value)
					}
					if key == "SM_UID" && value != nil {
						tun.Config.V2rayUid = dara.StringValue(value)
					}
				}
			}

			if f.Description != nil {
				*tun.Type = enum.TransTunnelType(dara.StringValue(f.Description))
			}
			// 从 CustomContainerConfig 获取 Port
			if f.CustomContainerConfig != nil && f.CustomContainerConfig.Port != nil {
				tun.Port = tea.Int32(dara.Int32Value(f.CustomContainerConfig.Port))
			}

			// 列出触发器
			listTriggersReq := &fc.ListTriggersRequest{}
			listTriggersResp, err := client.ListTriggers(f.FunctionName, listTriggersReq)
			if err == nil && listTriggersResp.Body != nil && listTriggersResp.Body.Triggers != nil {
				for _, t := range listTriggersResp.Body.Triggers {
					if t.TriggerType != nil && dara.StringValue(t.TriggerType) == "http" {
						if t.HttpTrigger != nil && t.HttpTrigger.UrlInternet != nil {
							*tun.Addr = strings.Replace(dara.StringValue(t.HttpTrigger.UrlInternet), "https://", "", -1)
						}
						// 解析 triggerConfig 获取 AuthType
						if t.TriggerConfig != nil {
							var triggerConfig apimodels.TriggerConfig
							if err := json.Unmarshal([]byte(dara.StringValue(t.TriggerConfig)), &triggerConfig); err == nil {
								tun.Config.FcAuthType = enum.TransAuthType(triggerConfig.AuthType)
							}
						}
					}
				}
				*tun.Status = enum.TunnelActive
			} else {
				*tun.Status = enum.TunnelError
				if err != nil {
					*tun.StatusMessage = err.Error()
				}
			}

			res = append(res, tun)
		}
	}

	return res, nil
}
