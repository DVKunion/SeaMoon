package client

import (
	"context"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "net/http/pprof"

	"github.com/DVKunion/SeaMoon/cmd/client/static"
	"github.com/DVKunion/SeaMoon/pkg/consts"
)

func Serve(ctx context.Context, verbose bool, debug bool) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sg := NewSigGroup()
	go API(sg, verbose, debug)
	go Control(ctx, sg)

	Config().Load(sg)
	<-sg.WatchChannel

	sg.StopProxy()
	cancel()

	sg.wg.Wait()
}

func API(sg *SigGroup, verbose bool, debug bool) {
	slog.Info(consts.CONTROLLER_START, "addr", Config().Control.ConfigAddr)

	if consts.Version != "dev" || !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	webLogger, err := os.OpenFile(Config().Control.LogPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else if verbose {
		gin.DefaultWriter = io.MultiWriter(webLogger, os.Stdout)
	} else {
		gin.DefaultWriter = io.MultiWriter(webLogger)
	}

	server := gin.Default()

	tmpl := template.Must(template.New("").ParseFS(static.HtmlFiles, "templates/*.html"))
	server.SetHTMLTemplate(tmpl)

	server.StaticFS("/static", http.FS(static.AssetFiles))
	server.StaticFileFS("/favicon.ico", "public/img/favicon.ico", http.FS(static.AssetFiles))

	// controller page
	server.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", Config())
	})

	// pprof
	if debug {
		server.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	}

	// controller set
	server.POST("/", func(ctx *gin.Context) {
		if err := ctx.ShouldBindJSON(Config()); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg":    "服务异常",
				"result": err.Error(),
			})
			return
		}
		sg.Detection()
		ctx.JSON(http.StatusOK, Config())
	})

	if err := server.Run(Config().Control.ConfigAddr); err != http.ErrServerClosed {
		slog.Error("client running error", "err", err)
	}
}
