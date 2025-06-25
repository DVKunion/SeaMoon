package openapi

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

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

func NewListFCParams() *openapi.Params {
	return newParams("ListFunctions", "2021-04-06", "FC", "GET", "/2021-04-06/services/seamoon/functions")
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
