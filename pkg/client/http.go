package client

import (
	"github.com/DVKunion/SeaMoon/pkg/consts"
	"github.com/DVKunion/SeaMoon/pkg/utils"
	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

func NewHttpClient(listenAddr string, proxyAddr string, verbose bool) {
	server := goproxy.NewProxyHttpServer()
	if err := InitCa(); err != nil {
		log.Error(consts.CA_ERROR, err)
	}
	server.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	server.Verbose = verbose || log.GetLevel() == log.DebugLevel
	server.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		req.Header.Set("SM-Host", utils.GetUrl(req))
		req.URL, _ = url.Parse(proxyAddr)
		req.Host = req.URL.Host
		return req, nil
	})

	log.Infof(consts.HTTP_LISTEN_START, listenAddr)

	if err := http.ListenAndServe(listenAddr, server); err != nil {
		log.Error(consts.HTTP_LISTEN_ERROR, err)
	}
}
