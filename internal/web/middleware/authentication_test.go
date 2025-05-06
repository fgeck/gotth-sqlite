package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fgeck/gotth-sqlite/internal/service/security/jwt"
	mw "github.com/fgeck/gotth-sqlite/internal/web/middleware"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwtAuthMiddleware(t *testing.T) {
	t.Parallel()
	jwtSecret := "testsecret"
	middleware := mw.NewAuthenticationMiddleware(jwtSecret).JwtAuthMiddleware()

	t.Run("No token provided", func(t *testing.T) {
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
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("Invalid token provided", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "token", Value: "invalidtoken"})
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})

		err := handler(c)
		require.Error(t, err)
		var httpErr *echo.HTTPError
		require.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	})

	t.Run("Valid token provided", func(t *testing.T) {
		t.Parallel()
		e := echo.New()
		token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &jwt.JwtCustomClaims{
			RegisteredClaims: gojwt.RegisteredClaims{
				Issuer: "test",
			},
		})
		tokenString, _ := token.SignedString([]byte(jwtSecret))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "success")
		})

		err := handler(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	})
}
