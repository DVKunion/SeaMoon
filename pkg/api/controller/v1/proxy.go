package v1

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/pkg/api/controller/servant"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/api/service"
	"github.com/DVKunion/SeaMoon/pkg/signal"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
	"github.com/DVKunion/SeaMoon/pkg/tools"
)

func ListProxies(c *gin.Context) {
	total, err := service.SVC.TotalProxies(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
		return
	}

	p, s := servant.GetPageSize(c)
	if res, err := service.SVC.ListProxies(c, p, s); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, total, res.ToApi())
	}
}

func GetProxyById(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}
	if res, err := service.SVC.GetProxyById(c, uint(id)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func CreateProxy(c *gin.Context) {
	var obj models.ProxyCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	if service.SVC.ExistProvider(c, obj.Name) {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsExist, nil))
		return
	}

	// 判断是否从账户来的，账户来的需要先创建 tunnel
	if obj.TunnelCreateApi != nil {
		if tun, err := service.SVC.CreateTunnel(c, obj.TunnelCreateApi.ToModel(true)); err != nil {
			*obj.Status = enum.ProxyStatusError
			*obj.StatusMessage = err.Error()
		} else {
			signal.Signal().SendTunnelSignal(tun.ID, enum.TunnelActive, nil)
			obj.TunnelID = tun.ID
		}
	}

	// 非账户来的，直接带着
	if res, err := service.SVC.CreateProxy(c, obj.ToModel(true)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
		return
	} else {
		// 发送启动通知
		signal.Signal().SendProxySignal(res.ID, enum.ProxyStatusActive, nil)
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func UpdateProxy(c *gin.Context) {
	var obj *models.ProxyCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	m := obj.ToModel(false)
	if m.Status != nil {
		signal.Signal().SendProxySignal(obj.ID, *obj.Status, nil)
		// 这里愚蠢de做一个特殊处理: 当代理关闭时，自动将连接数清零
		if *m.Status == enum.ProxyStatusInactive {
			m.Conn = tools.IntPtr(0)
		}
	}
	if res, err := service.SVC.UpdateProxy(c, uint(id), m); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func DeleteProxy(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	signal.Signal().SendProxySignal(uint(id), enum.ProxyStatusDelete, wg)
	wg.Wait()
	servant.SuccessMsg(c, 1, nil)
}

func SpeedRateProxy(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	// todo 全量测速
	if proxy, err := service.SVC.GetProxyById(c, uint(id)); err != nil || *proxy.Status != enum.ProxyStatusActive {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		signal.Signal().SendProxySignal(proxy.ID, enum.ProxyStatusSpeeding, nil)
		servant.SuccessMsg(c, 1, proxy.ToApi())
	}
}
