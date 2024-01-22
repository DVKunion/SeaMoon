package sdk

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

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
	"github.com/DVKunion/SeaMoon/pkg/xlog"
)

var (
	// 阿里云在 fc 上层还有一套 service 的概念，为了方便管理，这里硬编码了 service 的内容。
	serviceName = "seamoon"
	serviceDesc = "seamoon service"
)

type ALiYunSDK struct {
}

type Resp struct {
	StatusCode int                    `json:"statusCode"`
	Headers    map[string]interface{} `json:"headers"`
	Body       struct {
		Code      string                 `json:"Code"`
		Message   string                 `json:"Message"`
		RequestId string                 `json:"RequestId"`
		Success   bool                   `json:"Success"`
		Data      map[string]interface{} `json:"Data"`
	} `json:"body"`
}

func (a *ALiYunSDK) Auth(providerId uint) error {
	provider := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: &provider.CloudAuth.AccessKey,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: &provider.CloudAuth.AccessSecret,
	}
	// todo: seems ALiYunBillingMap is not right here
	config.Endpoint = tea.String("business.aliyuncs.com")
	client, err := bss.NewClient(config)
	if err != nil {
		return err
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
	res, err := client.CallApi(params, request, runtime)
	if err != nil {
		return err
	}
	bs, err := json.Marshal(res)
	if err != nil {
		return err
	}
	var r Resp
	err = json.Unmarshal(bs, &r)
	if err != nil {
		return err
	}
	if r.StatusCode != 200 || r.Body.Code != "200" {
		return errors.New(r.Body.Message)
	}
	amount, err := strconv.ParseFloat(strings.Replace(r.Body.Data["AvailableAmount"].(string), ",", "", -1), 64)
	if err != nil {
		return err
	}
	*provider.Amount = amount

	// todo: 查询总花费
	service.GetService("provider").Update(provider.ID, provider)
	return nil
}

