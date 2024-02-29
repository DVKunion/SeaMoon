package service

import (
	"github.com/DVKunion/SeaMoon/cmd/client/api/database"
	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
)

type SystemConfigService struct {
}

func (s SystemConfigService) Count(cond ...Condition) int64 {
	var res int64 = 0
	conn := database.GetConn()
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Model(&models.SystemConfig{}).Count(&res)
	return res
}

func (s SystemConfigService) List(page, size int, preload bool, cond ...Condition) interface{} {
	var data = make([]*models.SystemConfig, 0)
	conn := database.QueryPage(page, size)
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Find(&data)
	return data
}

func (s SystemConfigService) GetById(id uint) interface{} {
	var data = models.SystemConfig{}
	database.GetConn().Limit(1).Where("ID = ?", id).Find(&data)
	return &data
}

func (s SystemConfigService) GetByName(name string) *models.SystemConfig {
	var data = models.SystemConfig{}
	database.GetConn().Limit(1).Where("KEY = ?", name).Find(&data)
	return &data
}

func (s SystemConfigService) Create(obj interface{}) interface{} {
	database.GetConn().Create(obj.(*models.SystemConfig))
	return obj
}

func (s SystemConfigService) Update(id uint, obj interface{}) interface{} {
	for _, sys := range obj.([]models.SystemConfig) {
		if sys.Key == "version" || sys.Value == "" {
			continue
		}
		// api 传过来的不知道 id 的，从 key 去查吧
		target := s.GetByName(sys.Key)
		target.Value = sys.Value
		database.GetConn().Updates(target)
	}

	return nil
}

func (s SystemConfigService) Delete(id uint) {
	// 系统配置禁止删除
	return
}
