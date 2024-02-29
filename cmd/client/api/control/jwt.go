package control

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/DVKunion/SeaMoon/cmd/client/api/types"
	"github.com/DVKunion/SeaMoon/pkg/tools"
)

var secret = []byte(tools.GenerateRandomString(64))

func JWTAuth(user string, t types.AuthType) (string, error) {
	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"type": t,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(secret)
}

func JWTAuthMiddleware(debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// debug 模式跳过认证
		if debug {
			c.Next()
			return
		}
		// 这里简化了JWT验证逻辑，实际使用时应更复杂
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, errorMsg(10011))
			c.Abort()
			return
		}

		// 验证token是否有效
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil {
			c.JSON(http.StatusForbidden, errorMsg(10010))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if cast, ok := claims["type"].(float64); ok && types.AuthType(cast) == types.Admin {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, errorMsg(10012))
		c.Abort()
		return
	}
}
