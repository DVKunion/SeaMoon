package service

import (
	"github.com/DVKunion/SeaMoon/cmd/client/api/database"
	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
)

type AuthService struct {
}

func (a AuthService) Count(cond ...Condition) int64 {
	var res int64 = 0
	conn := database.GetConn()
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Model(&models.Auth{}).Count(&res)
	return res
}

func (a AuthService) List(page int, size int, preload bool, cond ...Condition) interface{} {
	var data = make([]models.Auth, 0)
	conn := database.QueryPage(page, size)
	for _, cs := range cond {
		conn = conn.Where(cs.Key+" = ?", cs.Value)
	}
	conn.Find(&data)
	return data
}

func (a AuthService) GetById(id uint) interface{} {
	var data = models.Auth{}
	database.GetConn().Limit(1).Where("ID = ?", id).Find(&data)
	return &data
}

func (a AuthService) Create(obj interface{}) interface{} {
	database.GetConn().Create(obj.(*models.Auth))
	return obj
}

func (a AuthService) Update(id uint, obj interface{}) interface{} {
	obj.(*models.Auth).ID = id
	database.GetConn().Updates(obj.(*models.Auth))
	return a.GetById(id)
}

func (a AuthService) Delete(id uint) {
	database.GetConn().Delete(&models.Auth{}, id)
}
