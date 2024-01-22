package service

import (
	"github.com/DVKunion/SeaMoon/cmd/client/api/database"
	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
)

type CloudProviderService struct {
}

func (c CloudProviderService) Count(cond ...Condition) int64 {
	var res int64 = 0
	conn := database.GetConn()
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Model(&models.CloudProvider{}).Count(&res)
	return res
}

func (c CloudProviderService) List(page int, size int, preload bool, cond ...Condition) interface{} {
	var data = make([]models.CloudProvider, 0)
	conn := database.QueryPage(page, size)
	if preload {
		conn = conn.Preload("Tunnels.Proxies")
	}
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Find(&data)
	return data
}

func (c CloudProviderService) GetById(id uint) interface{} {
	var data = models.CloudProvider{}
	//database.GetConn().Unscoped().Limit(1).Where("ID = ?", id).Find(&data)
	database.GetConn().Limit(1).Where("ID = ?", id).Find(&data)
	return &data
}

func (c CloudProviderService) Create(obj interface{}) interface{} {
	database.GetConn().Create(obj.(*models.CloudProvider))
	return obj
}

func (c CloudProviderService) Update(id uint, obj interface{}) interface{} {
	obj.(*models.CloudProvider).ID = id
	database.GetConn().Updates(obj.(*models.CloudProvider))
	return c.GetById(id)
}

func (c CloudProviderService) Delete(id uint) {
	database.GetConn().Delete(&models.CloudProvider{}, id)
}
