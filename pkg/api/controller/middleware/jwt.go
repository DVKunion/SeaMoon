package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/DVKunion/SeaMoon/pkg/api/controller/servant"
	"github.com/DVKunion/SeaMoon/pkg/api/enum"
	"github.com/DVKunion/SeaMoon/pkg/system/errors"
	"github.com/DVKunion/SeaMoon/pkg/system/tools"
	"github.com/DVKunion/SeaMoon/pkg/system/xlog"
)

func JWTAuthMiddleware(c *gin.Context) {
	// 这里简化了JWT验证逻辑，实际使用时应更复杂
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		tokenString = c.Query("token")
	}

	if tokenString == "" {
		servant.ErrorMsg(c, http.StatusUnauthorized, errors.ApiError(xlog.ApiParamsRequire, nil))
		c.Abort()
		return
	}

	// 验证token是否有效
	token, err := tools.JWTParse(tokenString)

	if err != nil {
		// 认证失败
		servant.ErrorMsg(c, http.StatusForbidden, errors.ApiError(xlog.ApiAuthError, err))
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 校验通过
		if cast, ok := claims["type"].(float64); ok && enum.AuthType(cast) == enum.AdminAuth {
			c.Next()
			return
		}
	}

	// 校验不通过
	servant.ErrorMsg(c, http.StatusForbidden, errors.ApiError(xlog.ApiAuthLimit, err))
	c.Abort()
	return
}
