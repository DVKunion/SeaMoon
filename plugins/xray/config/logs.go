package config

import "github.com/xtls/xray-core/infra/conf"

func WithLogs(level string) Options {
	return func(cf *conf.Config) {
		cf.LogConfig = &conf.LogConfig{
			AccessLog: "",
			ErrorLog:  "",
			LogLevel:  level,
		}
	}
}
