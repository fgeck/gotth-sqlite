//go:build unittest

package password_test

import (
	"errors"
	"testing"

	"github.com/fgeck/go-register/internal/service/security/password"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	mockHashFunc := func(password []byte, cost int) ([]byte, error) {
		if string(password) == "error" {
			return nil, errors.New("mock hashing error")
		}
		return []byte("mockHashedPassword"), nil
	}

	service := password.NewPasswordServiceWithCustomFuncs(mockHashFunc, nil)

	t.Run("successfully hashes password", func(t *testing.T) {
		hashedPassword, err := service.HashAndSaltPassword("validPassword")
		require.NoError(t, err)
		assert.Equal(t, "mockHashedPassword", hashedPassword)
	})

	t.Run("fails to hash password", func(t *testing.T) {
		_, err := service.HashAndSaltPassword("error")
		require.Error(t, err)
		assert.Equal(t, "mock hashing error", err.Error())
	})
}

func TestComparePassword(t *testing.T) {
	mockCompareFunc := func(hashedPassword, password []byte) error {
		if string(password) == "wrongPassword" {
			return errors.New("mock invalid password")
		}
		return nil
	}

	service := password.NewPasswordServiceWithCustomFuncs(nil, mockCompareFunc)

	t.Run("successfully compares password", func(t *testing.T) {
		err := service.ComparePassword("mockHashedPassword", "validPassword")
		require.NoError(t, err)
	})

	t.Run("fails to compare password", func(t *testing.T) {
		err := service.ComparePassword("mockHashedPassword", "wrongPassword")
		require.Error(t, err)
		assert.Equal(t, "mock invalid password", err.Error())
	})
}
