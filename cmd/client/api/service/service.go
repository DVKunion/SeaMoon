package service

import (
	"log/slog"
	"reflect"

	"github.com/DVKunion/SeaMoon/cmd/client/api/database"
	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
)

var serviceFactory = map[string]ApiService{}

type Condition struct {
	Key   string
	Value interface{}
}

type ApiService interface {
	Count(cond ...Condition) int64
	List(page, size int, preload bool, cond ...Condition) interface{}
	GetById(id uint) interface{}
	Create(obj interface{}) interface{}
	Update(id uint, obj interface{}) interface{}
	Delete(id uint)
}

func GetService(t string) ApiService {
	return serviceFactory[t]
}

func Exist(svc ApiService, cond ...Condition) bool {
	s := svc.List(0, 1, false, cond...)

	// 重复名称的账户不允许创建
	if reflect.TypeOf(s).Kind() == reflect.Slice && reflect.ValueOf(s).Len() > 0 {
		return true
	}

	return false
}

func init() {
	serviceFactory["proxy"] = ProxyService{}
	serviceFactory["provider"] = CloudProviderService{}
	serviceFactory["config"] = SystemConfigService{}
	serviceFactory["auth"] = AuthService{}
	serviceFactory["tunnel"] = TunnelService{}

	// 注册初始化相关信息
	database.RegisterMigrate(func() {
		slog.Info("未查询到本地数据，初始化默认配置......")
		for _, conf := range models.DefaultSysConfig {
			serviceFactory["config"].Create(&conf)
		}
		slog.Info("未查询到本地数据，初始化默认账户......")
		for _, auth := range models.DefaultAuth {
			serviceFactory["auth"].Create(&auth)
		}
	})
}
