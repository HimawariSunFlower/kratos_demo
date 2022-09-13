package auth

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"time"
)

type MyClaims struct {
	Uid uint `json:"uid"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, key string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		Uid: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)), //设置JWT过期时间,此处设置为2小时
			Issuer:    "test",                                            //设置签发人
		}})
	signedString, err := claims.SignedString([]byte(key))

	return signedString, err
}
