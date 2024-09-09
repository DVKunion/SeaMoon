package tools

import (
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/DVKunion/SeaMoon/pkg/api/enum"
)

var secret = []byte(GenerateRandomString(64))

func JWTAuth(user string) (string, error) {
	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"type": enum.AdminAuth,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(secret)
}

func JWTParse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func RollJWTSecret() {
	secret = []byte(GenerateRandomString(64))
}
