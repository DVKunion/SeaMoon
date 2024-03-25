package tencent

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	billing "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	fcError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/consts"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

const (
	BillEndPoint = "billing.tencentcloudapi.com"
	SCFEndPoint  = "scf.tencentcloudapi.com"
	CAMEndPoint  = "cam.tencentcloudapi.com"
)

type triggerDesc struct {
	Api     triggerDescApi     `json:"api"`
	Service triggerDescService `json:"service"`
	Release struct {
		EnvironmentName string `json:"environmentName"`
	} `json:"release"`
}

type triggerDescApi struct {
	AuthRequired  string `json:"authRequired"`
	AuthType      string `json:"authType"`
	RequestConfig struct {
		Method string `json:"method"`
		Path   string `json:"path"`
	} `json:"requestConfig"`
	IsIntegratedResponse string `json:"isIntegratedResponse"`
	IsBase64Encoded      string `json:"isBase64Encoded"`
}

type triggerDescService struct {
	ServiceName string   `json:"serviceName"`
	ServiceType string   `json:"serviceType"`
	SubDomain   string   `json:"subDomain"`
	Tags        []string `json:"tags"`
}

type triggerResp struct {
	Service struct {
		SubDomain string `json:"subDomain"`
	} `json:"service"`
}

type fcInfo struct {
	detail *scf.GetFunctionResponseParams
	region string
	addr   string
	auth   string
}

func createRole(ca *models.CloudAuth) error {
	credential := common.NewCredential(
		ca.AccessKey,
		ca.AccessSecret,
	)

	cpf := profile.NewClientProfile()
	// 需要授权: SCF、 API GateWay
	cpf.HttpProfile.Endpoint = "cam.tencentcloudapi.com"
	roleClient, _ := cam.NewClient(credential, "ap-guangzhou", cpf)

	apiGateWayPolicy := "{\"version\":\"2.0\",\"statement\":[{\"action\":\"name/sts:AssumeRole\",\"effect\":\"allow\",\"principal\":{\"service\":\"apigateway.qcloud.com\"}}]}"
	fcPolicy := "{\"version\":\"2.0\",\"statement\":[{\"action\":\"name/sts:AssumeRole\",\"effect\":\"allow\",\"principal\":{\"service\":\"scf.qcloud.com\"}}]}"

	roleRequest := cam.NewCreateRoleRequest()
	roleRequest.RoleName = common.StringPtr("ApiGateWay_QCSRole")
	roleRequest.PolicyDocument = common.StringPtr(apiGateWayPolicy)
	roleRequest.Description = common.StringPtr("API 网关(API Gateway)对您的腾讯云资源进行访问操作，含上传日志、获取日志游标、下载日志、获取日志主题信息等。")

	_, err := roleClient.CreateRole(roleRequest)
	if err != nil {
		if err, ok := err.(*fcError.TencentCloudSDKError); !ok || err.Code != "InvalidParameter.RoleNameInUse" ||
			err.Message != "role name in use" {
			return err
		}
	}

	attachRolePolicyRequest := cam.NewAttachRolePolicyRequest()
	attachRolePolicyRequest.PolicyName = common.StringPtr("QcloudAccessForAPIGatewayRoleInSCFTrigger")
	attachRolePolicyRequest.AttachRoleName = roleRequest.RoleName

	_, err = roleClient.AttachRolePolicy(attachRolePolicyRequest)
	if err != nil {
		return err
	}

	roleRequest = cam.NewCreateRoleRequest()
	roleRequest.RoleName = common.StringPtr("SCF_QcsRole")
	roleRequest.PolicyDocument = common.StringPtr(fcPolicy)
	roleRequest.Description = common.StringPtr("云函数(SCF)操作权限含创建对象存储(COS)触发器，拉取代码包等；含创建API网关(API Gateway)触发器等；含消创建息队列(CMQ)触发器等；含投递日志服务(CLS)日志等。")

	_, err = roleClient.CreateRole(roleRequest)
	if err != nil {
		if err, ok := err.(*fcError.TencentCloudSDKError); !ok || err.Code != "InvalidParameter.RoleNameInUse" ||
			err.Message != "role name in use" {
			return err
		}
	}

	attachRolePolicyRequest = cam.NewAttachRolePolicyRequest()
	attachRolePolicyRequest.PolicyName = common.StringPtr("QcloudAccessForScfRole")
	attachRolePolicyRequest.AttachRoleName = roleRequest.RoleName

	_, err = roleClient.AttachRolePolicy(attachRolePolicyRequest)
	if err != nil {
		return err
	}

	return nil
}

