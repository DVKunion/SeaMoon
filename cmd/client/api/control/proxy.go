package control

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/cmd/client/signal"
	"github.com/DVKunion/SeaMoon/pkg/tunnel"
)

type ProxyController struct {
	svc  service.ApiService
	tsvc service.ApiService
}

func (pc ProxyController) ListProxies(c *gin.Context) {
	var res = make([]*models.ProxyApi, 0)
	var data []models.Proxy
	var ok bool
	var count int64
	id, err := getPathId(c)
	page, size := getPageSize(c)
	if err != nil {
		count = pc.svc.Count()
		if count < 0 {
			c.JSON(http.StatusOK, successMsg(count, res))
			return
		}
		data, ok = pc.svc.List(page, size, true).([]models.Proxy)
	} else {
		count = pc.svc.Count(service.Condition{
			Key:   "ID",
			Value: strconv.Itoa(id),
		})
		if count < 0 {
			c.JSON(http.StatusOK, successMsg(count, res))
			return
		}
		data, ok = pc.svc.List(page, size, true, service.Condition{
			Key:   "ID",
			Value: strconv.Itoa(id),
		}).([]models.Proxy)
	}
	if !ok {
		c.JSON(http.StatusBadRequest, errorMsg(10004))
	}
	// 处理 API
	for _, d := range data {
		api := models.ToApi(d, &models.ProxyApi{})
		res = append(res, api.(*models.ProxyApi))
	}

	c.JSON(http.StatusOK, successMsg(count, res))
	return
}

func (pc ProxyController) GetProxy(c *gin.Context) {
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	proxy := pc.svc.GetById(uint(id)).(models.Proxy)
	c.JSON(http.StatusOK, successMsg(1, models.ToApi(proxy, &models.ProxyApi{})))
}

func (pc ProxyController) CreateProxy(c *gin.Context) {
	var obj *models.ProxyCreateApi
	var proxy *models.Proxy
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	// 先查询一下是否存在
	if exist := service.Exist(pc.svc, service.Condition{
		Key:   "NAME",
		Value: obj.Name,
	}, service.Condition{
		Key:   "tunnel_id",
		Value: obj.TunnelId,
	}); exist {
		c.JSON(http.StatusBadRequest, errorMsg(10001))
		return
	}

	// 判断是从账户还是从函数
	if obj.TunnelId != 0 {
		// 去查询队列
		tun := pc.tsvc.GetById(obj.TunnelId).(*models.Tunnel)
		if tun != nil && tun.ID != 0 {
			// 说明没问题，继续
			target := &models.Proxy{}
			// 转换成对应的结构
			models.ToModel(obj, target)
			// 填充 null 字段
			models.AutoFull(target)
			proxy = pc.svc.Create(target).(*models.Proxy)

			tun.Proxies = append(tun.Proxies, *proxy)
			pc.tsvc.Update(tun.ID, tun)
		}
	}

	// 判断是否从账户来的，账户来的需要先创建 tunnel
	if obj.TunnelCreateApi != nil {
		// 说明没问题，继续
		target := &models.Proxy{}
		// 转换成对应的结构
		models.ToModel(obj, target)
		// 填充 null 字段
		models.AutoFull(target)
		proxy = pc.svc.Create(target).(*models.Proxy)
		// 创建 tunnel
		tun, err := TunnelC.createTunnel(obj.TunnelCreateApi)
		if err != nil {
			*proxy.Status = types.ERROR
			pc.tsvc.Update(target.ID, target)
			c.JSON(http.StatusBadRequest, errorMsgInfo(err.Error()))
			return
		}
		// 再加入 proxy
		tun.Proxies = append(tun.Proxies, *proxy)
		*tun.Status = tunnel.ACTIVE
		pc.tsvc.Update(tun.ID, tun)
	}

	if proxy == nil {
		c.JSON(http.StatusBadRequest, errorMsg(10002))
		return
	}

	// 发送至队列
	signal.Signal().SendStartProxy(proxy.ID, proxy.Addr())

	c.JSON(http.StatusOK, successMsg(1, models.ToApi(proxy, &models.ProxyApi{})))
	return
}

func (pc ProxyController) UpdateProxy(c *gin.Context) {
	var obj *models.ProxyCreateApi
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	data := pc.svc.GetById(uint(id)).(*models.Proxy)
	target := &models.Proxy{}
	// 转换成对应的结构
	models.ToModel(obj, target)
	// 说明服务状态发生了改变
	if *target.Status != *data.Status {
		switch *target.Status {
		case types.ACTIVE:
			signal.Signal().SendStartProxy(data.ID, data.Addr())
		case types.INACTIVE:
			signal.Signal().SendStopProxy(data.ID, data.Addr())
		}
	}
	target = pc.svc.Update(uint(id), target).(*models.Proxy)

	c.JSON(http.StatusOK, successMsg(1, models.ToApi(target, &models.ProxyApi{})))
}

func (pc ProxyController) DeleteProxy(c *gin.Context) {
	id, err := getPathId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}

	obj := pc.svc.GetById(uint(id))

	if obj == nil {
		c.JSON(http.StatusBadRequest, errorMsg(100002))
		return
	}

	pc.svc.Delete(uint(id))
	c.JSON(http.StatusOK, successMsg(1, nil))
}
