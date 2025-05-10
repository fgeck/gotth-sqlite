package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	"github.com/fgeck/gotth-sqlite/internal/service/security/password"
	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
)

type DbAdminInterface interface {
	ConnectToDatabase() *Queries
	Migrate() error
	CreateAdminUser(email, userName, password string) error
}

type DbAdminService struct {
	dbFilePath      string
	queries         *Queries
	passwordService password.PasswordServiceInterface
}

const (
	DIR_PERMISSION = 0755 // Read, write, execute for owner; read, execute for group and others
)

func NewDbAdminService(passwordService password.PasswordServiceInterface) *DbAdminService {
	return &DbAdminService{
		passwordService: passwordService,
	}
}

func (d *DbAdminService) ConnectToDatabase(databasePath string) (*Queries, error) {
	dbDirPath := "./data/"
	if databasePath != "" {
		dbDirPath = databasePath
	}
	if err := os.MkdirAll(dbDirPath, DIR_PERMISSION); err != nil {
		return nil, err
	}

	d.dbFilePath = filepath.Join(dbDirPath, "database.db")

	db, err := sql.Open("sqlite", d.dbFilePath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	absPath, _ := filepath.Abs(d.dbFilePath)
	log.Printf("Using SQLite database at: %s", absPath)
	d.queries = New(db)
	return d.queries, nil
}

func (d *DbAdminService) Migrate(migrationsPath string) error {
	migrationPath := "file://" + migrationsPath
	m, err := migrate.New(migrationPath, "sqlite://"+d.dbFilePath)
	if err != nil {
		return err
	}

	// Apply migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("Database migrations applied successfully.")
	return nil
}

func (d *DbAdminService) CreateAdminUser(ctx context.Context, email, username, password string) error {
	adminId := uuid.NewString()

	hashedPassword, err := d.passwordService.HashAndSaltPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		return err
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	exists, err := d.queries.UserExistsByEmail(dbCtx, email)
	if err != nil {
		log.Printf("Error checking if admin user exists: %v\n", err)
		return err
	}

	if exists == 1 {
		log.Println("Admin user already exists, skipping creation.")
		return nil
	}

	userParams := CreateUserParams{
		ID:           adminId,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		UserRole:     "ADMIN",
	}
	user, err := d.queries.CreateUser(dbCtx, userParams)
	if err != nil {
		log.Printf("Error creating admin user: %v\n", err)
		return err
	}

	log.Printf("Admin user created successfully:\n"+
		"	id: %q\n	email: %q\n	username: %q\n",
		user.ID,
		user.Username,
		user.Email,
	)
	return nil
}
