package handlers

import (
	"errors"
	"fmt"
	"net/http"

	customErrors "github.com/fgeck/go-register/internal/service/errors"
	"github.com/fgeck/go-register/internal/service/loginRegister"
	"github.com/fgeck/go-register/internal/service/render"
	"github.com/fgeck/go-register/templates/views"
	echo "github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	loginRegisterService loginRegister.LoginRegisterServiceInterface
}

func NewRegisterHandler(loginRegisterService loginRegister.LoginRegisterServiceInterface) *RegisterHandler {
	return &RegisterHandler{
		loginRegisterService: loginRegisterService,
	}
}

func (r *RegisterHandler) RegisterFormHandler(ctx echo.Context) error {
	if err := render.Render(ctx, views.RegisterForm()); err != nil {
		return fmt.Errorf("failed to render register form: %w", err)
	}

	return nil
}

func (r *RegisterHandler) RegisterUserHandler(ctx echo.Context) error {
	username := ctx.FormValue("username")
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	user, err := r.loginRegisterService.RegisterUser(ctx.Request().Context(), username, email, password)
	if err != nil {
		var userfacingErr *customErrors.UserFacingError
		if errors.As(err, &userfacingErr) {
			jsonErr := ctx.JSON(http.StatusBadRequest, map[string]string{"error": userfacingErr.Error()})
			if jsonErr != nil {
				return fmt.Errorf("failed to send error response: %w", jsonErr)
			}

			return err
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}

	return ctx.JSON(http.StatusCreated, user)
}
