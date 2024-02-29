package control

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
)

var (
	StatisticC = StatisticController{}
	AuthC      = AuthController{
		service.GetService("auth"),
	}
	ProxyC = ProxyController{
		service.GetService("proxy"),
		service.GetService("tunnel"),
	}
	TunnelC = TunnelController{
		service.GetService("tunnel"),
		service.GetService("provider"),
	}
	ProviderC = CloudProviderController{
		service.GetService("provider"),
	}
	SysConfigC = SysConfigController{
		service.GetService("config"),
	}
)

// PageNotFound 404页面
func PageNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, errorMsg(10404))
}

// 通用获取页面信息
func getPageSize(c *gin.Context) (int, int) {
	defaultPageSize := 10
	defaultPageNum := 0

	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, errPage := strconv.Atoi(pageStr)
	if errPage != nil {
		page = defaultPageNum
	}

	size, errSize := strconv.Atoi(sizeStr)
	if errSize != nil {
		size = defaultPageSize
	}
	return page, size
}

// 通用获取 path 信息
func getPathId(c *gin.Context) (int, error) {
	return strconv.Atoi(c.Param("id"))
}

// 通用 正常响应
func successMsg(total int64, data interface{}) gin.H {
	return gin.H{
		"success": true,
		"total":   total,
		"data":    data,
	}
}

var errorCodeMaps = map[int]string{
	10000: "请求结构错误",
	10001: "请求数据重复",
	10002: "请求数据不存在",
	10003: "请求数据字段缺失",
	10004: "请求数据异常",

	10010: "认证失败",
	10011: "需要认证",
	10012: "权限不足",
	10013: "Limit限制",

	10404: "页面不存在",
}

// 通用 错误响应
func errorMsg(code int) gin.H {

	if msg, ok := errorCodeMaps[code]; ok {
		return gin.H{
			"success": false,
			"code":    code,
			"msg":     msg,
		}
	}

	return gin.H{
		"type": "error",
		"code": code,
		"msg":  "unknown error",
	}
}

func errorMsgInfo(m string) gin.H {

	return gin.H{
		"success": false,
		"code":    99999,
		"msg":     m,
	}
}
