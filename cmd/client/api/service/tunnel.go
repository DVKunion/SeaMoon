package service

import (
	"github.com/DVKunion/SeaMoon/cmd/client/api/database"
	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
)

type TunnelService struct {
}

func (t TunnelService) Count(cond ...Condition) int64 {
	var res int64 = 0
	conn := database.GetConn()
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Model(&models.Tunnel{}).Count(&res)
	return res
}

func (t TunnelService) List(page int, size int, preload bool, cond ...Condition) interface{} {
	var data = make([]models.Tunnel, 0)
	conn := database.QueryPage(page, size)
	if preload {
		conn = conn.Preload("Proxies")
	}
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Find(&data)
	return data
}

func (t TunnelService) GetById(id uint) interface{} {
	var data = models.Tunnel{}
	database.GetConn().Limit(1).Where("ID = ?", id).Find(&data)
	return &data
}

func (t TunnelService) Create(obj interface{}) interface{} {
	database.GetConn().Create(obj.(*models.Tunnel))
	return obj
}

func (t TunnelService) Update(id uint, obj interface{}) interface{} {
	obj.(*models.Tunnel).ID = id
	database.GetConn().Updates(obj.(*models.Tunnel))
	return t.GetById(id)
}

func (t TunnelService) Delete(id uint) {
	database.GetConn().Delete(&models.Tunnel{}, id)
}
