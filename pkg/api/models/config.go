package models

import (
	"gorm.io/gorm"

	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

var DefaultConfig = []Config{
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
		Key:   "auto_start",
		Value: "true",
	},
	{
		Key:   "version",
		Value: xlog.Version,
	},
}

// Config 系统标准配置表
type Config struct {
	gorm.Model

	Key   string
	Value string
}

type ConfigList []*Config

// ConfigApi 对外暴露接口
type ConfigApi struct {
	// 为了 web 方便, 直接转化成对应的 key 了
	ControlAddr string `json:"control_addr"`
	ControlPort string `json:"control_port"`
	ControlLog  string `json:"control_log"`
	AutoStart   string `json:"auto_start"`

	Version string `json:"version"`
}

func (c *ConfigApi) ToModel() ConfigList {
	// 由于目前东西比较少，懒得写反射了；后续如果也有这种 KV 存储转换的需求，可以抽到 models 公共方法中
	var res = make([]*Config, 0)
	res = append(res, &Config{
		Key:   "control_addr",
		Value: c.ControlAddr,
	})
	res = append(res, &Config{
		Key:   "control_port",
		Value: c.ControlPort,
	})
	res = append(res, &Config{
		Key:   "control_log",
		Value: c.ControlLog,
	})
	res = append(res, &Config{
		Key:   "auto_start",
		Value: c.AutoStart,
	})

	return res
}

func (cl ConfigList) ToApi() *ConfigApi {
	var res = &ConfigApi{}
	for _, s := range cl {
		switch s.Key {
		case "control_addr":
			res.ControlAddr = s.Value
		case "control_port":
			res.ControlPort = s.Value
		case "control_log":
			res.ControlLog = s.Value
		case "auto_start":
			res.AutoStart = s.Value
		case "version":
			res.Version = s.Value
		}
	}
	return res
}
