package control

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/DVKunion/SeaMoon/cmd/client/api/models"
	"github.com/DVKunion/SeaMoon/cmd/client/api/service"
	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/pkg/tools"
)

type AuthController struct {
	svc service.ApiService
}

func (a AuthController) Login(c *gin.Context) {
	var obj *models.Auth
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	if obj.Username == "" || obj.Password == "" {
		c.JSON(http.StatusBadRequest, errorMsg(10003))
		return
	}
	// 检查用户是否存在
	data := a.svc.List(0, 1, false, service.Condition{Key: "USERNAME", Value: obj.Username}, service.Condition{Key: "TYPE", Value: types.Admin})
	if len(data.([]models.Auth)) != 1 {
		c.JSON(http.StatusForbidden, errorMsg(10010))
		return
	}
	// 检查用户密码是否正确
	target := data.([]models.Auth)[0]
	hash := md5.New()

	// 写入数据到哈希实例中
	hash.Write([]byte(obj.Password))

	// 计算哈希值
	if target.Password != strings.ToLower(hex.EncodeToString(hash.Sum(nil))) {
		c.JSON(http.StatusForbidden, errorMsg(10010))
		return
	}

	token, err := JWTAuth(target.Username, target.Type)
	if err != nil {
		c.JSON(http.StatusForbidden, errorMsg(10010))
		return
	}

	c.JSON(http.StatusOK, successMsg(1, token))
}

func (a AuthController) Passwd(c *gin.Context) {
	var obj *models.Auth
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, errorMsg(10000))
		return
	}
	if obj.Username == "" || obj.Password == "" {
		c.JSON(http.StatusBadRequest, errorMsg(10003))
		return
	}
	// 检查用户是否存在
	data := a.svc.List(0, 1, false, service.Condition{Key: "USERNAME", Value: obj.Username}, service.Condition{Key: "TYPE", Value: types.Admin})
	if len(data.([]models.Auth)) != 1 {
		c.JSON(http.StatusBadRequest, errorMsg(10002))
		return
	}
	target := data.([]models.Auth)[0]
	hash := md5.New()

	// 写入数据到哈希实例中
	hash.Write([]byte(obj.Password))

	// 更新数据
	target.Password = strings.ToLower(hex.EncodeToString(hash.Sum(nil)))
	a.svc.Update(target.ID, &target)
	// 更新 jwt secret 让之前的 token 失效
	secret = []byte(tools.GenerateRandomString(64))
	c.JSON(http.StatusOK, successMsg(1, "更新成功"))
}
