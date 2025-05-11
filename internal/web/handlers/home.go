package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fgeck/gotth-sqlite/internal/service/render"
	components "github.com/fgeck/gotth-sqlite/templates/components"
	views "github.com/fgeck/gotth-sqlite/templates/views"
	echo "github.com/labstack/echo/v4"
)

func HomeHandler(ctx echo.Context) error {
	isCollapsed := isSidebarCollapsed(ctx.Request())
	if err := render.Render(ctx, views.Home(isCollapsed)); err != nil {
		return fmt.Errorf("failed to render home view: %w", err)
	}

	return nil
}

func SideBarHandler(ctx echo.Context) error {
	currentState := isSidebarCollapsed(ctx.Request())
	newState := !currentState
	setSidebarCollapsed(ctx.Response().Writer, newState)

	render.Render(ctx, components.Sidebar(newState))
	render.Render(ctx, components.MainContent(newState))
	return nil
}

const sidebarCookieName = "sidebar_collapsed"

func isSidebarCollapsed(r *http.Request) bool {
	cookie, err := r.Cookie(sidebarCookieName)
	return err == nil && cookie.Value == "true"
}

func setSidebarCollapsed(w http.ResponseWriter, collapsed bool) {
	http.SetCookie(w, &http.Cookie{
		Name:  sidebarCookieName,
		Value: strconv.FormatBool(collapsed),
		Path:  "/",
	})
}
