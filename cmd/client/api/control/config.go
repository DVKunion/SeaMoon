package control

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
)

type SysConfigController struct {
	svc service.ApiService
}

func (sc SysConfigController) ListSystemConfigs(c *gin.Context) {
	page, size := getPageSize(c)

	var obj = sc.svc.List(page, size, false).([]*models.SystemConfig)

	c.JSON(http.StatusOK, successMsg(1, models.ToSystemConfigApi(obj)))
}

func (sc SysConfigController) UpdateSystemConfig(c *gin.Context) {
	var obj models.SystemConfigApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	// model 转换
	target := obj.ToModel()
	sc.svc.Update(0, target)
	c.JSON(http.StatusOK, successMsg(1, obj))
}
