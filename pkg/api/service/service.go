package service

import (
	"context"
	"log/slog"

	"github.com/DVKunion/SeaMoon/pkg/api/database/drivers"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
)

var (
	SVC = &svc{
		&auth{},
		&config{},
		&provider{},
		&proxy{},
		&tunnel{},
	}

	// error list
)

type svc struct {
	*auth
	*config
	*provider
	*proxy
	*tunnel
}

func init() {
	drivers.RegisterMigrate(func() {
		slog.Info("未查询到本地数据，初始化默认配置......")
		for _, conf := range models.DefaultConfig {
			_ = SVC.CreateConfig(context.Background(), &conf)
		}
		slog.Info("未查询到本地数据，初始化默认账户......")
		for _, ca := range models.DefaultAuth {
			_ = SVC.CreateAuth(context.Background(), &ca)
		}
	})
}
