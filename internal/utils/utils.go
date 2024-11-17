package utils

import (
	"ai-qa-service/pkg/errno"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type UserClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

var aExpireDuration time.Duration
var jwtKey []byte
var issuer string

func InitToken() {
	// 读取配置
	aExpireDuration = time.Duration(viper.GetInt64("service.token.access_token_expire_duration")) * time.Second  // access token 过期时间
	jwtKey = []byte("114514")
	issuer = "Sky_Lee"
}

func GenToken(UserID int64) (string, error) {
	cliams := &UserClaims{
		UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(aExpireDuration)),
			Issuer:    issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (UserID int64, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, errno.ErrExpiredToken
		}
		return 0, err
	}

	cliams, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return 0, errno.ErrInvalidToken
	}

	return cliams.UserID, nil
}

func GetAccessTokenExpireDuration() time.Duration {
	return aExpireDuration
}

func GetJwtKey() []byte {
	return jwtKey
}