func getAmount(ca *models.CloudAuth) (float64, error) {
	credential := common.NewCredential(
		ca.AccessKey,
		ca.AccessSecret,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = BillEndPoint

	client, err := billing.NewClient(credential, "", cpf)

	if err != nil {
		return 0, err
	}

	// 构造请求
	request := billing.NewDescribeAccountBalanceRequest()

	// 发送请求
	response, err := client.DescribeAccountBalance(request)

	if err != nil {
		return 0, err
	}

	balance := *response.Response.Balance

	return float64(balance) / 100, nil
}

func deploy(ca *models.CloudAuth, tun *models.Tunnel) (string, error) {
	credential := common.NewCredential(
		ca.AccessKey,
		ca.AccessSecret,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = SCFEndPoint

	// 创建 SFC 客户端
	client, err := scf.NewClient(credential, tun.Config.Region, cpf)

	if err != nil {
		return "", err
	}

	// SCF 需要一个 namespace
	// 同阿里云一样，这里硬编码一个 code
	nsRequest := scf.NewCreateNamespaceRequest()

	nsRequest.Namespace = common.StringPtr(serviceName)
	nsRequest.Description = common.StringPtr(serviceDesc)
	_, err = client.CreateNamespace(nsRequest)
	if err != nil {
		// 如果错误是 ns 存在，则忽略。
		if err, ok := err.(*fcError.TencentCloudSDKError); !ok || err.Code != scf.RESOURCEINUSE_NAMESPACE {
			return "", err
		}
	}

	// 创建函数
	request := scf.NewCreateFunctionRequest()

	// 查询的时候只能用模糊匹配，sb, 得用个不会模糊的前缀区分
	fcName := "a" + strconv.Itoa(int(tun.CreatedAt.Unix())) + "-" + *tun.Name

	request.Namespace = common.StringPtr(serviceName)
	request.FunctionName = common.StringPtr(fcName)
	request.Description = common.StringPtr(string(*tun.Type))
	request.Type = common.StringPtr("HTTP")
	// 腾讯云没有 cpu 大小设置选项
	request.MemorySize = common.Int64Ptr(int64(tun.Config.Memory))
	// 需要记得打开 ws 支持
	request.ProtocolType = common.StringPtr("WS")
	request.ProtocolParams = &scf.ProtocolParams{
		WSParams: &scf.WSParams{
			IdleTimeOut: common.Uint64Ptr(60),
		},
	}
	request.Timeout = common.Int64Ptr(600)

	request.AutoCreateClsTopic = common.StringPtr("FALSE")

	request.Code = &scf.Code{
		ImageConfig: &scf.ImageConfig{
			ImageType: common.StringPtr("personal"),
			ImageUri:  common.StringPtr(strings.Join([]string{RegistryEndPoint[tun.Config.Region], consts.Version}, ":")),
			Command:   common.StringPtr("/app/seamoon"),
			Args:      common.StringPtr("server -p " + strconv.Itoa(int(*tun.Port)) + " -t " + string(*tun.Type)),
			ImagePort: common.Int64Ptr(int64(*tun.Port)),
		},
	}

	if tun.Config.Tor {
		request.Environment = &scf.Environment{
			Variables: []*scf.Variable{
				{
					Key:   common.StringPtr("SEAMOON_TOR"),
					Value: common.StringPtr("true"),
				},
			},
		}
	}

	request.PublicNetConfig = &scf.PublicNetConfigIn{
		PublicNetStatus: common.StringPtr("ENABLE"),
		EipConfig: &scf.EipConfigIn{
			EipStatus: common.StringPtr("DISABLE"),
		},
	}

	request.InstanceConcurrencyConfig = &scf.InstanceConcurrencyConfig{
		DynamicEnabled: common.StringPtr("FALSE"),
		MaxConcurrency: common.Uint64Ptr(uint64(tun.Config.Instance)),
	}

	_, err = client.CreateFunction(request)
	if err != nil {
		if err, ok := err.(*fcError.TencentCloudSDKError); !ok || err.Code != scf.RESOURCEINUSE_FUNCTION {
			return "", err
		}
	}

	// 查询等待状态正常
	// 尝试查询30次，每次等待2秒, 共计一分钟
	cnt := 0
	for cnt < 30 {
		eRequest := scf.NewListFunctionsRequest()
		eRequest.Namespace = common.StringPtr(serviceName)
		eRequest.SearchKey = common.StringPtr(fcName)

		fc, err := client.ListFunctions(eRequest)
		if err != nil {
			return "", err
		}
		if *fc.Response.TotalCount != 1 {
			return "", errors.New(errors.SDKFCInfoError)
		}
		xlog.Info(xlog.SDKWaitingFCStatus, "status", *fc.Response.Functions[0].Status, "cnt", cnt)
		switch *fc.Response.Functions[0].Status {
		case "Active":
			cnt = 31
		case "Creating":
			time.Sleep(2 * time.Second)
			cnt++
			continue
		default:
			return "", errors.New(*fc.Response.Functions[0].StatusDesc)
		}
	}

	// 尝试创建触发器
	r := scf.NewCreateTriggerRequest()
	r.TriggerName = common.StringPtr("apigw")
	r.FunctionName = common.StringPtr(fcName)
	r.Type = common.StringPtr("apigw")

	config, err := json.Marshal(&triggerDesc{
		Api: triggerDescApi{
			AuthRequired: func() string {
				if tun.Config.FcAuthType == enum.AuthParis {
					return "TRUE"
				}
				return "FALSE"
			}(),
			RequestConfig: struct {
				Method string `json:"method"`
				Path   string `json:"path"`
			}{Method: "GET", Path: "/"},
			IsIntegratedResponse: "FALSE",
			IsBase64Encoded:      "FALSE",
		},
		Service: triggerDescService{
			ServiceName: "SCF_API_SERVICE",
			ServiceType: "BASIC",
		},
		Release: struct {
			EnvironmentName string `json:"environmentName"`
		}{EnvironmentName: "release"},
	})

	// 触发器配置参数叫做 desc...
	r.TriggerDesc = common.StringPtr(string(config))
	r.Namespace = common.StringPtr(serviceName)

	response, err := client.CreateTrigger(r)
	if err != nil {
		return "", err
	}

	extractor := &triggerResp{}
	desc := *response.Response.TriggerInfo.TriggerDesc
	if err := json.Unmarshal([]byte(desc), extractor); err != nil {
		return "", err
	}

	return extractor.Service.SubDomain, nil
}

func destroy(ca *models.CloudAuth, tun *models.Tunnel) error {
	credential := common.NewCredential(
		ca.AccessKey,
		ca.AccessSecret,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = SCFEndPoint

	// 创建 SFC 客户端
	client, err := scf.NewClient(credential, tun.Config.Region, cpf)

	if err != nil {
		return err
	}

	fcName := "a" + strconv.Itoa(int(tun.CreatedAt.Unix())) + "-" + *tun.Name

	// 先删除触发器
	r := scf.NewDeleteTriggerRequest()
	r.TriggerName = common.StringPtr("apigw")
	r.FunctionName = common.StringPtr(fcName)
	r.Namespace = common.StringPtr(serviceName)
	if _, err = client.DeleteTrigger(r); err != nil {
		return err
	}
	// 再删除函数
	request := scf.NewDeleteFunctionRequest()
	request.FunctionName = common.StringPtr(fcName)
	request.Namespace = common.StringPtr(serviceName)
	if _, err = client.DeleteFunction(request); err != nil {
		return err
	}
	// 不要删除 ns, ns 又不花钱
	return nil
}

func sync(ca *models.CloudAuth, regions []string) ([]fcInfo, error) {
	var res = make([]fcInfo, 0)

	credential := common.NewCredential(
		ca.AccessKey,
		ca.AccessSecret,
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = SCFEndPoint
	for _, rg := range regions {
		// 创建 SFC 客户端
		client, err := scf.NewClient(credential, rg, cpf)

		if err != nil {
			return res, err
		}

		// 需要先看一下是否存在 ns
		nsRequest := scf.NewListNamespacesRequest()

		nsRequest.SearchKey = []*scf.SearchKey{
			{
				Key:   common.StringPtr("Namespace"),
				Value: common.StringPtr(serviceName),
			},
		}
		response, err := client.ListNamespaces(nsRequest)
		if err != nil {
			return res, err
		}

		if *response.Response.TotalCount == 0 {
			return res, nil
		}

		request := scf.NewListFunctionsRequest()
		request.Namespace = common.StringPtr(serviceName)
		request.Limit = common.Int64Ptr(999999)
		fcList, err := client.ListFunctions(request)
		if err != nil {
			return nil, err
		}
		// list 不够详细，需要继续处理
		for _, fc := range fcList.Response.Functions {
			target := fcInfo{
				region: rg,
			}
			req := scf.NewGetFunctionRequest()
			req.FunctionName = fc.FunctionName
			req.Namespace = fc.Namespace
			fcd, err := client.GetFunction(req)
			if err != nil {
				xlog.Error(errors.SDKFCDetailError, "name", *fc.FunctionName, "err", err)
				continue
			} else {
				target.detail = fcd.Response
				// 解析触发器
				trigger := fcd.Response.Triggers
				if len(trigger) < 1 {
					xlog.Error(errors.SDKTriggerError, "name", *fc.FunctionName, "err", err)
				} else {
					var tri triggerDesc
					err := json.Unmarshal([]byte(*trigger[0].TriggerDesc), &tri)
					if err != nil {
						xlog.Error(errors.SDKTriggerUnmarshalError, "name", *fc.FunctionName, "err", err)
					}
					target.addr = strings.Replace(tri.Service.SubDomain, "https://", "", -1)
					target.auth = tri.Api.AuthType
				}
			}

			res = append(res, target)
		}
	}
	return res, nil
}
