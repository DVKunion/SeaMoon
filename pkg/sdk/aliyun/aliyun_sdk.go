package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	bss "github.com/alibabacloud-go/bssopenapi-20171214/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/fc-go-sdk"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func getBilling(ca *models.CloudAuth) (float64, error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: &ca.AccessKey,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: &ca.AccessSecret,
	}

	config.Endpoint = tea.String("business.aliyuncs.com")
	client, err := bss.NewClient(config)
	if err != nil {
		return 0, err
	}
	params := &openapi.Params{
		// 接口名称
		Action: tea.String("QueryAccountBalance"),
		// 接口版本
		Version: tea.String("2017-12-14"),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String("POST"),
		AuthType: tea.String("AK"),
		Style:    tea.String("RPC"),
		// 接口 PATH
		Pathname: tea.String("/"),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	response, err := client.CallApi(params, request, runtime)
	if err != nil {
		return 0, err
	}
	bs, err := json.Marshal(response)
	if err != nil {
		return 0, err
	}
	var r Resp
	err = json.Unmarshal(bs, &r)
	if err != nil {
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
	client, err := fc.NewClient(
		fmt.Sprintf("%s.%s.fc.aliyuncs.com", ca.AccessId, tun.Config.Region),
		"2016-08-15", ca.AccessKey, ca.AccessSecret)
	if err != nil {
		return "", "", err
	}
	// 先尝试是否已经存在了 svc
	_, err = client.GetService(fc.NewGetServiceInput(serviceName))
	if err != nil {
		if fcErr, ok := err.(*fc.ServiceError); ok {
			if fcErr.HTTPStatus == http.StatusNotFound {
				// 说明 service 空了，先创建svc
				_, err := client.CreateService(fc.NewCreateServiceInput().
					WithServiceName(serviceName).
					WithDescription(serviceDesc))
				if err != nil {
					return "", "", err
				}
			}
		} else {
			return "", "", err
		}
	}

	funcName := *tun.Name
	// 有了服务了，现在来创建函数
	if res, err := client.CreateFunction(fc.NewCreateFunctionInput(serviceName).
		WithFunctionName(funcName).
		WithDescription(string(*tun.Type)).
		WithRuntime("custom-container").
		WithCPU(tun.Config.CPU).
		WithMemorySize(tun.Config.Memory).
		WithHandler("main").
		WithEnvironmentVariables(map[string]string{
			"SM_SS_PASS":  tun.Config.SSRPass,
			"SM_SS_CRYPT": tun.Config.SSRCrypt,
			"SM_UID":      tun.Config.V2rayUid,
		}).
		WithDisk(512).
		WithInstanceConcurrency(tun.Config.Instance).
		WithCAPort(*tun.Port).
		WithInstanceType("e1"). // 性能实例
		WithTimeout(300).
		WithCustomContainerConfig(fc.NewCustomContainerConfig().
			WithImage(fmt.Sprintf("%s:%s", registryEndPoint[tun.Config.Region], xlog.Version)).
			WithArgs(func() string {
				switch *tun.Type {
				case enum.TunnelTypeWST:
					return "[\"server\", \"-p\", \"9000\", \"-t\", \"websocket\"]"
				case enum.TunnelTypeGRT:
					return "[\"server\", \"-p\", \"8089\", \"-t\", \"grpc\"]"
				}
				return ""
			}()))); err != nil {
		return "", "", err
	} else {
		uid = *res.FunctionID
	}
	// 有了函数了，接下来创建 trigger
	if _, err = client.CreateTrigger(fc.NewCreateTriggerInput(serviceName, funcName).
		WithTriggerType("http").
		WithTriggerName(string(*tun.Type)).
		WithTriggerConfig(fc.TriggerConfig{
			Methods:            []string{"GET", "POST"},
			AuthType:           "anonymous",
			DisableURLInternet: false,
		})); err != nil {
		return "", "", err
	}
	// 创建成功了, 查一下
	respTS, err := client.GetTrigger(fc.NewGetTriggerInput(serviceName, funcName, string(*tun.Type)))
	if err != nil {
		return "", "", err
	}

	return strings.Replace(respTS.UrlInternet, "https://", "", -1), uid, nil
}

func destroy(ca *models.CloudAuth, tun *models.Tunnel) error {
	client, err := fc.NewClient(
		fmt.Sprintf("%s.%s.fc.aliyuncs.com", ca.AccessId, tun.Config.Region),
		"2016-08-15", ca.AccessKey, ca.AccessSecret)
	if err != nil {
		return err
	}
	// 先删除 trigger
	if _, err = client.DeleteTrigger(
		fc.NewDeleteTriggerInput(serviceName, *tun.Name, string(*tun.Type))); err != nil {
		return err
	}
	// 在删除 fc
	if _, err = client.DeleteFunction(
		fc.NewDeleteFunctionInput(serviceName, *tun.Name)); err != nil {
		return err
	}
	// 不要删除 service, service 又不花钱, 删了还得重新创建，还可能导致整个 svc 下服务不存在
	return nil
}

func sync(ca *models.CloudAuth, regions []string) (models.TunnelCreateApiList, error) {
	res := make(models.TunnelCreateApiList, 0)
	for _, rg := range regions {
		client, err := fc.NewClient(
			fmt.Sprintf("%s.%s.fc.aliyuncs.com", ca.AccessId, rg),
			"2016-08-15", ca.AccessKey, ca.AccessSecret)
		if err != nil {
			return res, err
		}
		// 先同步函数
		respC, err := client.ListFunctions(fc.NewListFunctionsInput(serviceName))
		if err != nil {
			e, ok := err.(*fc.ServiceError)
			if ok && e.HTTPStatus == http.StatusNotFound {
				// 说明没有 service，继续下一个地区好了
				continue
			} else {
				return res, err
			}
		}
		for _, c := range respC.Functions {
			// 检查下 Type, 现在临时塞到了desc中了
			if *c.Description == "" {
				xlog.Warn("sdk", "发现了不正确的隧道", "fc_name", *c.FunctionName, "fc_type", *c.Description)
				continue
			}

			var tun = models.NewTunnelCreateApi()
			tun.UniqID = c.FunctionID
			tun.Name = c.FunctionName
			tun.Config = &models.TunnelConfig{
				CPU:      *c.CPU,
				Region:   rg,
				Memory:   *c.MemorySize,
				Instance: *c.InstanceConcurrency,

				TLS: true, // 默认同步过来都打开
				Tor: false,
			}
			if len(c.EnvironmentVariables) > 0 {
				for key, value := range c.EnvironmentVariables {
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
			*tun.Type = enum.TransTunnelType(*c.Description)
			tun.Port = c.CAPort
			// 再同步触发器来填充隧道的地址
			respT, err := client.ListTriggers(fc.NewListTriggersInput(serviceName, *c.FunctionName))
			if err == nil {
				for _, t := range respT.Triggers {
					if *t.TriggerType == "http" {
						*tun.Addr = strings.Replace(t.UrlInternet, "https://", "", -1)
						tun.Config.FcAuthType = enum.TransAuthType(t.TriggerConfig.AuthType)
					}
				}
				*tun.Status = enum.TunnelActive
				res = append(res, tun)
			} else {
				*tun.Status = enum.TunnelError
				*tun.StatusMessage = err.Error()
			}
			res = append(res, tun)
		}
	}
	return res, nil
}
