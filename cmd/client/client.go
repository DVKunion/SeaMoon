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
	"github.com/DVKunion/SeaMoon/pkg/api/signal"
	"github.com/DVKunion/SeaMoon/pkg/system/version"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func Serve(ctx context.Context, debug bool) {
	// Signal 异步服务
	runSignal(ctx)
	// Restful API 服务
	runApi(ctx, debug)
}

func runSignal(ctx context.Context) {
	// 控制总线，用于管控服务相关
	go signal.Signal().Daemon(ctx)
	// 如果配置了自动恢复设置，尝试发送恢复信号
	rec, err := service.SVC.GetConfigByName(ctx, "auto_start")
	if err != nil {
		xlog.Error(xlog.SignalGetObjError, "err", err)
		return
	}
	signal.Signal().Recover(ctx, rec.Value)

	// 启动时同步云账户和执行健康检查（异步执行，不阻塞启动）
	go signal.Signal().StartupSync(ctx)
}

func runApi(ctx context.Context, debug bool) {
	logPath, err := service.SVC.GetConfigByName(ctx, "control_log")
	addr, err := service.SVC.GetConfigByName(ctx, "control_addr")
	port, err := service.SVC.GetConfigByName(ctx, "control_port")

	xlog.Info(xlog.ApiServiceStart, "addr", addr.Value, "port", port.Value)

	if version.Version != "dev" || !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	webLogger, err := os.OpenFile(logPath.Value, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		gin.DefaultWriter = io.MultiWriter(xlog.Logger())
	} else {
		gin.DefaultWriter = io.MultiWriter(xlog.Logger(), webLogger)
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
		xlog.Error(xlog.ApiServeError, "err", err)
	}
}
