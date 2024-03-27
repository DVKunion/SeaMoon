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
)

func ListProviders(c *gin.Context) {
	total, err := service.SVC.TotalProviders(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
		return
	}

	p, s := servant.GetPageSize(c)
	if res, err := service.SVC.ListProviders(c, p, s); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, total, res.ToApi())
	}
}

func GetProviderById(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}
	if res, err := service.SVC.GetProviderById(c, uint(id)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func ListActiveProviders(c *gin.Context) {
	total, err := service.SVC.TotalProviders(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
		return
	}

	if res, err := service.SVC.ListActiveProviders(c); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		servant.SuccessMsg(c, total, res.ToApi())
	}
}

func CreateProvider(c *gin.Context) {
	var obj models.ProviderCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	if service.SVC.ExistProvider(c, obj.Name) {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsExist, nil))
		return
	}

	if res, err := service.SVC.CreateProvider(c, obj.ToModel(true)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		signal.Signal().SendProviderSignal(res.ID, enum.ProvStatusSync, nil)
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func UpdateProvider(c *gin.Context) {
	var obj *models.ProviderCreateApi
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

	if res, err := service.SVC.UpdateProvider(c, obj.ToModel(false)); err != nil {
		servant.ErrorMsg(c, http.StatusInternalServerError, errors.ApiError(xlog.ApiServiceError, err))
	} else {
		signal.Signal().SendProviderSignal(res.ID, enum.ProvStatusSync, nil)
		servant.SuccessMsg(c, 1, res.ToApi())
	}
}

func DeleteProvider(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	signal.Signal().SendProviderSignal(uint(id), enum.ProvStatusDelete, wg)
	wg.Wait()
	servant.SuccessMsg(c, 1, nil)
}

func SyncProvider(c *gin.Context) {
	id, err := servant.GetPathId(c)
	if err != nil {
		servant.ErrorMsg(c, http.StatusBadRequest, errors.ApiError(xlog.ApiParamsError, err))
		return
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	signal.Signal().SendProviderSignal(uint(id), enum.ProvStatusSync, wg)
	wg.Wait()
	servant.SuccessMsg(c, 1, nil)
}
