package web

import (
	"context"
	"net/http"
	"time"

	"github.com/fgeck/gotth-sqlite/internal/repository"
	"github.com/fgeck/gotth-sqlite/internal/service/config"
	"github.com/fgeck/gotth-sqlite/internal/service/loginRegister"
	"github.com/fgeck/gotth-sqlite/internal/service/security/jwt"
	"github.com/fgeck/gotth-sqlite/internal/service/security/password"
	"github.com/fgeck/gotth-sqlite/internal/service/user"
	"github.com/fgeck/gotth-sqlite/internal/service/validation"
	"github.com/fgeck/gotth-sqlite/internal/web/handlers"
	mw "github.com/fgeck/gotth-sqlite/internal/web/middleware"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	TWENTY_FOUR_HOURS_IN_SECONDS = 24 * 60 * 60
	ISSUER                       = "gotth-sqlite"
	CONTEXT_TIMEOUT              = 10 * time.Second
	DIR_PERMISSION               = 0755 // Read, write, execute for owner; read, execute for group and others
)

func InitServer(e *echo.Echo, cfg *config.Config) {
	ctx := context.Background()

	// Initialize DB
	passwordService := password.NewPasswordService()
	dbAdmin := repository.NewDbAdminService(passwordService)
	queries, err := dbAdmin.ConnectToDatabase(cfg.Db.DataBasePath)
	if err != nil {
		e.Logger.Fatalf("Failed to connect to database", err)
	}
	err = dbAdmin.Migrate(cfg.Db.MigrationsPath)
	if err != nil {
		e.Logger.Fatalf("Failed to migrate database", err)
	}
	err = dbAdmin.CreateAdminUser(ctx, cfg.App.AdminEmail, cfg.App.AdminUser, cfg.App.AdminPassword)
	if err != nil {
		e.Logger.Fatalf("Failed to create Admin User", err)
	}

	// Services
	validator := validation.NewValidationService()
	userService := user.NewUserService(queries, validator)
	jwtService := jwt.NewJwtService(cfg.App.JwtSecret, ISSUER, TWENTY_FOUR_HOURS_IN_SECONDS)
	loginRegisterService := loginRegister.NewLoginRegisterService(userService, passwordService, jwtService)

	// Handlers
	homeHandler  := handlers.NewHomeHandler()
	registerHandler := handlers.NewRegisterHandler(loginRegisterService)
	loginHandler := handlers.NewLoginHandler(loginRegisterService)

	// Middlewares
	authenticationMiddleware := mw.NewAuthenticationMiddleware(cfg.App.JwtSecret)
	authorizationMiddleware := mw.NewAuthorizationMiddleware()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public Routes
	e.Static("/", "public")
	e.GET("/", homeHandler.HomeViewHandler)
	e.GET("/toggle-sidebar", homeHandler.SideBarHandler)
	e.GET("/login", loginHandler.LoginRegisterContainerHandler)
	e.GET("/loginForm", loginHandler.LoginFormHandler)
	e.POST("/api/login", loginHandler.LoginHandler)
	e.GET("/registerForm", registerHandler.RegisterFormHandler)
	e.POST("/api/register", registerHandler.RegisterUserHandler)

	// JWT Middleware only
	res := e.Group("/api/restricted")
	res.Use(authenticationMiddleware.JwtAuthMiddleware())
	// for testing purposes
	res.GET("", func(c echo.Context) error {
		token, ok := c.Get("user").(*gojwt.Token)
		if !ok {
			return echo.ErrForbidden
		}
		claims, ok := token.Claims.(*jwt.JwtCustomClaims)
		if !ok {
			return echo.ErrForbidden
		}
		id := claims.UserId
		role := claims.UserRole

		return c.String(http.StatusOK, "Welcome user with ID "+id+" and with role: "+role+"!")
	})

	// Admin Routes (requires "UserRole" == "admin")
	// Admin Routes (requires "UserRole" == "admin")
	adminGroup := e.Group("/api/admin")
	adminGroup.Use(authenticationMiddleware.JwtAuthMiddleware(), authorizationMiddleware.RequireAdminMiddleware())
	// for testing purposes
	adminGroup.GET("/users", func(c echo.Context) error {
		token, ok := c.Get("user").(*gojwt.Token)
		if !ok {
			return echo.ErrForbidden
		}
		claims, ok := token.Claims.(*jwt.JwtCustomClaims)
		if !ok {
			return echo.ErrForbidden
		}
		id := claims.UserId
		role := claims.UserRole

		return c.String(http.StatusOK, "Welcome user with ID "+id+" and with role: "+role+"!")
	})
}
