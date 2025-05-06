package jwt

import (
	goJwt "github.com/golang-jwt/jwt/v5"
)

const (
	USER_ID   = "userId"
	USER_ROLE = "userRole"
)

type JwtCustomClaims struct {
	UserId   string `json:"userId"`
	UserRole string `json:"userRole"`
	goJwt.RegisteredClaims
}

func NewJwtCustomClaims(userId, userRole string, standardClaims goJwt.RegisteredClaims) *JwtCustomClaims {
	return &JwtCustomClaims{
		UserId:           userId,
		UserRole:         userRole,
		RegisteredClaims: standardClaims,
	}
}
