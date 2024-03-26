package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/pkg/api/controller/servant"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/signal"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
)

func ListProxies(c *gin.Context) {
	total, err := service.SVC.TotalProxies(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
		return
	}

	p, s := servant.GetPageSize(c)
	if res, err := service.SVC.ListProxies(c, p, s); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, total, res.ToApi())
	}
}

func GetProxyById(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(errors.ApiParamsError, err))
		return
	}
	if res, err := service.SVC.GetProxyById(c, uint(id)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func CreateProxy(c *gin.Context) {
	var obj models.ProxyCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(errors.ApiParamsError, err))
		return
	}

	if service.SVC.ExistProvider(c, obj.Name) {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(errors.ApiParamsExist, nil))
		return
	}

	// 判断是否从账户来的，账户来的需要先创建 tunnel
	if obj.TunnelCreateApi != nil {
		if tun, err := service.SVC.CreateTunnel(c, obj.TunnelCreateApi.ToModel(true)); err != nil {
			*obj.Status = enum.ProxyStatusError
			*obj.StatusMessage = err.Error()
		} else {
			// todo: 发送部署请求
			obj.TunnelId = tun.ID
		}
	}

	if res, err := service.SVC.CreateProxy(c, obj.ToModel(true)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
		return
	} else {
		// 发送启动通知
		signal.Signal().SendProxySignal(res.ID, *res.Status)
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func UpdateProxy(c *gin.Context) {
	var obj *models.ProxyCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(errors.ApiParamsError, err))
		return
	}
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(errors.ApiParamsError, err))
		return
	}

	// 朝着队列发送控制信号
	signal.Signal().SendProxySignal(obj.ID, *obj.Status)

	if res, err := service.SVC.UpdateProxy(c, uint(id), obj.ToModel(false)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func DeleteProxy(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(errors.ApiParamsError, err))
		return
	}

	if err = service.SVC.DeleteProxy(c, uint(id)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, nil)
	}
}

func SpeedRateProxy(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(errors.ApiParamsError, err))
		return
	}

	// todo 全量测速
	if proxy, err := service.SVC.GetProxyById(c, uint(id)); err != nil || *proxy.Status != enum.ProxyStatusActive {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
	} else {
		signal.Signal().SendProxySignal(proxy.ID, enum.ProxyStatusSpeeding)
		if err := service.SVC.UpdateProxyStatus(c, proxy.ID, enum.ProxyStatusSpeeding, ""); err != nil {
			servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(errors.ApiServiceError, err))
		}
		// 发送测速信号
		servant.SuccessMsg(c, 1, proxy.ToApi())
	}
}
