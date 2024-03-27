package v1

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/pkg/api/controller/servant"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func ListTunnels(c *gin.Context) {
	total, err := service.SVC.TotalTunnels(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
		return
	}

	p, s := servant.GetPageSize(c)
	if res, err := service.SVC.ListTunnels(c, p, s); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, total, res.ToApi(extra()))
	}
}

func GetTunnelById(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}
	if res, err := service.SVC.GetTunnelById(c, uint(id)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi(extra()))
	}
}

func CreateTunnel(c *gin.Context) {
	var obj models.TunnelCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	if service.SVC.ExistTunnel(c, obj.Name, nil) {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsExist, nil))
		return
	}

	if res, err := service.SVC.CreateTunnel(c, obj.ToModel(true)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi(extra()))
	}
}

func UpdateTunnel(c *gin.Context) {
	var obj *models.TunnelCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	obj.ID = uint(id)

	if res, err := service.SVC.UpdateTunnel(c, obj.ToModel(false)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi(extra()))
	}
}

func DeleteTunnel(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	if err = service.SVC.DeleteTunnel(c, uint(id)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, nil)
	}
}

func extra() func(api interface{}) {
	return func(api interface{}) {
		ref := reflect.ValueOf(api).Elem()
		idField := ref.FieldByName("ProviderId")
		prv, err := service.SVC.GetProviderById(context.Background(), uint(idField.Uint()))
		if err != nil {
			// todo: deal with this error
			return
		}
		field := ref.FieldByName("ProviderType")
		if field.CanSet() {
			field.Set(reflect.ValueOf(prv.Type))
		}
	}
}
