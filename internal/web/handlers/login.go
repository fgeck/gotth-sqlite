package handlers

import (
	"fmt"
	"net/http"

	loginregister "github.com/fgeck/go-register/internal/service/loginRegister"
	"github.com/fgeck/go-register/internal/service/render"

	"github.com/fgeck/go-register/templates/views"
	"github.com/labstack/echo/v4"
)

type LoginHandlerInterface interface {
	LoginRegisterContainerHandler(ctx echo.Context) error
	LoginFormHandler(ctx echo.Context) error
	LoginHandler(ctx echo.Context) error
}

type LoginHandler struct {
	loginRegisterService loginregister.LoginRegisterServiceInterface
}

func NewLoginHandler(loginRegisterService loginregister.LoginRegisterServiceInterface) *LoginHandler {
	return &LoginHandler{
		loginRegisterService: loginRegisterService,
	}
}

func (h *LoginHandler) LoginRegisterContainerHandler(ctx echo.Context) error {
	if err := render.Render(ctx, views.LoginRegister()); err != nil {
		return fmt.Errorf("failed to render login register container: %w", err)
	}

	return nil
}

func (h *LoginHandler) LoginFormHandler(ctx echo.Context) error {
	if err := render.Render(ctx, views.LoginForm()); err != nil {
		return fmt.Errorf("failed to render login form: %w", err)
	}

	return nil
}

func (h *LoginHandler) LoginHandler(ctx echo.Context) error {
	username := ctx.FormValue("email")
	password := ctx.FormValue("password")

	token, err := h.loginRegisterService.LoginUser(ctx.Request().Context(), username, password)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to login user: %w", err)
		jsonErr := ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to login user"})

		if jsonErr != nil {
			return fmt.Errorf("failed to send error response: %w", jsonErr)
		}

		return wrappedErr
	}
	ctx.SetCookie(
		&http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",                  // Cookie is valid for the entire site
			HttpOnly: true,                 // Prevent access via JavaScript
			Secure:   true,                 // Only send the cookie over HTTPS
			SameSite: http.SameSiteLaxMode, // Prevent CSRF attacks
		},
	)

	if err := ctx.String(http.StatusOK, "success"); err != nil {
		return fmt.Errorf("failed to send success response: %w", err)
	}

	return nil
}
