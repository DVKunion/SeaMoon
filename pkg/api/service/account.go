package service

import (
	"context"
	"errors"

	"github.com/DVKunion/SeaMoon/pkg/api/models"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
)

var paramsMissingError = errors.New("missing important params")

type account struct {
}

func (a *account) Login(ctx context.Context, auth *models.Account) (string, error) {
	//if auth.Username == "" || auth.Password == "" {
	//	return "", paramsMissingError
	//}

	// 检查用户密码是否正确
	//hash := md5.New()

	// 写入数据到哈希实例中
	//hash.Write([]byte(auth.Password))

	// 检查用户是否存在
	//data, err := dao.Q.Auth.WithContext(ctx).Where(
	//	dao.Q.Auth.Username.Eq(auth.Username),
	//	dao.Q.Auth.Password.Eq(strings.ToLower(hex.EncodeToString(hash.Sum(nil)))),
	//	dao.Q.Auth.Type.Eq(int8(enum.AuthAdmin)),
	//).First()

	//if err != nil || data == nil {
	//	return "", err
	//}

	//return tools.JWTAuth(auth.Username)
	return "", nil
}

func (a *account) CreateAuth(ctx context.Context, obj *models.Account) error {
	//return dao.Q.WithContext(ctx).Create(obj)
	return nil
}

func (a *account) UpdatePassword(ctx context.Context, auth *models.Account) error {
	/*if auth.Username == "" || auth.Password == "" {
		return paramsMissingError
	}*/

	// 检查用户密码是否正确
	//hash := md5.New()

	// 写入数据到哈希实例中
	//hash.Write([]byte(auth.Password))
	//
	//res, err := dao.Q.Auth.WithContext(ctx).
	//	Where(dao.Q.Auth.Username.Eq(auth.Username), dao.Q.Auth.Type.Eq(int8(enum.AuthAdmin))).
	//	Update(dao.Q.Auth.Password, strings.ToLower(hex.EncodeToString(hash.Sum(nil))))
	//
	//if err != nil || res.Error != nil {
	//	return err
	//}
	tools.RollJWTSecret()

	return nil
}
