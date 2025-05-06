//go:build unittest

package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/fgeck/gotth-sqlite/internal/service/security/jwt"
	"github.com/fgeck/gotth-sqlite/internal/service/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TEST_SECRET = "SuperVal!d1@asdawe36"
)

func generatePrivateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "Failed to generate RSA private key")
	return privateKey
}

func TestGenerateToken(t *testing.T) {
	t.Parallel()
	jwtService := jwt.NewJwtService(TEST_SECRET, "test-issuer", 3600)

	t.Run("Valid Input", func(t *testing.T) {
		t.Parallel()
		userDto := &user.UserDto{
			ID:   uuid.New(),
			Role: user.UserRoleAdmin,
		}

		token, err := jwtService.GenerateToken(userDto)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Empty User ID", func(t *testing.T) {
		t.Parallel()
		userDto := &user.UserDto{
			ID:   uuid.Nil,
			Role: user.UserRoleAdmin,
		}

		token, err := jwtService.GenerateToken(userDto)
		require.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Empty User Role", func(t *testing.T) {
		t.Parallel()
		userDto := &user.UserDto{
			ID:   uuid.New(),
			Role: user.UserRole{Name: ""},
		}

		token, err := jwtService.GenerateToken(userDto)
		require.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestValidateAndExtractClaims(t *testing.T) {
	t.Parallel()
	jwtService := jwt.NewJwtService(TEST_SECRET, "test-issuer", 3600)

	t.Run("Valid Token", func(t *testing.T) {
		t.Parallel()
		userID := uuid.New()
		userDto := &user.UserDto{
			ID:   userID,
			Role: user.UserRoleAdmin,
		}

		token, err := jwtService.GenerateToken(userDto)
		require.NoError(t, err)

		extractedClaims, err := jwtService.ValidateAndExtractClaims(token)
		require.NoError(t, err)
		assert.Equal(t, userID.String(), extractedClaims.UserId)
		assert.Equal(t, user.UserRoleAdmin.Name, extractedClaims.UserRole)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		t.Parallel()
		token := "invalid-token"

		extractedUser, err := jwtService.ValidateAndExtractClaims(token)
		require.Error(t, err)
		assert.Nil(t, extractedUser)
	})

	t.Run("No HMAC Signing Method Used", func(t *testing.T) {
		t.Parallel()
		privateKey := generatePrivateKey(t)

		// Create a token with RS256 signing method
		claims := gojwt.MapClaims{
			"userId":   uuid.New().String(),
			"userRole": "admin",
			"iss":      "test-issuer",
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Hour).Unix(),
			"nbf":      time.Now().Unix(),
		}

		token := gojwt.NewWithClaims(gojwt.SigningMethodRS256, claims)
		signedToken, err := token.SignedString(privateKey)
		require.NoError(t, err)

		// Validate the token using the JwtService
		extractedClaims, err := jwtService.ValidateAndExtractClaims(signedToken)
		require.Error(t, err)
		assert.Nil(t, extractedClaims)
		assert.Contains(t, err.Error(), "unexpected signing method")
	})

	t.Run("No User ID in Parsed Token", func(t *testing.T) {
		t.Parallel()
		// Create a token without a USER_ID claim
		claims := gojwt.MapClaims{
			"userRole": "admin",
			"iss":      "test-issuer",
			"iat":      time.Now().Unix(),
			"exp":      time.Now().Add(time.Hour).Unix(),
		}

		token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(TEST_SECRET))
		require.NoError(t, err)

		extractedClaims, err := jwtService.ValidateAndExtractClaims(signedToken)
		require.Error(t, err)
		assert.Nil(t, extractedClaims)
		assert.Contains(t, err.Error(), "missing userId claim")
	})

	t.Run("No User Role in Parsed Token", func(t *testing.T) {
		t.Parallel()
		// Create a token without a USER_ROLE claim
		claims := gojwt.MapClaims{
			"userId": uuid.New().String(),
			"iss":    "test-issuer",
			"iat":    time.Now().Unix(),
			"exp":    time.Now().Add(time.Hour).Unix(),
		}

		token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(TEST_SECRET))
		require.NoError(t, err)

		extractedClaims, err := jwtService.ValidateAndExtractClaims(signedToken)
		require.Error(t, err)
		assert.Nil(t, extractedClaims)
		assert.Contains(t, err.Error(), "missing userRole claim")
	})
}
