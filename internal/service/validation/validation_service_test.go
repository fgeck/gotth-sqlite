//go:build unittest

package validation_test

import (
	"testing"

	validation "github.com/fgeck/go-register/internal/service/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		email    string
		expected error
	}{
		{"valid.email@example.com", nil},
		{"invalid-email", validation.ErrInvalidEmailFormat},
		{"", validation.ErrInvalidEmailFormat},
	}

	for _, test := range tests {
		err := vs.ValidateEmail(test.email)
		if test.expected == nil {
			require.NoError(t, err, "expected no error for email: %s", test.email)
		} else {
			require.Error(t, err, "expected an error for email: %s", test.email)
			assert.Equal(t, test.expected, err, "unexpected error for email: %s", test.email)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		password string
		expected error
	}{
		{"SuperVal!d1@", nil},
		{"Valid1@", validation.ErrInvalidPassword},
		{"short", validation.ErrInvalidPassword},
		{"NoSpecialChar1", validation.ErrInvalidPassword},
		{"nouppercase1@", validation.ErrInvalidPassword},
	}

	for _, test := range tests {
		err := vs.ValidatePassword(test.password)
		if test.expected == nil {
			require.NoError(t, err, "expected no error for password: %s", test.password)
		} else {
			require.Error(t, err, "expected an error for password: %s", test.password)
			assert.Equal(t, test.expected, err, "unexpected error for password: %s", test.password)
		}
	}
}

func TestValidateUsername(t *testing.T) {
	vs := validation.NewValidationService()

	tests := []struct {
		username string
		expected error
	}{
		{"validUser", nil},
		{"val1dUs3r", nil},
		{"ab", validation.ErrInvalidUsername},
		{"invalid_user!", validation.ErrInvalidUsername},
	}

	for _, test := range tests {
		err := vs.ValidateUsername(test.username)
		if test.expected == nil {
			require.NoError(t, err, "expected no error for username: %s", test.username)
		} else {
			require.Error(t, err, "expected an error for username: %s", test.username)
			assert.Equal(t, test.expected, err, "unexpected error for username: %s", test.username)
		}
	}
}
