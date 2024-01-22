package service

import (
	"github.com/DVKunion/SeaMoon/cmd/client/api/database"
	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
)

type ProxyService struct {
}

func (p ProxyService) Count(cond ...Condition) int64 {
	var res int64 = 0
	conn := database.GetConn()
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Model(&models.Proxy{}).Count(&res)
	return res
}

func (p ProxyService) List(page, size int, preload bool, cond ...Condition) interface{} {
	var data = make([]models.Proxy, 0)
	conn := database.QueryPage(page, size)
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Find(&data)
	return data
}

func (p ProxyService) GetById(id uint) interface{} {
	var data = models.Proxy{}
	database.GetConn().Limit(1).Where("ID = ?", id).Find(&data)
	return &data
}

func (p ProxyService) Create(obj interface{}) interface{} {
	database.GetConn().Create(obj.(*models.Proxy))
	return obj
}

func (p ProxyService) Update(id uint, obj interface{}) interface{} {
	obj.(*models.Proxy).ID = id
	database.GetConn().Updates(obj.(*models.Proxy))
	return p.GetById(id)
}

func (p ProxyService) Delete(id uint) {
	database.GetConn().Delete(&models.Proxy{}, id)
}
