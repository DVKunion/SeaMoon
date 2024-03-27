package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/pkg/api/controller/servant"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func Login(c *gin.Context) {
	var obj *models.AuthApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	if token, err := service.SVC.Login(c, obj); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, token)
	}
}

func Passwd(c *gin.Context) {
	var obj *models.AuthApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	if err := service.SVC.UpdatePassword(c, obj); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, "更新成功")
	}
}
