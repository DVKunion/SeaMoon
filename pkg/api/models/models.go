package models

import (
	"github.com/DVKunion/SeaMoon/pkg/api/models/model"
)

func init() {
	ModelList = append(ModelList, &model.Account{})
	ModelList = append(ModelList, &model.Proxy{})
	ModelList = append(ModelList, &model.Tunnel{})
	ModelList = append(ModelList, &model.Provider{})
	ModelList = append(ModelList, &model.Config{})
}

var ModelList = make([]interface{}, 0)