func (a *ALiYunSDK) Deploy(providerId uint, tun *models.Tunnel) error {
	provider := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	// 原生的库是真tm的难用，
	client, err := fc.NewClient(
		fmt.Sprintf("%s.%s.fc.aliyuncs.com", provider.CloudAuth.AccessId, *provider.Region),
		"2016-08-15", provider.CloudAuth.AccessKey, provider.CloudAuth.AccessSecret)
	if err != nil {
		return err
	}
	// 先尝试是否已经存在了 svc
	_, err = client.GetService(fc.NewGetServiceInput(serviceName))
	if err != nil {
		err := err.(*fc.ServiceError)
		if err.HTTPStatus == http.StatusNotFound {
			// 说明 service 空了，先创建svc
			_, err := client.CreateService(fc.NewCreateServiceInput().
				WithServiceName(serviceName).
				WithDescription(serviceDesc))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	p, err := strconv.Atoi(*tun.Port)
	funcName := *tun.Name
	// 有了服务了，现在来创建函数
	respC, err := client.CreateFunction(fc.NewCreateFunctionInput(serviceName).
		WithFunctionName(funcName).
		WithDescription(string(*tun.Type)).
		WithRuntime("custom-container").
		WithCPU(tun.TunnelConfig.CPU).
		WithMemorySize(tun.TunnelConfig.Memory).
		WithHandler("main").
		WithDisk(512).
		WithInstanceConcurrency(tun.TunnelConfig.Instance).
		WithCAPort(int32(p)).
		WithInstanceType("e1"). // 性能实例
		WithTimeout(300).
		WithCustomContainerConfig(fc.NewCustomContainerConfig().
			WithImage("registry.cn-hongkong.aliyuncs.com/seamoon/seamoon:dev").
			WithCommand("[\"./seamoon\"]").
			WithArgs("[\"server\"]")))
	if err != nil {
		return err
	}
	fmt.Println(respC)
	// 有了函数了，接下来创建 trigger
	respT, err := client.CreateTrigger(fc.NewCreateTriggerInput(serviceName, funcName).
		WithTriggerType("http").
		WithTriggerName(string(*tun.Type)).
		WithTriggerConfig(fc.TriggerConfig{
			Methods:            []string{"GET", "POST"},
			AuthType:           "anonymous",
			DisableURLInternet: false,
		}))
	if err != nil {
		return err
	}
	fmt.Println(respT)
	// 创建成功了, 查一下
	respTS, err := client.GetTrigger(fc.NewGetTriggerInput(serviceName, funcName, string(*tun.Type)))
	if err != nil {
		return err
	}

	*tun.Addr = strings.Replace(respTS.UrlInternet, "https://", "", -1)
	*tun.Status = tunnel.ACTIVE
	// 更新 tunnel
	service.GetService("tunnel").Update(tun.ID, tun)
	return nil
}

func (a *ALiYunSDK) Destroy(providerId uint, tun *models.Tunnel) error {
	return nil
}

func (a *ALiYunSDK) SyncFC(providerId uint) error {
	provider := service.GetService("provider").GetById(providerId).(*models.CloudProvider)
	client, err := fc.NewClient(
		fmt.Sprintf("%s.%s.fc.aliyuncs.com", provider.CloudAuth.AccessId, *provider.Region),
		"2016-08-15", provider.CloudAuth.AccessKey, provider.CloudAuth.AccessSecret)
	if err != nil {
		return err
	}
	// 先同步函数
	respC, err := client.ListFunctions(fc.NewListFunctionsInput(serviceName))
	if err != nil {
		e, ok := err.(*fc.ServiceError)
		if ok && e.HTTPStatus == http.StatusNotFound {
			// 说明没有 service，甭同步了
			return nil
		} else {
			return err
		}
	}
	for _, c := range respC.Functions {
		// 判断下存不存在吧，不然每次同步都会整出来一堆
		if exist := service.Exist(service.GetService("tunnel"), service.Condition{
			Key:   "NAME",
			Value: *c.FunctionName,
		}, service.Condition{
			Key:   "cloud_provider_id",
			Value: provider.ID,
		}); exist {
			continue
		}
		// 还得检查下 Type, 现在临时塞到了desc中了
		if *c.Description == "" {
			xlog.Warn("sdk", "发现了不正确的隧道", "provider_id", provider.ID,
				"fc_name", *c.FunctionName, "fc_type", *c.Description,
			)
			continue
		}

		var tun = models.Tunnel{
			CloudProviderId: provider.ID,
			Name:            c.FunctionName,

			TunnelConfig: &models.TunnelConfig{
				CPU:      *c.CPU,
				Memory:   *c.MemorySize,
				Instance: *c.InstanceConcurrency,

				// todo: 这里太糙了
				TLS: true, // 默认同步过来都打开
				Tor: func() bool {
					// 如果是 开启 Tor 的隧道，需要有环境变量
					return len(c.EnvironmentVariables) > 0
				}(),
			},
		}
		// 自动填充防止空指针
		models.AutoFull(&tun)
		*tun.Type = tunnel.TranType(*c.Description)
		*tun.Port = strconv.Itoa(int(*c.CAPort))
		*tun.Status = tunnel.ACTIVE
		// 再同步触发器来填充隧道的地址
		respT, err := client.ListTriggers(fc.NewListTriggersInput(serviceName, *c.FunctionName))
		if err == nil {
			for _, t := range respT.Triggers {
				if *t.TriggerType == "http" {
					// todo: 这里太糙了
					*tun.Addr = strings.Replace(t.UrlInternet, "https://", "", -1)
					tun.TunnelConfig.TunnelAuthType = types.TransAuthType(t.TriggerConfig.AuthType)
				}
			}
		} else {
			// todo: 打印一下失败的的日志
		}
		// 最后更新一下
		service.GetService("tunnel").Create(&tun)
	}

	return nil
}
