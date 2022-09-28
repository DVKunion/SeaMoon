package client

import (
	"bufio"
	"context"
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
)

type Client struct {
	net.Conn
	br *bufio.Reader
}

func (c *Client) Read(b []byte) (int, error) {
	return c.br.Read(b)
}

func Serve(ctx context.Context, verbose bool, debug bool) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sg := NewSigGroup()
	go Controller(sg, verbose, debug)
	go HttpController(ctx, sg)
	go Socks5Controller(ctx, sg)

	Config().Load(sg)
	<-sg.WatchChannel

	sg.StopProxy()
	cancel()

	sg.wg.Wait()
}

func Controller(sg *SigGroup, verbose bool, debug bool) {
	log.Infof(consts.CONTROLLER_START, Config().Control.ConfigAddr)

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
		log.Error(err)
	}
}
