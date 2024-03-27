package service

import (
	"context"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
)

type config struct {
}

func (c *config) ListConfigs(ctx context.Context, page, size int) (models.ConfigList, error) {
	query := dao.Q.Config
	return query.WithContext(ctx).Offset(page * size).Limit(size).Find()
}

func (c *config) GetConfigByName(ctx context.Context, name string) (*models.Config, error) {
	query := dao.Q.Config
	return query.WithContext(ctx).Where(query.Key.Eq(name)).First()
}

func (c *config) CreateConfig(ctx context.Context, obj *models.Config) error {
	return dao.Q.Config.WithContext(ctx).Create(obj)
}

func (c *config) UpdateConfig(ctx context.Context, configList models.ConfigList) error {
	query := dao.Q.Config
	for _, sys := range configList {
		_, err := query.WithContext(ctx).Where(query.Key.Eq(sys.Key)).Update(query.Value, sys.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
