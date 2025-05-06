package web

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"

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
	// Initialize DB
	ctx, cancel := context.WithTimeout(context.Background(), CONTEXT_TIMEOUT)
	defer cancel()
	queries := connectToDatabase(cfg)
	createAdminUser(ctx, queries, cfg)

	// Services
	validator := validation.NewValidationService()
	userService := user.NewUserService(queries, validator)
	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJwtService(cfg.App.JwtSecret, ISSUER, TWENTY_FOUR_HOURS_IN_SECONDS)
	loginRegisterService := loginRegister.NewLoginRegisterService(userService, passwordService, jwtService)

	// Handlers
	registerHandler := handlers.NewRegisterHandler(loginRegisterService)
	loginHandler := handlers.NewLoginHandler(loginRegisterService)

	// Middlewares
	authenticationMiddleware := mw.NewAuthenticationMiddleware(cfg.App.JwtSecret)
	authorizationMiddleware := mw.NewAuthorizationMiddleware()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public Routes
	e.Static("/", "public")
	e.GET("/", handlers.HomeHandler)
	e.GET("/login", loginHandler.LoginRegisterContainerHandler)
	e.GET("/loginForm", loginHandler.LoginFormHandler)
	e.POST("/api/login", loginHandler.LoginHandler)
	e.GET("/registerForm", registerHandler.RegisterFormHandler)
	e.POST("/api/register", registerHandler.RegisterUserHandler)

	// JWT Middleware only
	res := e.Group("/restricted")
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
		name := claims.UserId
		role := claims.UserRole

		return c.String(http.StatusOK, "Welcome "+name+" with role: "+role+"!")
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
		name := claims.UserId
		role := claims.UserRole

		return c.String(http.StatusOK, "Welcome "+name+" with role: "+role+"!e")
	})
}

func connectToDatabase(cfg *config.Config) *repository.Queries {
	DATABASE_PATH := "../../data/"
	dbFilePath := DATABASE_PATH + "database.db"

	var dbPersistence string
	switch cfg.Db.Persistence {
	case "FILE":
		if err := os.MkdirAll(DATABASE_PATH, DIR_PERMISSION); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
		dbPersistence = dbFilePath
	default:
		dbPersistence = "file::memory:?cache=shared"
	}

	db, err := sql.Open("sqlite", dbPersistence)
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations using golang-migrate
	migrationPath := "file://" + cfg.Db.MigrationsPath
	m, err := migrate.New(migrationPath, "sqlite://"+dbPersistence)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Database migrations applied successfully.")
	absPath, _ := filepath.Abs(dbPersistence)
	log.Printf("Using SQLite database at: %s", absPath)
	return repository.New(db)
}

func createAdminUser(ctx context.Context, queries *repository.Queries, cfg *config.Config) {
	adminId := uuid.NewString()
	adminName := cfg.App.AdminUser
	adminPassword := cfg.App.AdminPassword
	adminEmail := cfg.App.AdminEmail
	hashedPassword, err := password.NewPasswordService().HashAndSaltPassword(adminPassword)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		return
	}

	exists, err := queries.UserExistsByEmail(ctx, cfg.App.AdminEmail)
	if err != nil {
		log.Printf("Error checking if admin user exists: %v\n", err)
		return
	}

	if exists == 1 {
		log.Println("Admin user already exists, skipping creation.")
		return
	}

	userParams := repository.CreateUserParams{
		ID:           adminId,
		Username:     adminName,
		Email:        adminEmail,
		PasswordHash: hashedPassword,
		UserRole:     "ADMIN",
	}
	user, err := queries.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating admin user: %v\n", err)
		return
	}

	log.Printf("Admin user created successfully:\n"+
		"	id: %q\n	email: %q\n	username: %q\n",
		user.ID,
		user.Username,
		user.Email,
	)
}
