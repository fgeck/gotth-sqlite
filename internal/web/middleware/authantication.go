package middleware

import (
	"github.com/fgeck/gotth-sqlite/internal/service/security/jwt"
	gojwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthenticationMiddlewareInterface interface {
	JwtAuthMiddleware(jwtSecret string) echo.MiddlewareFunc
}
type AuthenticationMiddleware struct {
	jwtSecret string
}

func NewAuthenticationMiddleware(jwtSecret string) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		jwtSecret: jwtSecret,
	}
}

func (a *AuthenticationMiddleware) JwtAuthMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(a.jwtSecret),
		TokenLookup: "cookie:token",
		NewClaimsFunc: func(c echo.Context) gojwt.Claims {
			return new(jwt.JwtCustomClaims)
		},
	})
}
