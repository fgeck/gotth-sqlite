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

type HomeHandlerInterface interface {
	HomeHandler(ctx echo.Context) error
	SideBarHandler(ctx echo.Context) error
}

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

const (
	SIDEBAR_COLLAPSED_COOKIE_KEY = "sidebar_collapsed"
)

func (h *HomeHandler) HomeViewHandler(ctx echo.Context) error {
	sidebarCollapsed := h.isSidebarCollapsed(ctx.Request())
	if err := render.Render(ctx, views.Home(sidebarCollapsed, h.isUserLoggedIn(ctx))); err != nil {
		return fmt.Errorf("failed to render home view: %w", err)
	}

	return nil
}

func (h *HomeHandler) SideBarHandler(ctx echo.Context) error {
	currentState := h.isSidebarCollapsed(ctx.Request())
	newState := !currentState
	h.setSidebarCollapsed(ctx.Response().Writer, newState)

	render.Render(ctx, components.Sidebar(newState, h.isUserLoggedIn(ctx)))
	render.Render(ctx, components.MainContent(newState))
	return nil
}

func (h *HomeHandler) isSidebarCollapsed(r *http.Request) bool {
	cookie, err := r.Cookie(SIDEBAR_COLLAPSED_COOKIE_KEY)
	return err == nil && cookie.Value == "true"
}

func (h *HomeHandler) setSidebarCollapsed(w http.ResponseWriter, collapsed bool) {
	http.SetCookie(w, &http.Cookie{
		Name:  SIDEBAR_COLLAPSED_COOKIE_KEY,
		Value: strconv.FormatBool(collapsed),
		Path:  "/",
	})
}

func (h *HomeHandler) isUserLoggedIn(ctx echo.Context) bool {
	user := ctx.Get("user")
	return user != ""
}
