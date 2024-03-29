package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/DVKunion/SeaMoon/pkg/api/database/dao"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

var paramsMissingError = errors.New("missing important params")

type auth struct {
}

func (a *auth) Login(ctx context.Context, auth *models.AuthApi) (string, error) {
	if auth.Username == "" || auth.Password == "" {
		return "", paramsMissingError
	}

	// 检查用户密码是否正确
	hash := md5.New()

	// 写入数据到哈希实例中
	hash.Write([]byte(auth.Password))

	// 检查用户是否存在
	data, err := dao.Q.Auth.WithContext(ctx).Where(
		dao.Q.Auth.Username.Eq(auth.Username),
		dao.Q.Auth.Password.Eq(strings.ToLower(hex.EncodeToString(hash.Sum(nil)))),
		dao.Q.Auth.Type.Eq(int8(enum.AuthAdmin)),
	).First()

	if err != nil || data == nil {
		return "", err
	}

	return tools.JWTAuth(auth.Username)
}

func (a *auth) CreateAuth(ctx context.Context, obj *models.Auth) error {
	return dao.Q.Auth.WithContext(ctx).Create(obj)
}

func (a *auth) UpdatePassword(ctx context.Context, auth *models.AuthApi) error {
	if auth.Username == "" || auth.Password == "" {
		return paramsMissingError
	}

	// 检查用户密码是否正确
	hash := md5.New()

	// 写入数据到哈希实例中
	hash.Write([]byte(auth.Password))

	res, err := dao.Q.Auth.WithContext(ctx).
		Where(dao.Q.Auth.Username.Eq(auth.Username), dao.Q.Auth.Type.Eq(int8(enum.AuthAdmin))).
		Update(dao.Q.Auth.Password, strings.ToLower(hex.EncodeToString(hash.Sum(nil))))

	if err != nil || res.Error != nil {
		return err
	}
	tools.RollJWTSecret()

	return nil
}
