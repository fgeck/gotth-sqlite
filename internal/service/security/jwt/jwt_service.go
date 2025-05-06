package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/fgeck/go-register/internal/service/user"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtServiceInterface interface {
	GenerateToken(user *user.UserDto) (string, error)
	ValidateAndExtractClaims(givenToken string) (*JwtCustomClaims, error)
}

type JwtService struct {
	secretKey  string
	issuer     string
	expiration int64
}

func NewJwtService(secretKey, issuer string, expiration int64) *JwtService {
	return &JwtService{
		secretKey:  secretKey,
		issuer:     issuer,
		expiration: expiration,
	}
}

var (
	ErrEmptyUserRole        = errors.New("userRole role is empty")
	ErrEmptyUserId          = errors.New("userId is empty")
	ErrInvalidTokenClaims   = errors.New("invalid token claims")
	ErrMissingUserIdClaim   = errors.New("missing userId claim")
	ErrMissingUserRoleClaim = errors.New("missing userRole claim")
	ErrInvalidClaims        = errors.New("userId or userRole claim is nil")
)

func (s *JwtService) GenerateToken(user *user.UserDto) (string, error) {
	if user.ID == uuid.Nil || user.ID.String() == "" {
		return "", ErrEmptyUserId
	}
	if user.Role.Name == "" {
		return "", ErrEmptyUserRole
	}

	now := time.Now()
	claims := NewJwtCustomClaims(
		user.ID.String(),
		user.Role.Name,
		gojwt.RegisteredClaims{
			Issuer:    s.issuer,
			IssuedAt:  gojwt.NewNumericDate(now),
			ExpiresAt: gojwt.NewNumericDate(now.Add(time.Duration(s.expiration) * time.Second)),
			NotBefore: gojwt.NewNumericDate(now),
		},
	)

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) ValidateAndExtractClaims(givenToken string) (*JwtCustomClaims, error) {
	token, err := gojwt.ParseWithClaims(
		givenToken,
		&JwtCustomClaims{},
		func(token *gojwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
				//nolint
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims if the token is valid
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		if claims.UserId == "" {
			return nil, ErrMissingUserIdClaim
		}
		if claims.UserRole == "" {
			return nil, ErrMissingUserRoleClaim
		}
		return claims, nil
	}

	return nil, ErrInvalidTokenClaims
}
