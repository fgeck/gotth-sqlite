package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fgeck/go-register/internal/repository"
	"github.com/fgeck/go-register/internal/service/validation"
	"github.com/google/uuid"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, username, email, passwordHash string) (*UserCreatedDto, error)
	GetUserByEmail(ctx context.Context, email string) (*UserDto, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	ValidateCreateUserParams(username, email, password string) error
}

type UserService struct {
	queries   repository.Querier
	validator validation.ValidationServiceInterface
}

func NewUserService(queries repository.Querier, validator validation.ValidationServiceInterface) *UserService {
	return &UserService{
		queries:   queries,
		validator: validator,
	}
}

func NewUserCreatedDto(username, email string) *UserCreatedDto {
	return &UserCreatedDto{username, email}
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*UserDto, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return NewUserDto(user)
}

func (s *UserService) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := s.queries.UserExistsByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	return exists == 1, err
}

func (s *UserService) CreateUser(ctx context.Context, username, email, hashedPassword string) (*UserCreatedDto, error) {
	user, err := s.queries.CreateUser(
		ctx,
		repository.CreateUserParams{
			ID:           uuid.NewString(),
			Username:     username,
			Email:        email,
			PasswordHash: hashedPassword,
			UserRole:     UserRoleUser.Name,
		},
	)
	if err != nil {
		// Todo log error
		return nil, err
	}

	return NewUserCreatedDto(user.Username, user.Email), nil
}

func (s *UserService) ValidateCreateUserParams(username, email, password string) error {
	if err := s.validator.ValidateEmail(email); err != nil {
		return err
	}

	if err := s.validator.ValidatePassword(password); err != nil {
		return err
	}

	if err := s.validator.ValidateUsername(username); err != nil {
		return err
	}

	return nil
}
