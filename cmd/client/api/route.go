package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api/control"
	"github.com/DVKunion/SeaMoon/cmd/client/api/database"
)

func init() {
	database.Init()
}

func Register(router *gin.Engine, debug bool) {
	// pprof
	if debug {
		router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	}
	// statistic
	router.GET("/api/statistic", control.JWTAuthMiddleware(debug), control.StatisticC.Get)
	// user
	router.POST("/api/user/login", control.AuthC.Login)
	router.PUT("/api/user/passwd", control.JWTAuthMiddleware(debug), control.AuthC.Passwd)
	// proxy
	router.GET("/api/proxy", control.JWTAuthMiddleware(debug), control.ProxyC.ListProxies)
	router.GET("/api/proxy/:id", control.JWTAuthMiddleware(debug), control.ProxyC.ListProxies)
	router.POST("/api/proxy", control.JWTAuthMiddleware(debug), control.ProxyC.CreateProxy)
	router.PUT("/api/proxy/:id", control.JWTAuthMiddleware(debug), control.ProxyC.UpdateProxy)
	router.DELETE("/api/proxy/:id", control.JWTAuthMiddleware(debug), control.ProxyC.DeleteProxy)
	// tunnel
	router.GET("/api/tunnel", control.JWTAuthMiddleware(debug), control.TunnelC.ListTunnels)
	router.GET("/api/tunnel/:id", control.JWTAuthMiddleware(debug), control.TunnelC.ListTunnels)
	router.POST("/api/tunnel", control.JWTAuthMiddleware(debug), control.TunnelC.CreateTunnel)
	router.PUT("/api/tunnel/:id", control.JWTAuthMiddleware(debug), control.TunnelC.UpdateTunnel)
	router.DELETE("/api/tunnel/:id", control.JWTAuthMiddleware(debug), control.TunnelC.DeleteTunnel)
	// cloud provider
	router.GET("/api/provider", control.JWTAuthMiddleware(debug), control.ProviderC.ListCloudProviders)
	router.GET("/api/provider/active", control.JWTAuthMiddleware(debug), control.ProviderC.ListActiveCloudProviders)
	router.GET("/api/provider/:id", control.JWTAuthMiddleware(debug), control.ProviderC.ListCloudProviders)
	router.POST("/api/provider", control.JWTAuthMiddleware(debug), control.ProviderC.CreateCloudProvider)
	router.PUT("/api/provider/sync/:id", control.JWTAuthMiddleware(debug), control.ProviderC.SyncCloudProvider)
	router.PUT("/api/provider/:id", control.JWTAuthMiddleware(debug), control.ProviderC.UpdateCloudProvider)
	router.DELETE("/api/provider/:id", control.JWTAuthMiddleware(debug), control.ProviderC.DeleteCloudProvider)
	// system config
	router.GET("/api/config", control.JWTAuthMiddleware(debug), control.SysConfigC.ListSystemConfigs)
	router.PUT("/api/config/", control.JWTAuthMiddleware(debug), control.SysConfigC.UpdateSystemConfig)
}
