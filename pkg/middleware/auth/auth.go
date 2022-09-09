package auth

import jwt "github.com/golang-jwt/jwt/v4"

type MyClaims struct {
	Uid uint64 `json:"uid"`
	jwt.RegisteredClaims
}
