package service

import (
	"context"

	"github.com/DVKunion/SeaMoon/pkg/api/database/drivers"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/models/abstract"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

var SVC *svc

type svc struct {
	*account
	*config
	*provider
	*proxy
	*tunnel
}

func Init() {
	SVC = &svc{
		&account{},
		&config{},
		&provider{},
		&proxy{},
		&tunnel{},
	}
	drivers.RegisterMigrate(func() {
		xlog.Info(xlog.DatabaseConfigInit)
		for _, conf := range defaultConfig {
			_ = SVC.CreateConfig(context.Background(), &conf)
		}
		xlog.Info(xlog.DatabaseUserInit)
		var ca = models.Account{
			Type:     enum.Admin,
			AuthType: enum.AdminAuth,
			Auth: abstract.AdminAuth{
				Username: "seamoon",
				Password: "2575a6f37310dd27e884a0305a2dd210",
			},
		}
		_ = SVC.CreateAuth(context.Background(), &ca)
	})
}

var defaultConfig = []models.Config{
	{
		Key:   "control_addr",
		Value: "0.0.0.0",
	},
	{
		Key:   "control_port",
		Value: "7777",
	},
	{
		Key:   "xray_api_port",
		Value: "10085",
	},
	{
		Key:   "control_log",
		Value: "seamoon.log",
	},
	{
		Key:   "auto_start",
		Value: "true",
	},
	{
		Key:   "auto_sync",
		Value: "true",
	},
}
