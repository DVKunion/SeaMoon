package config

import (
	"github.com/xtls/xray-core/infra/conf"

	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

func WithDefaultPolicy() Options {
	return func(cf *conf.Config) {
		cf.Policy = &conf.PolicyConfig{
			Levels: map[uint32]*conf.Policy{
				// 0为代理策略
				0: {
					UplinkOnly:   tools.Uint32Ptr(2), // 出站时，当服务端断开连接后最大中断等待时间，默认2
					DownlinkOnly: tools.Uint32Ptr(5), // 进站时，当客户端断开链接后最大中断等待时间，默认5
					// 暂时未开启用户流量统计
					StatsUserUplink:   false,
					StatsUserDownlink: false,
				},
				// 1为管理api策略，不计入统计相关
				1: {
					StatsUserUplink:   false,
					StatsUserDownlink: false,
				},
			},
			System: &conf.SystemPolicy{
				// 默认不开启全局流量统计
				StatsOutboundDownlink: false,
				StatsOutboundUplink:   false,
				StatsInboundDownlink:  false,
				StatsInboundUplink:    false,
			},
		}
	}
}

func WithInboundCalculate() Options {
	return func(conf *conf.Config) {
		conf.Policy.System.StatsInboundUplink = true
		conf.Policy.System.StatsInboundDownlink = true
	}
}

func WithOutboundCalculate() Options {
	return func(conf *conf.Config) {
		conf.Policy.System.StatsOutboundDownlink = true
		conf.Policy.System.StatsOutboundUplink = true
	}
}
