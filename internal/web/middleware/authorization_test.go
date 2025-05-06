package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fgeck/go-register/internal/service/security/jwt"
	"github.com/fgeck/go-register/internal/service/user"
	mw "github.com/fgeck/go-register/internal/web/middleware"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAdminMiddleware(t *testing.T) {
	t.Parallel()
	middleware := mw.NewAuthorizationMiddleware().RequireAdminMiddleware()

	t.Run("No token in context", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})

		err := handler(c)
		require.Error(t, err)
		var httpErr *echo.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("Invalid token type in context", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", "invalid_token_type")

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})

		err := handler(c)
		require.Error(t, err)
		var httpErr *echo.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("Invalid claims type in token", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		token := &gojwt.Token{Claims: gojwt.MapClaims{}}
		c.Set("user", token)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})

		err := handler(c)
		require.Error(t, err)
		var httpErr *echo.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("User role is not ADMIN", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		claims := &jwt.JwtCustomClaims{UserRole: "USER"}
		token := &gojwt.Token{Claims: claims}
		c.Set("user", token)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})

		err := handler(c)
		require.Error(t, err)
		var httpErr *echo.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusForbidden, httpErr.Code)
	})

	t.Run("User role is ADMIN", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		claims := &jwt.JwtCustomClaims{UserRole: strings.ToUpper(user.UserRoleAdmin.Name)}
		token := &gojwt.Token{Claims: claims}
		c.Set("user", token)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})

		err := handler(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})
}
