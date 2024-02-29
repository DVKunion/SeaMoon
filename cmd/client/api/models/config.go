package models

import (
	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/consts"
)

var DefaultSysConfig = []SystemConfig{
	{
		Key:   "control_addr",
		Value: "0.0.0.0",
	},
	{
		Key:   "control_port",
		Value: "7778",
	},
	{
		Key:   "control_log",
		Value: "seamoon-web.log",
	},
	{
		Key:   "version",
		Value: consts.Version,
	},
}

// SystemConfig 系统标准配置表
type SystemConfig struct {
	gorm.Model

	Key   string
	Value string
}

// SystemConfigApi 对外暴露接口
type SystemConfigApi struct {
	// 为了 web 方便, 直接转化成对应的 key 了
	ControlAddr string `json:"control_addr"`
	ControlPort string `json:"control_port"`
	ControlLog  string `json:"control_log"`

	Version string `json:"version"`
}

// ToModel SysConfig 比较特殊，存储时候是一个 KEY-VALUE 模式，
// 因此让他自己重新实现一下 ToModel
func (p *SystemConfigApi) ToModel() []SystemConfig {
	// 由于目前东西比较少，懒得写反射了；后续如果也有这种 KV 存储转换的需求，可以抽到 models 公共方法中
	var res = make([]SystemConfig, 0)
	res = append(res, SystemConfig{
		Key:   "control_addr",
		Value: p.ControlAddr,
	})
	res = append(res, SystemConfig{
		Key:   "control_port",
		Value: p.ControlPort,
	})
	res = append(res, SystemConfig{
		Key:   "control_log",
		Value: p.ControlLog,
	})
	res = append(res, SystemConfig{
		Key:   "version",
		Value: p.Version,
	})

	return res
}

func ToSystemConfigApi(sc []*SystemConfig) SystemConfigApi {
	var res = SystemConfigApi{}
	for _, s := range sc {
		switch s.Key {
		case "control_addr":
			res.ControlAddr = s.Value
		case "control_port":
			res.ControlPort = s.Value
		case "control_log":
			res.ControlLog = s.Value
		case "version":
			res.Version = s.Value
		}
	}
	return res
}
