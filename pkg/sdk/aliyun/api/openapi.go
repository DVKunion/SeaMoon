package api

import (
	"encoding/json"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	fc "github.com/alibabacloud-go/fc-20230330/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// NewBillingClient 创建用于 Billing API 的客户端（仍使用旧的 openapi 客户端）
func NewBillingClient(ak, sk, ep string) (*openapi.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(ak),
		AccessKeySecret: tea.String(sk),
	}
	config.Endpoint = tea.String(ep)
	return openapi.NewClient(config)
}

// NewFCClient 创建新的 FC SDK 客户端
func NewFCClient(ak, sk, ep string) (*fc.Client, error) {
	config := &openapiutil.Config{
		AccessKeyId:     tea.String(ak),
		AccessKeySecret: tea.String(sk),
		Endpoint:        tea.String(ep),
	}
	return fc.NewClient(config)
}

// CallBilling 调用 Billing API（仍使用旧的 openapi 客户端）
func CallBilling(client *openapi.Client, params *openapi.Params, body map[string]interface{}, obj interface{}) error {
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

// NewGetBillingsParams 创建 Billing API 参数
func NewGetBillingsParams() *openapi.Params {
	return newParams("QueryAccountBalance", "2017-12-14", "RPC", "POST", "/")
}

func newParams(action, version, style, method, path string) *openapi.Params {
	return &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String(version),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String(method),
		AuthType:    tea.String("AK"),
		Style:       tea.String(style),
		Pathname:    tea.String(path),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
}
