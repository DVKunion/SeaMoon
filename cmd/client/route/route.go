package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/pkg/api/controller/middleware"
	api_v1 "github.com/DVKunion/SeaMoon/pkg/api/controller/v1"
	"github.com/DVKunion/SeaMoon/pkg/api/database/drivers"
)

func init() {
	drivers.Init()
}

func Register(router *gin.Engine, debug bool) {
	var middles = make([]gin.HandlerFunc, 0)

	if debug {
		router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	} else {
		middles = append(middles, middleware.JWTAuthMiddleware)
	}

	registerV1(router, middles)
}

func registerV1(router *gin.Engine, middles []gin.HandlerFunc) {
	v1 := router.Group("/api/v1")
	// user
	v1.POST("/user/login", api_v1.Login)
	v1.PUT("/user/passwd", append(middles, api_v1.Passwd)...)

	// proxy
	v1.GET("/proxy", append(middles, api_v1.ListProxies)...)
	v1.GET("/proxy/:id", append(middles, api_v1.GetProxyById)...)
	v1.GET("/proxy/speed/:id", append(middles, api_v1.SpeedRateProxy)...)
	v1.POST("/proxy", append(middles, api_v1.CreateProxy)...)
	v1.PUT("/proxy/:id", append(middles, api_v1.UpdateProxy)...)
	v1.DELETE("/proxy/:id", append(middles, api_v1.DeleteProxy)...)

	// tunnel
	v1.GET("/tunnel", append(middles, api_v1.ListTunnels)...)
	v1.GET("/tunnel/:id", append(middles, api_v1.GetTunnelById)...)
	v1.POST("/tunnel", append(middles, api_v1.CreateTunnel)...)
	v1.PUT("/tunnel/:id", append(middles, api_v1.UpdateTunnel)...)
	v1.DELETE("/tunnel/:id", append(middles, api_v1.DeleteTunnel)...)

	// cloud provider
	v1.GET("/provider", append(middles, api_v1.ListProviders)...)
	v1.GET("/provider/active", append(middles, api_v1.ListActiveProviders)...)
	v1.GET("/provider/:id", append(middles, api_v1.GetProviderById)...)
	v1.POST("/provider", append(middles, api_v1.CreateProvider)...)
	v1.PUT("/provider/sync/:id", append(middles, api_v1.SyncProvider)...)
	v1.PUT("/provider/:id", append(middles, api_v1.UpdateProvider)...)
	v1.DELETE("/provider/:id", append(middles, api_v1.DeleteProvider)...)

	// config
	v1.GET("/config", append(middles, api_v1.ListConfigs)...)
	v1.PUT("/config", append(middles, api_v1.UpdateConfig)...)
}
