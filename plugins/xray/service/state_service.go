package service

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/app/stats/command"
)

type StateService struct {
	cc command.StatsServiceClient
}

// GetSysStates 获取系统状态
func (s *StateService) GetSysStates(ctx context.Context) (*command.SysStatsResponse, error) {
	return s.cc.GetSysStats(ctx, &command.SysStatsRequest{})
}

// GetUploadTrafficByEmail 获取用户的流量数据: 上行
func (s *StateService) GetUploadTrafficByEmail(ctx context.Context, email string) (int64, error) {
	return s.queryTraffic(ctx, fmt.Sprintf("user>>>%s>>>traffic>>>uplink", email), false)
}

// GetDownloadTrafficByEmail 获取用户的流量数据: 下行
func (s *StateService) GetDownloadTrafficByEmail(ctx context.Context, email string) (int64, error) {
	return s.queryTraffic(ctx, fmt.Sprintf("user>>>%s>>>traffic>>>downlink", email), false)
}

// GetUploadTrafficByTag 获取全局流量: 上行
// bt: inbound / outbound
// tag: 对应要查询的 tag
func (s *StateService) GetUploadTrafficByTag(ctx context.Context, bt string, tag string) (int64, error) {
	return s.queryTraffic(ctx, fmt.Sprintf("%s>>>%s>>>traffic>>>uplink", bt, tag), false)
}

// GetDownloadTrafficByTag 获取全局流量: 下行
// bt: bound_type: inbound / outbound
// tag: 对应要查询的 tag
func (s *StateService) GetDownloadTrafficByTag(ctx context.Context, bt string, tag string) (int64, error) {
	return s.queryTraffic(ctx, fmt.Sprintf("%s>>>%s>>>traffic>>>downlink", bt, tag), false)
}

// queryTraffic: 查询流量原生方法 ptn: 查询语句 reset: 是否重置
func (s *StateService) queryTraffic(ctx context.Context, ptn string, reset bool) (int64, error) {
	// 如果查无此用户或 bound 则返回-1, 默认值 -1
	var traffic int64 = -1
	resp, err := s.cc.QueryStats(ctx, &command.QueryStatsRequest{
		// 这里是查询语句，例如 “user>>>love@xray.com>>>traffic>>>uplink” 表示查询用户 email 为 love@xray.com 在所有入站中的上行流量
		Pattern: ptn,
		// 是否重置流量信息(true, false)，即完成查询后是否把流量统计归零
		Reset_: reset,
	})
	if err != nil {
		return traffic, err
	}
	// Get traffic data
	stat := resp.GetStat()
	// 判断返回 是否成功
	// 返回样例，value 值是我们需要的: [name:"inbound>>>proxy0>>>traffic>>>downlink" value:348789]
	if len(stat) != 0 {
		// 返回流量数据 byte
		traffic = stat[0].Value
	}

	return traffic, nil
}
