package service

import (
	"context"

	"github.com/DVKunion/SeaMoon/pkg/api/database/drivers"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
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
		xlog.Info(xlog.DatabaseConfigInit)
		for _, conf := range models.DefaultConfig {
			_ = SVC.CreateConfig(context.Background(), &conf)
		}
		xlog.Info(xlog.DatabaseUserInit)
		for _, ca := range models.DefaultAuth {
			_ = SVC.CreateAuth(context.Background(), &ca)
		}
	})
}
