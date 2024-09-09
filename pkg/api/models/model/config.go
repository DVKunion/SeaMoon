package model

import (
	"gorm.io/gorm"
)

// Config 系统标准配置表
type Config struct {
	gorm.Model

	Key   string
	Value string
}

// ConfigApi 对外暴露接口
type ConfigApi map[string]string

func (c ConfigApi) ToModel() ConfigList {
	var res = make([]*Config, 0)
	for k, v := range c {
		cf := &Config{
			Key:   k,
			Value: v,
		}
		res = append(res, cf)
	}
	return res
}

type ConfigList []*Config

func (cl ConfigList) ToApi() ConfigApi {
	var res = make(ConfigApi)
	for _, s := range cl {
		res[s.Key] = s.Value
	}
	return res
}
