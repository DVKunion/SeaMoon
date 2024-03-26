package client

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/route"
	"github.com/DVKunion/SeaMoon/cmd/client/static"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/signal"
	"github.com/DVKunion/SeaMoon/pkg/system/consts"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func Serve(ctx context.Context, debug bool) {
	// 控制总线，用于管控服务相关
	go signal.Signal().Daemon(ctx)
	// Restful API 服务
	runApi(ctx, debug)
}

func runApi(ctx context.Context, debug bool) {
	logPath, err := service.SVC.GetConfigByName(ctx, "control_log")
	addr, err := service.SVC.GetConfigByName(ctx, "control_addr")
	port, err := service.SVC.GetConfigByName(ctx, "control_port")

	xlog.Info(xlog.ApiServiceStart, "addr", addr.Value, "port", port.Value)

	if consts.Version != "dev" || !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	webLogger, err := os.OpenFile(logPath.Value, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else {
		gin.DefaultWriter = io.MultiWriter(webLogger)
	}

	server := gin.Default()

	route.Register(server, debug)

	subFS, err := fs.Sub(static.WebViews, "dist")

	if err != nil {
		panic("web static file error: " + err.Error())
	}

	server.NoRoute(func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(subFS))
	})

	if err := server.Run(strings.Join([]string{addr.Value, port.Value}, ":")); err != http.ErrServerClosed {
		xlog.Error(errors.ApiServeError, "err", err)
	}
}
