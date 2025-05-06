package loginRegister

import (
	"context"
	"fmt"

	customErrors "github.com/fgeck/gotth-sqlite/internal/service/errors"
	"github.com/fgeck/gotth-sqlite/internal/service/security/jwt"
	"github.com/fgeck/gotth-sqlite/internal/service/security/password"
	"github.com/fgeck/gotth-sqlite/internal/service/user"
)

type LoginRegisterServiceInterface interface {
	LoginUser(ctx context.Context, email, password string) (string, error)
	RegisterUser(ctx context.Context, username, email, password string) (*user.UserCreatedDto, error)
}

type LoginRegisterService struct {
	userService     user.UserServiceInterface
	passwordService password.PasswordServiceInterface
	jwtService      jwt.JwtServiceInterface
}

func NewLoginRegisterService(
	userService user.UserServiceInterface,
	passwordService password.PasswordServiceInterface,
	jwtService jwt.JwtServiceInterface,
) *LoginRegisterService {
	return &LoginRegisterService{userService: userService, passwordService: passwordService, jwtService: jwtService}
}

func (s *LoginRegisterService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := s.passwordService.ComparePassword(user.PasswordHash, password); err != nil {
		return "", customErrors.NewInternal("invalid password")
	}

	return s.jwtService.GenerateToken(user)
}

func (s *LoginRegisterService) RegisterUser(
	ctx context.Context,
	username string,
	email string,
	password string,
) (*user.UserCreatedDto, error) {
	userExists, err := s.userService.UserExistsByEmail(ctx, email)
	if err != nil {
		// Todo log error
		return nil, err
	}

	if userExists {
		return nil, customErrors.NewUserFacing("user already exists")
	}

	if err := s.userService.ValidateCreateUserParams(username, email, password); err != nil {
		return nil, customErrors.NewUserFacing("failed to validate create user parameters: " + err.Error())
	}

	hashedPassword, err := s.passwordService.HashAndSaltPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to salt and hash password: %w", err)
	}

	userCreatedDto, err := s.userService.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return userCreatedDto, nil
}
