package middleware

import (
	"strings"

	"github.com/fgeck/gotth-sqlite/internal/service/security/jwt"
	user "github.com/fgeck/gotth-sqlite/internal/service/user"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthorizationMiddlewareInterface interface {
	RequireAdminMiddleware() echo.MiddlewareFunc
}

type AuthorizationMiddleware struct{}

func NewAuthorizationMiddleware() *AuthorizationMiddleware {
	return &AuthorizationMiddleware{}
}

func (a *AuthorizationMiddleware) RequireAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("user").(*gojwt.Token)
			if !ok {
				return echo.ErrForbidden
			}
			claims, ok := token.Claims.(*jwt.JwtCustomClaims)
			if !ok {
				return echo.ErrForbidden
			}
			if claims.UserRole == "" || strings.ToUpper(claims.UserRole) != user.UserRoleAdmin.Name {
				return echo.ErrForbidden
			}
			return next(c)
		}
	}
}
