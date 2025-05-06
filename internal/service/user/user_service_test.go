//go:build unittest

package user_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/fgeck/go-register/internal/repository"
	repositoryMocks "github.com/fgeck/go-register/internal/repository/mocks"
	userfacing_errors "github.com/fgeck/go-register/internal/service/errors"
	"github.com/fgeck/go-register/internal/service/user"
	validationMocks "github.com/fgeck/go-register/internal/service/validation/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupUserServiceTest(t *testing.T) (*repositoryMocks.MockQuerier, *validationMocks.MockValidationServiceInterface, *user.UserService) {
	mockQueries := repositoryMocks.NewMockQuerier(t)
	mockValidator := validationMocks.NewMockValidationServiceInterface(t)
	userService := user.NewUserService(mockQueries, mockValidator)
	return mockQueries, mockValidator, userService
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	username := "testuser"
	email := "testuser@example.com"
	passwordHash := "hashedpassword"

	t.Run("successfully creates a user", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)

		mockQueries.On("CreateUser", ctx, mock.Anything).Return(repository.User{
			Username: username,
			Email:    email,
		}, nil)

		userDto, err := userService.CreateUser(ctx, username, email, passwordHash)

		require.NoError(t, err)
		assert.NotNil(t, userDto)
		assert.Equal(t, username, userDto.Username)
		assert.Equal(t, email, userDto.Email)

		mockQueries.AssertExpectations(t)
	})

	t.Run("fails when database error occurs", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("CreateUser", ctx, mock.Anything).Return(repository.User{}, errors.New("database error"))

		userDto, err := userService.CreateUser(ctx, username, email, passwordHash)

		require.Error(t, err)
		assert.Nil(t, userDto)
		assert.Equal(t, "database error", err.Error())

		mockQueries.AssertExpectations(t)
	})
}

func TestValidateCreateUserParams(t *testing.T) {
	username := "testuser"
	email := "testuser@example.com"
	password := "Valid1@"

	t.Run("successfully validates parameters", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(nil)
		mockValidator.On("ValidatePassword", password).Return(nil)
		mockValidator.On("ValidateUsername", username).Return(nil)

		err := userService.ValidateCreateUserParams(username, email, password)

		require.NoError(t, err)

		mockValidator.AssertExpectations(t)
	})

	t.Run("fails when email validation fails", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(userfacing_errors.NewUserFacing("invalid email format"))

		err := userService.ValidateCreateUserParams(username, email, password)

		require.Error(t, err)

		// Check for UserFacingError
		ufe, ok := err.(*userfacing_errors.UserFacingError)
		assert.True(t, ok, "expected a UserFacingError")
		assert.Equal(t, "invalid email format", ufe.Message)

		mockValidator.AssertExpectations(t)
	})

	t.Run("fails when password validation fails", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(nil)
		mockValidator.On("ValidatePassword", password).Return(userfacing_errors.NewUserFacing("password too weak"))

		err := userService.ValidateCreateUserParams(username, email, password)

		require.Error(t, err)

		// Check for UserFacingError
		ufe, ok := err.(*userfacing_errors.UserFacingError)
		assert.True(t, ok, "expected a UserFacingError")
		assert.Equal(t, "password too weak", ufe.Message)

		mockValidator.AssertExpectations(t)
	})

	t.Run("fails when username validation fails", func(t *testing.T) {
		_, mockValidator, userService := setupUserServiceTest(t)
		mockValidator.On("ValidateEmail", email).Return(nil)
		mockValidator.On("ValidatePassword", password).Return(nil)
		mockValidator.On("ValidateUsername", username).Return(userfacing_errors.NewUserFacing("username too short"))

		err := userService.ValidateCreateUserParams(username, email, password)

		require.Error(t, err)

		// Check for UserFacingError
		ufe, ok := err.(*userfacing_errors.UserFacingError)
		assert.True(t, ok, "expected a UserFacingError")
		assert.Equal(t, "username too short", ufe.Message)

		mockValidator.AssertExpectations(t)
	})
}

func TestUserExistsByEmail(t *testing.T) {
	ctx := context.Background()
	email := "testuser@example.com"

	t.Run("returns true when user exists", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("UserExistsByEmail", ctx, email).Return(int64(1), nil)

		exists, err := userService.UserExistsByEmail(ctx, email)

		require.NoError(t, err)
		assert.True(t, exists)

		mockQueries.AssertExpectations(t)
	})

	t.Run("returns false when user does not exist", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("UserExistsByEmail", ctx, email).Return(int64(0), nil)

		exists, err := userService.UserExistsByEmail(ctx, email)

		require.NoError(t, err)
		assert.False(t, exists)

		mockQueries.AssertExpectations(t)
	})

	t.Run("returns error when query fails", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("UserExistsByEmail", ctx, email).Return(int64(0), errors.New("database error"))

		exists, err := userService.UserExistsByEmail(ctx, email)

		require.Error(t, err)
		assert.False(t, exists)
		assert.Equal(t, "database error", err.Error())

		mockQueries.AssertExpectations(t)
	})
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()
	email := "testuser@example.com"

	t.Run("successfully retrieves user", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("GetUserByEmail", ctx, email).Return(repository.User{
			ID:       uuid.NewString(),
			Username: "testuser",
			Email:    email,
		}, nil)

		userDto, err := userService.GetUserByEmail(ctx, email)

		require.NoError(t, err)
		assert.NotNil(t, userDto)
		assert.Equal(t, "testuser", userDto.Username)
		assert.Equal(t, email, userDto.Email)

		mockQueries.AssertExpectations(t)
	})

	t.Run("returns nil when no rows found", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("GetUserByEmail", ctx, email).Return(repository.User{}, sql.ErrNoRows)

		userDto, err := userService.GetUserByEmail(ctx, email)

		require.Error(t, err)
		assert.Nil(t, userDto)
		assert.Equal(t, err, user.ErrUserNotFound)

		mockQueries.AssertExpectations(t)
	})

	t.Run("returns error when query fails", func(t *testing.T) {
		mockQueries, _, userService := setupUserServiceTest(t)
		mockQueries.On("GetUserByEmail", ctx, email).Return(repository.User{}, errors.New("database error"))

		_, err := userService.GetUserByEmail(ctx, email)

		require.Error(t, err)
		assert.Equal(t, "database error", err.Error())

		mockQueries.AssertExpectations(t)
	})
}
