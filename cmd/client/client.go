package client

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/api/signal"
	"github.com/DVKunion/SeaMoon/cmd/client/static"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/xlog"
)

func Serve(ctx context.Context, debug bool) {
	go signal.Signal().Run(ctx)
	// Restful API 服务
	runApi(debug)
}

func runApi(debug bool) {
	logPath := service.GetService("config").(service.SystemConfigService).GetByName("control_log").Value
	addr := service.GetService("config").(service.SystemConfigService).GetByName("control_addr").Value
	port := service.GetService("config").(service.SystemConfigService).GetByName("control_port").Value

	xlog.Info("API", xlog.CONTROLLER_START, "addr", addr, "port", port)

	if consts.Version != "dev" || !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	webLogger, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else {
		gin.DefaultWriter = io.MultiWriter(webLogger)
	}

	server := gin.Default()

	api.Register(server, debug)

	subFS, err := fs.Sub(static.WebViews, "dist")

	if err != nil {
		panic("web static file error: " + err.Error())
	}

	server.NoRoute(func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(subFS))
	})

	if err := server.Run(strings.Join([]string{addr, port}, ":")); err != http.ErrServerClosed {
		slog.Error("client running error", "err", err)
	}
}
