package api

import (
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

func Call(client *openapi.Client, params *openapi.Params, body map[string]interface{}, obj interface{}) error {
	// runtime options
	runtime := &util.RuntimeOptions{}
	request := &openapi.OpenApiRequest{}
	if body != nil {
		request.Body = body
	}
	if obj == nil {
		obj = make(map[string]interface{})
	}
	resp, err := client.CallApi(params, request, runtime)
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	bs, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, obj)
}

func NewClient(ak, sk, ep string) (*openapi.Client, error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(ak),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(sk),
	}
	config.Endpoint = tea.String(ep)
	return openapi.NewClient(config)
}

func NewGetBillingsParams() *openapi.Params {
	return newParams("QueryAccountBalance", "2017-12-14", "RPC", "POST", "/")
}

func NewCreateServiceParams() *openapi.Params {
	return newParams("CreateService", "2021-04-06", "FC", "POST", "/2021-04-06/services")
}

func NewListFCParams() *openapi.Params {
	return newParams("ListFunctions", "2021-04-06", "FC", "GET", "/2021-04-06/services/seamoon/functions")
}

func NewCreateFCParams() *openapi.Params {
	return newParams("CreateFunction", "2021-04-06", "FC", "POST", "/2021-04-06/services/seamoon/functions")
}

func NewDeleteFCParams(fc string) *openapi.Params {
	return newParams("DeleteFunction", "2021-04-06", "FC", "DELETE", "/2021-04-06/services/seamoon/functions/"+fc)
}

func NewListTriggerParams(fc string) *openapi.Params {
	return newParams("ListTriggers", "2021-04-06", "FC", "GET", "/2021-04-06/services/seamoon/functions/"+fc+"/triggers")
}

func NewCreateTriggerParams(fc string) *openapi.Params {
	return newParams("CreateTrigger", "2021-04-06", "FC", "POST", "/2021-04-06/services/seamoon/functions/"+fc+"/triggers")
}

func NewDeleteTriggerParams(fc, tr string) *openapi.Params {
	return newParams("DeleteTrigger", "2021-04-06", "FC", "DELETE", "/2021-04-06/services/seamoon/functions/"+fc+"/triggers/"+tr)
}

func newParams(action, version, style, method, path string) *openapi.Params {
	return &openapi.Params{
		// 接口名称
		Action: tea.String(action),
		// 接口版本
		Version: tea.String(version),
		// 接口协议
		Protocol: tea.String("HTTPS"),
		// 接口 HTTP 方法
		Method:   tea.String(method),
		AuthType: tea.String("AK"),
		Style:    tea.String(style),
		// 接口 PATH
		Pathname: tea.String(path),
		// 接口请求体内容格式
		ReqBodyType: tea.String("json"),
		// 接口响应体内容格式
		BodyType: tea.String("json"),
	}
}
