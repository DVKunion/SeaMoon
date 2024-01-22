package control

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/sdk"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type TunnelController struct {
	svc  service.ApiService
	ksvc service.ApiService
}

func (tc TunnelController) ListTunnels(c *gin.Context) {
	var res = make([]interface{}, 0)
	var data []models.Tunnel
	var ok bool
	var count int64
	id, err := getPathId(c)
	page, size := getPageSize(c)
	if err != nil {
		count = tc.svc.Count()
		if count < 0 {
			c.JSON(http.StatusOK, successMsg(count, res))
			return
		}
		data, ok = tc.svc.List(page, size, true).([]models.Tunnel)
	} else {
		count = tc.svc.Count(service.Condition{
			Key:   "ID",
			Value: strconv.Itoa(id),
		})
		if count < 0 {
			c.JSON(http.StatusOK, successMsg(count, res))
			return
		}
		data, ok = tc.svc.List(page, size, true, service.Condition{
			Key:   "ID",
			Value: strconv.Itoa(id),
		}).([]models.Tunnel)
	}
	if !ok {
		c.JSON(http.StatusBadRequest, errorMsg(10004))
	}
	// 处理 API
	for _, d := range data {
		prd := tc.ksvc.GetById(d.CloudProviderId).(*models.CloudProvider)
		api := models.ToApi(d, &models.TunnelApi{}, func(api interface{}) {
			ref := reflect.ValueOf(api).Elem()
			field := ref.FieldByName("CloudProviderType")
			if field.CanSet() {
				field.Set(reflect.ValueOf(prd.Type))
			}
			field = ref.FieldByName("CloudProviderRegion")
			if field.CanSet() {
				field.Set(reflect.ValueOf(prd.Region))
			}
		})
		res = append(res, api.(*models.TunnelApi))
	}
	c.JSON(http.StatusOK, successMsg(count, res))
	return
}

func (tc TunnelController) GetTunnel(c *gin.Context) {
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	c.JSON(http.StatusOK, successMsg(1, tc.svc.GetById(uint(id))))
}

func (tc TunnelController) CreateTunnel(c *gin.Context) {
	var obj models.TunnelCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	// 先查询一下是否存在
	if exist := service.Exist(tc.svc, service.Condition{
		Key:   "NAME",
		Value: obj.Name,
	}, service.Condition{
		Key:   "cloud_provider_id",
		Value: obj.CloudProviderId,
	}); exist {
		c.JSON(http.StatusBadRequest, errorMsg(10001))
		return
	}

	target, err := tc.createTunnel(&obj)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsgInfo(err.Error()))
		return
	}

	c.JSON(http.StatusOK, successMsg(1, models.ToApi(target, &models.TunnelApi{})))
}

func (tc TunnelController) UpdateTunnel(c *gin.Context) {
	var obj *models.TunnelCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}

	target := &models.Tunnel{}
	// 转换成对应的结构
	models.ToModel(obj, target)
	target = tc.svc.Update(uint(id), target).(*models.Tunnel)

	c.JSON(http.StatusOK, successMsg(1, models.ToApi(target, &models.TunnelApi{})))
}

func (tc TunnelController) DeleteTunnel(c *gin.Context) {
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}

	obj := tc.svc.GetById(uint(id))

	if obj == nil {
		c.JSON(http.StatusBadRequest, errorMsg(100002))
		return
	}

	tc.svc.Delete(uint(id))
	c.JSON(http.StatusOK, successMsg(0, nil))
}

func (tc TunnelController) createTunnel(tun *models.TunnelCreateApi) (*models.Tunnel, error) {
	target := &models.Tunnel{}

	// 检查账户provider是否正确
	prv := tc.ksvc.GetById(tun.CloudProviderId).(*models.CloudProvider)
	if prv == nil || prv.ID == 0 {
		return target, errors.New("provider 数据不存在")
	}
	if *prv.MaxLimit != 0 && *prv.MaxLimit < len(prv.Tunnels)+1 {
		return target, errors.New("超出最大 limit 限制")
	}

	// 转换成对应的结构
	models.ToModel(tun, target)
	// 填充默认值
	models.AutoFull(target)
	target.CloudProviderId = 0
	target = tc.svc.Create(target).(*models.Tunnel)

	// 添加到 provider 中去
	prv.Tunnels = append(prv.Tunnels, *target)
	tc.ksvc.Update(prv.ID, prv)

	err := sdk.GetSDK(*prv.Type).Deploy(prv.ID, target)
	if err != nil {
		// 部署失败了，更新状态
		*target.Status = tunnel.ERROR
		tc.svc.Update(target.ID, target)
		return target, err
	}
	return target, nil
}
