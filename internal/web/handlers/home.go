package handlers

import (
	"fmt"

	"github.com/fgeck/gotth-sqlite/internal/service/render"
	views "github.com/fgeck/gotth-sqlite/templates/views"
	echo "github.com/labstack/echo/v4"
)

func HomeHandler(ctx echo.Context) error {
	if err := render.Render(ctx, views.Home()); err != nil {
		return fmt.Errorf("failed to render home view: %w", err)
	}

	return nil
}
