package render

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, cmp templ.Component) error {
	if err := cmp.Render(ctx.Request().Context(), ctx.Response()); err != nil {
		return fmt.Errorf("failed to render component: %w", err)
	}

	return nil
}
