package control

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/cmd/client/sdk"
)

type CloudProviderController struct {
	svc service.ApiService
}

func (cp CloudProviderController) ListCloudProviders(c *gin.Context) {
	var ok bool
	var count int64
	var res = make([]*models.CloudProviderApi, 0)
	var data []models.CloudProvider
	id, err := getPathId(c)
	page, size := getPageSize(c)
	if err != nil {
		count = cp.svc.Count()
		if count < 0 {
			c.JSON(http.StatusOK, successMsg(count, res))
			return
		}
		data, ok = cp.svc.List(page, size, true).([]models.CloudProvider)
	} else {
		count = cp.svc.Count(service.Condition{
			Key:   "ID",
			Value: strconv.Itoa(id),
		})
		if count < 0 {
			c.JSON(http.StatusOK, successMsg(count, res))
			return
		}
		data, ok = cp.svc.List(page, size, true, service.Condition{
			Key:   "ID",
			Value: strconv.Itoa(id),
		}).([]models.CloudProvider)
	}
	if !ok {
		c.JSON(http.StatusBadRequest, errorMsg(10004))
	}
	// 处理 API
	for _, d := range data {
		api := models.ToApi(d, &models.CloudProviderApi{}, d.Extra())
		res = append(res, api.(*models.CloudProviderApi))
	}
	c.JSON(http.StatusOK, successMsg(count, res))
	return
}

func (cp CloudProviderController) ListActiveCloudProviders(c *gin.Context) {
	var ok bool
	var count int64
	var res = make([]*models.CloudProviderApi, 0)
	var data []models.CloudProvider
	page, size := getPageSize(c)
	count = cp.svc.Count(service.Condition{
		Key:   "Status",
		Value: types.SUCCESS,
	})
	if count < 0 {
		c.JSON(http.StatusOK, successMsg(count, res))
		return
	}
	data, ok = cp.svc.List(page, size, true,
		service.Condition{
			Key:   "Status",
			Value: types.SUCCESS,
		}).([]models.CloudProvider)
	if !ok {
		c.JSON(http.StatusBadRequest, errorMsg(10004))
	}
	for _, d := range data {
		api := models.ToApi(d, &models.CloudProviderApi{}, d.Extra())
		res = append(res, api.(*models.CloudProviderApi))
	}
	c.JSON(http.StatusOK, successMsg(count, res))
	return
}

func (cp CloudProviderController) GetCloudProvider(c *gin.Context) {
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	data, ok := cp.svc.GetById(uint(id)).(models.CloudProvider)
	if !ok {
		c.JSON(http.StatusBadRequest, errorMsg(10004))
	}
	c.JSON(http.StatusOK, successMsg(1, models.ToApi(data, &models.CloudProviderApi{}, data.Extra())))
}

func (cp CloudProviderController) CreateCloudProvider(c *gin.Context) {
	var obj models.CloudProviderCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	// 去重
	if exist := service.Exist(cp.svc, service.Condition{
		Key:   "NAME",
		Value: obj.Name,
	}); exist {
		c.JSON(http.StatusBadRequest, errorMsg(10001))
		return
	}

	target := &models.CloudProvider{}
	// 转换成对应的结构
	models.ToModel(obj, target)
	// 填充 null 字段
	models.AutoFull(target)
	target = cp.svc.Create(target).(*models.CloudProvider)
	cp.sync(target.ID, c)
}

func (cp CloudProviderController) UpdateCloudProvider(c *gin.Context) {
	var obj *models.CloudProviderCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}

	target := &models.CloudProvider{}
	// 转换成对应的结构
	models.ToModel(obj, target)
	target = cp.svc.Update(uint(id), target).(*models.CloudProvider)

	c.JSON(http.StatusOK, successMsg(1, models.ToApi(target, &models.CloudProviderApi{}, target.Extra())))
}

func (cp CloudProviderController) SyncCloudProvider(c *gin.Context) {
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	cp.sync(uint(id), c)
}

func (cp CloudProviderController) DeleteCloudProvider(c *gin.Context) {
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}

	obj := cp.svc.GetById(uint(id))

	if obj == nil {
		c.JSON(http.StatusBadRequest, errorMsg(100002))
		return
	}

	cp.svc.Delete(uint(id))
	c.JSON(http.StatusOK, successMsg(0, nil))
}

func (cp CloudProviderController) sync(id uint, c *gin.Context) {
	target := cp.svc.GetById(id).(*models.CloudProvider)

	err := sdk.GetSDK(*target.Type).Auth(target.ID)
	if err != nil {
		*target.Status = types.AUTH_ERROR
		cp.svc.Update(target.ID, target)
		c.JSON(http.StatusBadRequest, errorMsgInfo(err.Error()))
		return
	}
	// 自动同步函数
	err = sdk.GetSDK(*target.Type).SyncFC(target.ID)
	if err != nil {
		*target.Status = types.SYNC_ERROR
		cp.svc.Update(target.ID, target)
		c.JSON(http.StatusBadRequest, errorMsgInfo(err.Error()))
		return
	}

	// 更新完的
	target = cp.svc.GetById(target.ID).(*models.CloudProvider)
	*target.Status = types.SUCCESS

	cp.svc.Update(target.ID, target)
	c.JSON(http.StatusOK, successMsg(1, models.ToApi(target, &models.CloudProviderApi{}, target.Extra())))
}
